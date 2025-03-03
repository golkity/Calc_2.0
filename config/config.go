package config

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"lms/internal/custom_errors"
	"os"
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

type YAMLTestCases struct {
	Tests []YAMLTestCase `yaml:"tests"`
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

func LoadYML(path string) (*YAMLTestCases, error) {
	data, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var cases YAMLTestCases
	decoder := yaml.NewDecoder(data)
	if err := decoder.Decode(&cases); err != nil {
		return nil, err
	}
	return &cases, nil
}
