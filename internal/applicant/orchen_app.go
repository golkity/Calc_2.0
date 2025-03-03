package applicant

import (
	"fmt"
	"lms/config"
	"lms/internal/http/server"
	"lms/pkg/logger"
	"net/http"
)

type OrchestratorApp struct {
	config *config.Config
	log    *logger.Logger
}

func OrchestratorNew(configPath string) (*OrchestratorApp, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load config: %v", err)
	}

	log, err := logger.NewFile(cfg.LogFile)
	if err != nil {
		return nil, fmt.Errorf("could not open log file: %v", err)
	}

	return &OrchestratorApp{
		config: cfg,
		log:    log,
	}, nil
}

func (app *OrchestratorApp) Run() {
	mux := http.NewServeMux()
	server.RegRoutes(mux, app.log)

	addr := ":" + app.config.ServerPort
	app.log.Infof("Orchestrator server is running on http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		app.log.Fatalf("Server error: %v", err)
	}
	app.log.Info("Orchestrator server stopped.")
}
