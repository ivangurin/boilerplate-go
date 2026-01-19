package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"boilerplate/internal/model"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/pkg/metadata"
)

const (
	headerKey       = "X-Key"
	headerRequestID = "X-Request-ID"
	headerOrgID     = "X-Org-ID"
	headerUserID    = "X-User-ID"
	headerIP        = "X-IP"
)

type client struct {
	logger   logger_pkg.Logger
	name     string
	url      string
	conn     net.Conn
	nc       *nats.Conn
	js       jetstream.JetStream
	contexts []jetstream.ConsumeContext
}

func NewClient(logger logger_pkg.Logger, opts ...Option) (model.BrokerClient, error) {
	c := &client{
		logger: logger,
	}

	for _, opt := range opts {
		opt(c)
	}

	var err error
	if c.conn == nil {
		c.nc, err = nats.Connect(c.url, nats.Name(c.name))
		if err != nil {
			return nil, fmt.Errorf("create nats connection: %w", err)
		}
	} else {
		c.nc, err = nats.Connect(c.url, nats.Name(c.name), nats.InProcessServer(c))
		if err != nil {
			return nil, fmt.Errorf("create nats connection: %w", err)
		}
	}

	c.js, err = jetstream.New(c.nc)
	if err != nil {
		return nil, fmt.Errorf("create nats jetstream: %w", err)
	}

	return c, nil
}

func (c *client) Publish(ctx context.Context, topic string, partition *int, key, data any) error {
	c.logger.DebugKV(ctx, "publish to topic", "topic", topic, "partition", partition, "key", key)

	if partition != nil {
		topic = fmt.Sprintf("%s.p%d", topic, *partition)
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data for subject %s: %w", topic, err)
	}

	// Создаем сообщение с заголовками из контекста
	msg := nats.NewMsg(topic)
	msg.Data = dataBytes

	msg.Header.Add(headerKey, fmt.Sprintf("%v", key))

	// Извлекаем метаданные из контекста и добавляем в заголовки
	if requestID, ok := metadata.GetRequestID(ctx); ok {
		msg.Header.Add(headerRequestID, requestID)
	}
	if userID, ok := metadata.GetUserID(ctx); ok {
		msg.Header.Add(headerUserID, fmt.Sprintf("%d", userID))
	}
	if ip, ok := metadata.GetIP(ctx); ok {
		msg.Header.Add(headerIP, ip)
	}

	pa, err := c.js.PublishMsg(ctx, msg)
	if err != nil {
		return fmt.Errorf("publish to topic %s: %w", topic, err)
	}

	c.logger.DebugKV(ctx, "published to topic", "topic", topic, "sequence", pa.Sequence)

	return nil
}

// nolint: gocognit
func (c *client) Subscribe(ctx context.Context, consumerName, description string, topic model.BrokerTopic, handler model.BrokerHandler) error {
	c.logger.InfoKV(ctx, "subscribing to topic", "consumer", consumerName, "topic", topic.Name)

	stream, err := c.js.Stream(ctx, topic.Name)
	if err != nil {
		return fmt.Errorf("get stream for topic %s: %w", topic.Name, err)
	}

	info, err := stream.Info(ctx)
	if err != nil {
		return fmt.Errorf("get stream info %s: %w", topic.Name, err)
	}

	for _, subject := range info.Config.Subjects {
		cn := consumerName + "-" + subject
		cn = strings.ReplaceAll(cn, ".", "-")
		natsConsumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
			Name:          cn,
			Durable:       cn,
			Description:   description,
			FilterSubject: subject,
			AckPolicy:     jetstream.AckExplicitPolicy,
			MaxDeliver:    topic.Retries + 1,
			DeliverPolicy: jetstream.DeliverAllPolicy,
		})
		if err != nil {
			return fmt.Errorf("create or update consumer for subject %s: %w", subject, err)
		}

		natsContext, err := natsConsumer.Consume(func(msg jetstream.Msg) {
			md, err := msg.Metadata()
			if err != nil {
				c.logger.ErrorKV(ctx, "get message metadata error", "consumer", consumerName, "subject", msg.Subject(), "error", err.Error())
				return
			}

			key := msg.Headers().Get(headerKey)

			requestID := msg.Headers().Get(headerRequestID)
			if requestID != "" {
				ctx = metadata.WithRequestID(ctx, requestID)
			}

			userIDstr := msg.Headers().Get(headerUserID)
			if userIDstr != "" {
				userID, err := strconv.Atoi(userIDstr)
				if err != nil {
					c.logger.ErrorKV(ctx, "invalid user ID in message header", "consumer", consumerName, "subject", msg.Subject(), "error", err.Error())
				} else {
					ctx = metadata.WithUserID(ctx, userID)
				}
			}

			ip := msg.Headers().Get(headerIP)
			if ip != "" {
				ctx = metadata.WithIP(ctx, ip)
			}

			c.logger.DebugKV(ctx, "message received", "consumer", cn, "subject", msg.Subject())

			err = handler(ctx, msg.Subject(), msg.Data())
			if err != nil {
				c.logger.ErrorKV(ctx, "handle message error", "consumer", cn, "subject", msg.Subject(), "error", err.Error())
				if md.NumDelivered >= uint64(topic.Retries) {
					c.logger.WarnKV(ctx, "message reached max delivery attempts", "consumer", cn, "subject", msg.Subject(), "attempts", md.NumDelivered)

					// Публикуем в DLQ перед Ack
					err = c.Publish(ctx, topic.DLQTopicName, nil, key, msg.Data())
					if err != nil {
						c.logger.ErrorKV(ctx, "publish to DLQ error", "consumer", cn, "subject", msg.Subject(), "error", err.Error())
						return
					}
					c.logger.InfoKV(ctx, "message sent to DLQ", "consumer", cn, "subject", msg.Subject(), "dlq_topic", topic.DLQTopicName)

					if err = msg.Ack(); err != nil {
						c.logger.ErrorKV(ctx, "ack message error", "consumer", cn, "subject", msg.Subject(), "error", err.Error())
						return
					}

					return
				}

				if topic.Retries > 0 && md.NumDelivered < uint64(topic.Retries) {
					if err := msg.NakWithDelay(topic.RetriesDelay); err != nil {
						c.logger.ErrorKV(ctx, "nak with delay message", "consumer", cn, "subject", msg.Subject(), "error", err.Error())
						return
					}
					c.logger.DebugKV(ctx, "nak message with delay", "consumer", cn, "subject", msg.Subject(), "attempts", md.NumDelivered)
					return
				}

				if err := msg.TermWithReason(err.Error()); err != nil {
					c.logger.ErrorKV(ctx, "term message error", "consumer", cn, "subject", msg.Subject(), "error", err.Error())
				}

				return
			}

			if err = msg.Ack(); err != nil {
				c.logger.ErrorKV(ctx, "ack message error", "consumer", cn, "subject", msg.Subject(), "error", err.Error())
				return
			}

			c.logger.DebugKV(ctx, "message processed successfully", "consumer", cn, "subject", msg.Subject())
		})
		if err != nil {
			return fmt.Errorf("start to consume messages for subject %s: %w", subject, err)
		}

		c.logger.DebugKV(ctx, "subscribed to subject", "consumer", cn, "subject", subject)

		c.contexts = append(c.contexts, natsContext)
	}

	return nil
}

func (c *client) InProcessConn() (net.Conn, error) {
	return c.conn, nil
}

func (c *client) Close() error {
	for _, ctx := range c.contexts {
		ctx.Stop()
	}
	c.contexts = nil
	c.nc.Close()
	return nil
}

func (c *client) CreateOrUpdateTopic(ctx context.Context, topic model.BrokerTopic) error {
	subjects := []string{}
	if topic.Partitions > 0 {
		for i := 0; i < topic.Partitions; i++ {
			subjects = append(subjects, fmt.Sprintf("%s.p%d", topic.Name, i))
		}
	} else {
		subjects = append(subjects, topic.Name)
	}

	stream, err := c.js.Stream(ctx, topic.Name)
	if err == nil {
		info, err := stream.Info(ctx)
		if err != nil {
			return fmt.Errorf("get stream info %s: %w", topic.Name, err)
		}

		info.Config.Subjects = subjects
		info.Config.MaxBytes = topic.MaxBytes
		info.Config.MaxAge = topic.MaxAge

		_, err = c.js.UpdateStream(ctx, info.Config)
		if err != nil {
			return fmt.Errorf("update stream %s with new subject: %w", topic.Name, err)
		}

		return nil
	}

	_, err = c.js.CreateStream(ctx, jetstream.StreamConfig{
		Name:        topic.Name,
		Description: topic.Description,
		Subjects:    subjects,
		MaxBytes:    topic.MaxBytes,
		MaxAge:      topic.MaxAge,
	})
	if err != nil {
		return fmt.Errorf("create stream %s: %w", topic.Name, err)
	}

	return nil
}
