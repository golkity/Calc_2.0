package stuct

type OrchestratorTestCase struct {
	Name         string  `yaml:"name"`
	Operation    string  `yaml:"operation"`
	Expression   string  `yaml:"expression,omitempty"`
	ExpressionID string  `yaml:"expressionID,omitempty"`
	TaskID       int     `yaml:"taskID,omitempty"`
	Result       float64 `yaml:"result,omitempty"`
	Expected     struct {
		ExpressionCount  int     `yaml:"expressionCount,omitempty"`
		TaskCount        int     `yaml:"taskCount,omitempty"`
		ExpressionStatus string  `yaml:"expressionStatus,omitempty"`
		ExpressionResult float64 `yaml:"expressionResult,omitempty"`
		TaskExists       *bool   `yaml:"taskExists,omitempty"`
	} `yaml:"expected,omitempty"`
	ExpectedError string `yaml:"expectedError,omitempty"`
	PreOperation  *struct {
		Op         string `yaml:"op"`
		Expression string `yaml:"expression"`
	} `yaml:"preOperation,omitempty"`
}

type OrchestratorTestSuite struct {
	Tests []OrchestratorTestCase `yaml:"tests"`
}
