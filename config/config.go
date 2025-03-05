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
