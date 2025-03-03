package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golkity/Calc_2.0/pkg/calc"
	"github.com/golkity/Calc_2.0/pkg/logger"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type TaskResponse struct {
	Task struct {
		ID            int    `json:"id"`
		Arg1          string `json:"arg1"`
		Arg2          string `json:"arg2"`
		Operation     string `json:"operation"`
		OperationTime int    `json:"operation_time"`
	} `json:"task"`
}

type taskInfo struct {
	ID            int
	Arg1          string
	Arg2          string
	Operation     string
	OperationTime int
}

type Agent struct {
	log    *logger.Logger
	client *http.Client
}

func NewAgent(log *logger.Logger) *Agent {
	return &Agent{
		log: log,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (a *Agent) Start() error {
	cp := getComputingPower()
	var wg sync.WaitGroup
	for i := 0; i < cp; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			a.workerLoop(workerID)
		}(i)
	}
	wg.Wait()
	return nil
}

func (a *Agent) workerLoop(workerID int) {
	for {
		task, err := a.getTask()
		if err != nil {
			a.log.WithError(err).Errorf("Worker %d: error getting task", workerID)
			time.Sleep(2 * time.Second)
			continue
		}
		if task == nil {
			a.log.Infof("Worker %d: no tasks available, sleeping...", workerID)
			time.Sleep(2 * time.Second)
			continue
		}

		a.log.Infof("Worker %d: got task=%d, expr='%s' op='%s'", workerID, task.ID, task.Arg1, task.Operation)

		resultValue, err := calc.Calc(task.Arg1)
		if err != nil {
			a.log.WithError(err).Errorf("Worker %d: calc error for task %d", workerID, task.ID)
			time.Sleep(1 * time.Second)
			continue
		}

		time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

		if err := a.sendResult(task.ID, resultValue); err != nil {
			a.log.WithError(err).Errorf("Worker %d: send result error for task %d", workerID, task.ID)
		} else {
			a.log.Infof("Worker %d: completed task %d, result=%.2f", workerID, task.ID, resultValue)
		}

		time.Sleep(1 * time.Second)
	}
}

func (a *Agent) getTask() (*taskInfo, error) {
	url := orchestratorURL() + "/internal/task"
	resp, err := a.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status %d", resp.StatusCode)
	}

	var body TaskResponse
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decode task error: %v", err)
	}

	return &taskInfo{
		ID:            body.Task.ID,
		Arg1:          body.Task.Arg1,
		Arg2:          body.Task.Arg2,
		Operation:     body.Task.Operation,
		OperationTime: body.Task.OperationTime,
	}, nil
}

func (a *Agent) sendResult(taskID int, result float64) error {
	url := orchestratorURL() + "/internal/task"
	data := map[string]interface{}{
		"id":     taskID,
		"result": result,
	}
	payload, _ := json.Marshal(data)

	resp, err := a.client.Post(url, "application/json", bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("POST %s failed: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}

func orchestratorURL() string {
	if val := os.Getenv("ORCHESTRATOR_HOST"); val != "" {
		return val
	}
	return "http://localhost:8080"
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
