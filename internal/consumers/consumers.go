package consumers

import (
	"context"
	"fmt"

	"boilerplate/internal/consumers/user_created"
	"boilerplate/internal/model"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/service_provider"
	"boilerplate/internal/topics"
)

type consumers struct {
	logger    logger_pkg.Logger
	client    model.BrokerClient
	consumers []model.BrokerConsumer
}

func NewConsumers(logger logger_pkg.Logger, client model.BrokerClient, _ *service_provider.Provider) *consumers {
	c := &consumers{
		logger: logger,
		client: client,
	}

	c.consumers = []model.BrokerConsumer{
		user_created.NewConsumer(
			logger.With("consumer", "user_created")),
	}

	return c
}

func (c *consumers) Start(ctx context.Context) error {
	for _, consumer := range c.consumers {
		topic, exists := topics.Topics[consumer.MainTopic()]
		if !exists {
			return fmt.Errorf("topic %s not found", consumer.MainTopic())
		}

		err := c.client.Subscribe(ctx, consumer.Name(), consumer.Description(), topic, consumer.HandleMessage)
		if err != nil {
			return fmt.Errorf("subscribe consumer %s: %w", consumer.Name(), err)
		}
	}

	return nil
}
