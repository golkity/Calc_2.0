package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	stuct "github.com/golkity/Calc_2.0/config/struct"
	"github.com/golkity/Calc_2.0/pkg/calc"
	"github.com/golkity/Calc_2.0/pkg/logger"
)

type TestResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
}

var (
	testResults []TestResult
	resultsMu   sync.Mutex
)

func recordResult(name string, passed bool, d time.Duration) {
	resultsMu.Lock()
	defer resultsMu.Unlock()
	testResults = append(testResults, TestResult{Name: name, Passed: passed, Duration: d})
}

func TestMain(m *testing.M) {
	fmt.Println(`░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░ЗАПУСКАЕМ░ГУСЕЙ-РАЗВЕДЧИКОВ░░░░
░░░░░▄▀▀▀▄░░░▄▀▀▀▀▄░░░▄▀▀▀▄░░░░░
▄███▀░◐░░░▌░▐0░░░░0▌░▐░░░◐░▀███▄
░░░░▌░░░░░▐░▌░▐▀▀▌░▐░▌░░░░░▐░░░░
░░░░▐░░░░░▐░▌░▌▒▒▐░▐░▌░░░░░▌░░░░`)

	exitCode := m.Run()

	totalTests := len(testResults)
	var passedCount int
	for _, res := range testResults {
		if res.Passed {
			passedCount++
		}
	}

	barLength := 20
	ratio := float64(passedCount) / float64(totalTests)
	filled := int(ratio * float64(barLength))

	if filled < 0 {
		filled = 0
	} else if filled > barLength {
		filled = barLength
	}

	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled) + "]"
	fmt.Printf("\nШкала выполнения тестов: %s (%.2f%%)\n", bar, ratio*100)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nTEST NAME\tSTATUS\tTIME (ms)")
	fmt.Fprintln(w, "---------\t------\t---------")
	for _, res := range testResults {
		status := "PASS"
		if !res.Passed {
			status = "FAIL"
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", res.Name, status, res.Duration.Milliseconds())
	}
	w.Flush()
	fmt.Printf("\nОбщая оценка: %d из %d тестов пройдено (%.2f%%)\n", passedCount, totalTests, ratio*100)
	os.Exit(exitCode)
}

func TestAgent(t *testing.T) {
	yamlTests, err := stuct.LoadYML[stuct.AgentTestCases]("agent_test.yaml")
	if err != nil {
		t.Fatalf("Ошибка загрузки тестов из YAML: %v", err)
	}

	log := logger.New()

	for _, tc := range yamlTests.Tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			defer func() {
				duration := time.Since(start)
				recordResult(tc.Name, !t.Failed(), duration)
			}()

			mux := http.NewServeMux()
			mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodGet:
					if tc.Simulate.GetStatus == http.StatusOK {
						w.WriteHeader(http.StatusOK)
						json.NewEncoder(w).Encode(map[string]interface{}{
							"task": map[string]interface{}{
								"id":             tc.Task.ID,
								"arg1":           tc.Task.Arg1,
								"arg2":           tc.Task.Arg2,
								"operation":      tc.Task.Operation,
								"operation_time": tc.Task.OperationTime,
							},
						})
					} else {
						http.Error(w, "No tasks", tc.Simulate.GetStatus)
					}
				case http.MethodPost:
					if tc.Simulate.PostStatus == http.StatusOK {
						w.WriteHeader(http.StatusOK)
						io.WriteString(w, `{"status": "ok"}`)
					} else {
						http.Error(w, "Error processing task", tc.Simulate.PostStatus)
					}
				}
			})
			ts := httptest.NewServer(mux)
			defer ts.Close()

			os.Setenv("ORCHESTRATOR_HOST", ts.URL)
			os.Setenv("COMPUTING_POWER", "1")

			ag := NewAgent(log)

			task, err := ag.getTask()
			if err != nil {
				if tc.Expected.NoTask {
					t.Logf("Ожидалось отсутствие задачи - получено ОК")
					return
				}
				t.Errorf("Ошибка получения задачи: %v", err)
				return
			}

			if tc.Expected.NoTask {
				t.Errorf("Ожидался статус 404 (нет задачи), но получена задача: %+v", task)
				return
			}

			if !tc.Expected.NoTask {
				result, err := calc.Calc(tc.Task.Arg1)
				if err != nil {
					if tc.Expected.Error == "" || !strings.Contains(err.Error(), tc.Expected.Error) {
						t.Errorf("Ожидалась ошибка %q, но получено: %v", tc.Expected.Error, err)
					} else {
						t.Logf("Ожидалась ошибка, получили корректный результат: %v", err)
					}
				} else {
					if result != tc.Expected.Result {
						t.Errorf("Ожидалось значение %v, получено %v", tc.Expected.Result, result)
					} else {
						t.Logf("Ожидалось и получено корректное значение: %v", result)
					}
				}
			}
		})
	}
}
