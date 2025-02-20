package queue

const (
	TaskTypeAsync = "async_task"
)

type TaskMessage struct {
	TaskID string   `json:"task_id"`
	Args   []string `json:"args"`
}

type TaskInterface interface {
	Execute(args []string) error
	GetTaskID() string
	GetTaskType() string
}

type Task struct {
	TaskID   string `json:"task_id"`
	TaskType string `json:"task_type"`
}

func NewTask(taskID, taskType string) *Task {
	return &Task{
		TaskID:   taskID,
		TaskType: taskType,
	}
}

func (t *Task) Execute(args []string) error {
	panic("implement me")
}

func (t *Task) GetTaskID() string {
	return t.TaskID
}

func (t *Task) GetTaskType() string {
	return t.TaskType
}
