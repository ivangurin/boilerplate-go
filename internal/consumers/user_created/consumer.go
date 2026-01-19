package user_created

import (
	"context"

	"boilerplate/internal/model"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/topics"
)

const (
	Name        = "user-created-consumer"
	Description = "Consumer for handling user created events"
)

type consumer struct {
	logger logger_pkg.Logger
}

func NewConsumer(logger logger_pkg.Logger) model.BrokerConsumer {
	return &consumer{
		logger: logger,
	}
}

func (c *consumer) Name() string {
	return Name
}

func (c *consumer) Description() string {
	return Description
}

func (c *consumer) MainTopic() string {
	return topics.TopicUserCreated
}

func (c *consumer) DLQTopic() string {
	return topics.TopicUserCreatedDLQ
}

func (c *consumer) HandleMessage(_ context.Context, _ string, _ []byte) error {
	return nil
}
