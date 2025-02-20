package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ngtrvu/zen-go/log"
)

var Publisher *pubsub.Topic

func InitPublisher(ctx context.Context, projectID, topic string) {
	if Publisher == nil {
		client, err := pubsub.NewClient(ctx, projectID)
		if err != nil {
			log.Error("failed to create pubsub client: %v", err)
			panic(fmt.Sprintf("failed to create pubsub client: %v", err))
		}

		Publisher = client.Topic(topic)
	}
}

func SendDelayTask(ctx context.Context, taskMessage TaskMessage) (string, error) {
	if Publisher == nil {
		log.Error("Publisher is not initialized")
		return "", fmt.Errorf("Publisher is not initialized")
	}

	log.Info("AddTask %v", taskMessage)

	taskData, err := json.Marshal(taskMessage)
	if err != nil {
		log.Error("failed to marshal task: %v", err)
		return "", err
	}

	result := Publisher.Publish(ctx, &pubsub.Message{
		Data: taskData,
	})

	// Get the result and handle any errors
	id, err := result.Get(ctx)
	if err != nil {
		log.Error("failed to publish message: %v", err)
		return "", err
	}

	return id, nil
}
