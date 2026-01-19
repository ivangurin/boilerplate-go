package topics

import (
	"context"
	"fmt"
	"time"

	"boilerplate/internal/model"
)

const (
	TopicUserCreated    = "user-created"
	TopicUserCreatedDLQ = "user-created-dlq"
)

var Topics = map[string]model.BrokerTopic{
	TopicUserCreated: {
		Name:         TopicUserCreated,
		Description:  "Main topic for user created events",
		Partitions:   3,
		MaxAge:       30 * 24 * time.Hour, // 30 days
		MaxBytes:     1024 * 1024 * 1024,  // 1 GB
		Retries:      3,
		RetriesDelay: time.Duration(5 * time.Second),
		DLQTopicName: TopicUserCreatedDLQ,
	},
	TopicUserCreatedDLQ: {
		Name:        TopicUserCreatedDLQ,
		Description: "DLQ topic for user created events",
		MaxAge:      365 * 24 * time.Hour, // 365 days
		MaxBytes:    1024 * 1024 * 1024,   // 1 GB
	},
}

func CreateOrUpdateTopics(ctx context.Context, client model.BrokerClient) error {
	for _, topic := range Topics {
		if err := client.CreateOrUpdateTopic(ctx, topic); err != nil {
			return fmt.Errorf("create or update topic %s: %w", topic.Name, err)
		}
	}
	return nil
}
