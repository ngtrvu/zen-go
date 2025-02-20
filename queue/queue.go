package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/ngtrvu/zen-go/log"
)

type Queue struct {
	Subscriber *pubsub.Subscription
	Tasks      map[string]TaskInterface
}

func NewQueue(ctx context.Context, cfg QueueConfig) *Queue {
	client, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		log.Error("failed to create pubsub client: %v", err)
		panic(fmt.Sprintf("failed to create pubsub client: %v", err))
	}

	InitPublisher(ctx, cfg.ProjectID, cfg.TopicName)

	sub := client.Subscription(cfg.SubName)
	tasks := make(map[string]TaskInterface)

	return &Queue{Subscriber: sub, Tasks: tasks}
}

func (b *Queue) Start(ctx context.Context) {
	err := b.Subscriber.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Info("received message: %v", string(msg.Data))

		var taskMessage TaskMessage
		err := json.Unmarshal(msg.Data, &taskMessage)
		if err != nil {
			log.Error("failed to unmarshal task: %v", err)
			msg.Ack()
			return
		}

		// search task by task type
		task := b.Tasks[taskMessage.TaskID]
		err = task.Execute(taskMessage.Args)
		if err != nil {
			log.Error("failed to execute task: %v", err)
			msg.Nack()
			return
		}

		// ack message when task executed successfully
		msg.Ack()
	})

	if err != context.Canceled {
		log.Error("cancelled: %v", err)
	}

	if err != nil {
		log.Error("failed to receive message: %v", err)
	}

	log.Info("Queue started")
}

func (b *Queue) AddTask(task TaskInterface) {
	b.Tasks[task.GetTaskID()] = task
	log.Debug("Added task: %v %v", task.GetTaskType(), task.GetTaskID())
}
