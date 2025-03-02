package main

import (
	"log"
	"os"

	"lms/internal/applicant"
)

func main() {
	configPath := "./config/config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	agentApp, err := applicant.AgentNew(configPath)
	if err != nil {
		log.Fatalf("Failed to init agent: %v", err)
	}

	agentApp.Run()
}
