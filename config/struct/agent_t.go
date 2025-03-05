package stuct

type AgentTestCases struct {
	Tests []AgentTestCase `yaml:"tests"`
}

type AgentTestCase struct {
	Name     string        `yaml:"name"`
	Simulate SimulateBlock `yaml:"simulate"`
	Task     TaskBlock     `yaml:"task"`
	Expected ExpectedBlock `yaml:"expected"`
}

type SimulateBlock struct {
	GetStatus  int `yaml:"get_status"`
	PostStatus int `yaml:"post_status"`
}

type TaskBlock struct {
	ID            int    `yaml:"id"`
	Arg1          string `yaml:"arg1"`
	Arg2          string `yaml:"arg2"`
	Operation     string `yaml:"operation"`
	OperationTime int    `yaml:"operation_time"`
}

type ExpectedBlock struct {
	Result float64 `yaml:"result,omitempty"`
	Error  string  `yaml:"error,omitempty"`
	NoTask bool    `yaml:"no_task,omitempty"`
}
