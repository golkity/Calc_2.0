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

	orchApp, err := applicant.OrchestratorNew(configPath)
	if err != nil {
		log.Fatalf("Failed to init orchestrator: %v", err)
	}

	orchApp.Run()
}
