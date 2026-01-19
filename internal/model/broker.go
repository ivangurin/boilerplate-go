package model

import (
	"context"
	"net"
	"time"
)

type BrokerServer interface {
	Start() error
	Stop() error
	GetConn() (net.Conn, error)
}

type BrokerClient interface {
	Publish(ctx context.Context, topic string, partition *int, key, data any) error
	Subscribe(ctx context.Context, consumerName, description string, topic BrokerTopic, handler BrokerHandler) error
	CreateOrUpdateTopic(ctx context.Context, topic BrokerTopic) error
	Close() error
}

type BrokerHandler func(ctx context.Context, subject string, data []byte) error

type BrokerConsumer interface {
	Name() string
	Description() string
	MainTopic() string
	DLQTopic() string
	HandleMessage(ctx context.Context, subject string, data []byte) error
}

type BrokerTopic struct {
	Name         string
	Description  string
	Partitions   int
	MaxAge       time.Duration
	MaxBytes     int64
	Retries      int
	RetriesDelay time.Duration
	DLQTopicName string
}
