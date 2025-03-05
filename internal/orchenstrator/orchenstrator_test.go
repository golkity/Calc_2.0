package orchestrator

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	stuct "github.com/golkity/Calc_2.0/config/struct"
	"github.com/golkity/Calc_2.0/internal/store"
	"github.com/golkity/Calc_2.0/internal/task"
)

func contains(str, substr string) bool {
	return len(substr) > 0 && strings.Contains(str, substr)
}

func resetStore() {
	store.StoreMu.Lock()
	defer store.StoreMu.Unlock()
	store.NextExprID = 0
	store.NextTaskID = 0
	store.DB.Expressions = make(map[string]*task.Expression)
	store.DB.Tasks = make(map[int]*task.Task)
}

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
	banner := `░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░ЗАПУСКАЕМ░ГУСЕЙ-РАЗВЕДЧИКОВ░░░░
░░░░░▄▀▀▀▄░░░▄▀▀▀▀▄░░░▄▀▀▀▄░░░░░
▄███▀░◐░░░▌░▐0░░░░0▌░▐░░░◐░▀███▄
░░░░▌░░░░░▐░▌░▐▀▀▌░▐░▌░░░░░▐░░░░
░░░░▐░░░░░▐░▌░▌▒▒▐░▐░▌░░░░░▌░░░░`
	fmt.Println(banner)

	exitCode := m.Run()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nTEST NAME\tSTATUS\tTIME (ms)")
	fmt.Fprintln(w, "---------\t------\t---------")
	var passedCount int
	for _, res := range testResults {
		status := "PASS"
		if !res.Passed {
			status = "FAIL"
		} else {
			passedCount++
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", res.Name, status, res.Duration.Milliseconds())
	}
	w.Flush()
	totalTests := len(testResults)
	overallRating := float64(passedCount) / float64(totalTests) * 100
	fmt.Printf("\nШкала выполнения тестов: [%s] (%.2f%%)\n", strings.Repeat("█", int(overallRating/5))+strings.Repeat("░", 20-int(overallRating/5)), overallRating)
	fmt.Printf("\nОбщая оценка: %d из %d тестов пройдено (%.2f%%)\n", passedCount, totalTests, overallRating)
	os.Exit(exitCode)
}

func TestOrchestrator(t *testing.T) {
	suite, err := stuct.LoadYML[stuct.OrchestratorTestSuite]("orchestrator_test.yaml")
	if err != nil {
		t.Fatalf("Error loading YAML file: %v", err)
	}
	resetStore()

	for _, tc := range suite.Tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()

			if tc.PreOperation != nil {
				switch tc.PreOperation.Op {
				case "AddNewExpression":
					id, err := AddNewExpression(tc.PreOperation.Expression)
					if err != nil {
						t.Fatalf("PreOperation AddNewExpression error: %v", err)
					}
					if tc.ExpressionID == "" {
						tc.ExpressionID = id
					}
				}
			}

			switch tc.Operation {
			case "AddNewExpression":
				id, err := AddNewExpression(tc.Expression)
				if tc.ExpectedError != "" {
					if err == nil {
						t.Errorf("Expected error %q, got nil", tc.ExpectedError)
					} else if !contains(err.Error(), tc.ExpectedError) {
						t.Errorf("Expected error %q, got %q", tc.ExpectedError, err.Error())
					}
					goto record
				}
				if err != nil {
					t.Errorf("AddNewExpression returned error: %v", err)
					goto record
				}
				if _, err := strconv.Atoi(id); err != nil {
					t.Errorf("Invalid expression ID format: %s", id)
				}
			case "GetExpressionByID":
				expr, err := GetExpressionByID(tc.ExpressionID)
				if tc.ExpectedError != "" {
					if err == nil {
						t.Errorf("Expected error %q, got nil", tc.ExpectedError)
					} else if !contains(err.Error(), tc.ExpectedError) {
						t.Errorf("Expected error %q, got %q", tc.ExpectedError, err.Error())
					}
					goto record
				}
				if err != nil {
					t.Errorf("GetExpressionByID returned error: %v", err)
					goto record
				}
				if expr.Status != tc.Expected.ExpressionStatus {
					t.Errorf("Expected expression status %q, got %q", tc.Expected.ExpressionStatus, expr.Status)
				}
			case "CompleteTask":
				err := CompleteTask(tc.TaskID, tc.Result)
				if tc.ExpectedError != "" {
					if err == nil {
						t.Errorf("Expected error %q, got nil", tc.ExpectedError)
					} else if !contains(err.Error(), tc.ExpectedError) {
						t.Errorf("Expected error %q, got %q", tc.ExpectedError, err.Error())
					}
					goto record
				}
				if err != nil {
					t.Errorf("CompleteTask returned error: %v", err)
					goto record
				}
				expr, err := GetExpressionByID(tc.ExpressionID)
				if err != nil {
					t.Errorf("Error retrieving expression: %v", err)
					goto record
				}
				if expr.Status != tc.Expected.ExpressionStatus {
					t.Errorf("Expected expression status %q, got %q", tc.Expected.ExpressionStatus, expr.Status)
				}
				if expr.Result == nil {
					t.Errorf("Expected expression result, got nil")
				} else if *expr.Result != tc.Expected.ExpressionResult {
					t.Errorf("Expected expression result %v, got %v", tc.Expected.ExpressionResult, *expr.Result)
				}
			case "GetNextTask":
				task := GetNextTask()
				if tc.Expected.TaskExists != nil {
					if *tc.Expected.TaskExists && task == nil {
						t.Errorf("Expected next task to exist, but got nil")
					}
					if !*tc.Expected.TaskExists && task != nil {
						t.Errorf("Expected no next task, but got one")
					}
				}
			default:
				t.Errorf("Unknown operation %q", tc.Operation)
			}
		record:
			duration := time.Since(start)
			recordResult(tc.Name, !t.Failed(), duration)
		})
	}
}
