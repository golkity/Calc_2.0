package task

type StatusTask string

const (
	StatusPending    StatusTask = "pending"
	StatusInProgress StatusTask = "in_progress"
	StatusCompleted  StatusTask = "completed"
)

type Expression struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Result *float64 `json:"result"`
	Tasks  []int
}

type Task struct {
	ID            int `json:"id"`
	ExpressionID  string
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
	Status        StatusTask
	Done          bool
	Result        *float64
}

type TaskResult struct {
	TaskID string  `json:"task_id"`
	Result float64 `json:"result"`
}
