package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"

	"boilerplate/internal/model"
	"boilerplate/internal/pkg/logger"
	"boilerplate/internal/pkg/metadata"
)

const (
	headerRequestID = "X-Request-ID"
	headerOrgID     = "X-Org-ID"
	headerUserID    = "X-User-ID"
	headerIP        = "X-IP"
)

type client struct {
	logger       logger.Logger
	name         string
	url          string
	conn         net.Conn
	nc           *nats.Conn
	js           jetstream.JetStream
	streams      map[string]jetstream.Stream
	streamsMutex sync.RWMutex
	contexts     []jetstream.ConsumeContext
}

func NewClient(logger logger.Logger, opts ...Option) (model.BrokerClient, error) {
	c := &client{
		logger:  logger,
		streams: map[string]jetstream.Stream{},
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

func (c *client) Publish(ctx context.Context, subject string, key, data any) error {
	c.logger.DebugKV(ctx, "publish to subject", "subject", subject, "key", key)
	if _, err := c.createOrUpdateStream(ctx, subject); err != nil {
		return fmt.Errorf("create or update stream for subject %s: %w", subject, err)
	}

	fullSubject := fmt.Sprintf("%s.%v", subject, key)

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data for subject %s: %w", subject, err)
	}

	// Создаем сообщение с заголовками из контекста
	msg := nats.NewMsg(fullSubject)
	msg.Data = dataBytes

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
		return fmt.Errorf("publish to subject %s: %w", fullSubject, err)
	}

	c.logger.DebugKV(ctx, "published to subject", "subject", fullSubject, "sequence", pa.Sequence)

	return nil
}

func (c *client) Subscribe(ctx context.Context, consumerName string, subjects []string, handler model.BrokerHandler) error {
	for _, subject := range subjects {
		c.logger.InfoKV(ctx, "subscribing to subject", "subject", subject, "consumer", consumerName)

		stream, err := c.createOrUpdateStream(ctx, subject)
		if err != nil {
			return fmt.Errorf("create or update stream for subject %s: %w", subject, err)
		}

		natsConsumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
			Name:          consumerName,
			Durable:       consumerName,
			Description:   consumerName,
			FilterSubject: subject + ".>",
			AckPolicy:     jetstream.AckExplicitPolicy,
		})
		if err != nil {
			return fmt.Errorf("create or update consumer for subject %s: %w", subject, err)
		}

		natsContext, err := natsConsumer.Consume(func(msg jetstream.Msg) {
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

			c.logger.DebugKV(ctx, "message received", "consumer", consumerName, "subject", msg.Subject())

			err := handler(ctx, msg.Subject(), msg.Data())
			if err != nil {
				c.logger.ErrorKV(ctx, "handle message error", "consumer", consumerName, "subject", msg.Subject(), "error", err.Error())
				if err := msg.TermWithReason(err.Error()); err != nil {
					c.logger.ErrorKV(ctx, "term message error", "consumer", consumerName, "subject", msg.Subject(), "error", err.Error())
				}
				return
			}

			if err = msg.Ack(); err != nil {
				c.logger.ErrorKV(ctx, "ack message error", "consumer", consumerName, "subject", msg.Subject(), "error", err.Error())
				return
			}

			c.logger.DebugKV(ctx, "message processed successfully", "consumer", consumerName, "subject", msg.Subject())
		})
		if err != nil {
			return fmt.Errorf("start to consume messages for subject %s: %w", subject, err)
		}

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

func (c *client) createOrUpdateStream(ctx context.Context, subject string) (jetstream.Stream, error) {
	streamName := strings.ReplaceAll(subject, ".", "_")

	c.streamsMutex.RLock()
	if stream, exists := c.streams[streamName]; exists {
		c.streamsMutex.RUnlock()
		return stream, nil
	}
	c.streamsMutex.RUnlock()

	c.streamsMutex.Lock()
	defer c.streamsMutex.Unlock()

	if stream, exists := c.streams[streamName]; exists {
		return stream, nil
	}

	// Try to get existing stream
	stream, err := c.js.Stream(ctx, streamName)
	if err == nil {
		// Stream exists, check subjects
		info, err := stream.Info(ctx)
		if err != nil {
			return nil, fmt.Errorf("get stream info %s: %w", streamName, err)
		}

		subjectPattern := subject + ".>"
		subjectExists := false
		for _, s := range info.Config.Subjects {
			if s == subjectPattern {
				subjectExists = true
				break
			}
		}

		if !subjectExists {
			// Add new subject to existing ones
			info.Config.Subjects = append(info.Config.Subjects, subjectPattern)
			stream, err = c.js.UpdateStream(ctx, info.Config)
			if err != nil {
				return nil, fmt.Errorf("update stream %s with new subject: %w", streamName, err)
			}
		}

		c.streams[streamName] = stream
		return stream, nil
	}

	// Stream doesn't exist, create new one
	stream, err = c.js.CreateStream(ctx, jetstream.StreamConfig{
		Name:        streamName,
		Description: fmt.Sprintf("Stream for %s", streamName),
		Subjects:    []string{subject + ".>"},
		MaxBytes:    -1,
	})
	if err != nil {
		return nil, fmt.Errorf("create stream %s: %w", streamName, err)
	}

	c.streams[streamName] = stream

	return stream, nil
}
