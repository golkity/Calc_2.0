package applicant

import (
	"lms/internal/custom_errors"
	"os"
	"strconv"

	"lms/config"
	"lms/internal/agent"
	"lms/pkg/logger"
)

type AgentApp struct {
	config *config.Config
	log    *logger.Logger
	agent  *agent.Agent
}

func AgentNew(configPath string) (*AgentApp, error) {
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, custom_errors.ErrLoadConfiguration
	}

	log, err := logger.NewFile(cfg.LogFile)
	if err != nil {
		return nil, custom_errors.ErrOpenLogFile
	}
	ag := agent.NewAgent(log)

	return &AgentApp{
		config: cfg,
		log:    log,
		agent:  ag,
	}, nil
}

func (app *AgentApp) Run() {
	app.log.Info("Agent application starting...")

	cp := getComputingPower()
	app.log.Infof("Starting agent with computing power = %d", cp)

	if err := app.agent.Start(); err != nil {
		app.log.Fatal(custom_errors.ErrStartAgent, err)
	}
	app.log.Info("Agent application stopped.")
}

func getComputingPower() int {
	cp := 1
	if val := os.Getenv("COMPUTING_POWER"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			cp = parsed
		}
	}
	return cp
}
