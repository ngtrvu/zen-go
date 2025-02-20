package queue

type QueueConfig struct {
	TopicName string `config:"BOPQ_TOPIC_NAME"`
	SubName   string `config:"BOPQ_SUB_NAME"`
	ProjectID string `config:"BOPQ_PROJECT_ID"`
}

func NewQueueConfig(TopicName string, SubName string, ProjectID string) QueueConfig {
	return QueueConfig{
		TopicName: TopicName,
		SubName:   SubName,
		ProjectID: ProjectID,
	}
}
