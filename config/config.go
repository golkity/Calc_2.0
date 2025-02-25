package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	ServerPort string `json:"server_port"`
	LogFile    string `json:"log_file"`
}

func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	byteValue, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}

	return &config, nil
}
