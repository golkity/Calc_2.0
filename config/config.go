package config

import (
	"encoding/json"
	"lms/internal/custom_errors"
	"os"
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
