package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Port string `json:"Port"`
}

func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Error to open jon: %v", err)
	}
	defer file.Close()
	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatalf("Cannot decode config file: %v", err)
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg
}
