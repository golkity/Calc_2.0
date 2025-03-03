package main

import (
	"lms/internal/applicant"
	"lms/internal/custom_errors"
	"log"
	"os"
)

func main() {
	configPath := "./config/config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	orchApp, err := applicant.OrchestratorNew(configPath)
	if err != nil {
		log.Fatal(custom_errors.ErrInitOrchestrator, err)
	}

	orchApp.Run()
}
