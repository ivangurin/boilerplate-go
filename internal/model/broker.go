package model

import (
	"context"
)

type BrokerServer interface {
	Start() error
	Stop() error
}

type BrokerClient interface {
	Publish(ctx context.Context, subject string, key, data any) error
	Subscribe(ctx context.Context, consumerName string, subjects []string, handler BrokerHandler) error
	Close() error
}

type BrokerHandler func(ctx context.Context, subject string, data []byte) error

type BrokerConsumer interface {
	Name() string
	Description() string
	Subjects() []string
	HandleMessage(ctx context.Context, subject string, data []byte) error
}
