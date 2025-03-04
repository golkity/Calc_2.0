package config

import (
	"encoding/json"
	"os"

	"github.com/golkity/Calc_2.0/internal/custom_errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ServerPort string `json:"Port"`
	LogFile    string `json:"Log"`
}

type YAMLTestCase struct {
	Name   string `yaml:"name"`
	Method string `yaml:"method"`
	Path   string `yaml:"path"`
	Body   string `yaml:"body"`
	Status int    `yaml:"status"`
}

type HandlerTestCases struct {
	Tests []YAMLTestCase `yaml:"tests"`
}

type CalcTestSuite struct {
	Tests []CalcTestCase `yaml:"tests"`
}

type CalcTestCase struct {
	Expression string  `yaml:"expression"`
	Expected   float64 `yaml:"expected"`
}

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

// тут агент структура
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

func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, custom_errors.ErrOpenConfiguration
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, custom_errors.ErrParceConfiguration
	}

	return &config, nil
}

func LoadYML[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result T
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
