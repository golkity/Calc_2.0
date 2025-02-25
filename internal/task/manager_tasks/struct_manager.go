package manager_tasks

type StatusTask string

const (
	StatusPending    StatusTask = "pending"
	StatusInProgress StatusTask = "in_progress"
	StatusCompleted  StatusTask = "completed"
)

type Task struct {
	ID           string     `json:"id"`
	ExpressionID string     `json:"expression_id"`
	Status       StatusTask `json:"status"`
	Arg1         string     `json:"arg1"`
	Arg2         string     `json:"arg2"`
	Operation    string     `json:"operation"`
	Result       float64    `json:"result,omitempty"`
}

type Expression struct {
	ID     string     `json:"id"`
	Raw    string     `json:"raw"`
	Status StatusTask `json:"status"`
	Result *float64   `json:"result,omitempty"`
}

type TaskResult struct {
	TaskID string  `json:"task_id"`
	Result float64 `json:"result"`
}
