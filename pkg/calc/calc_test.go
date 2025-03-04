package calc

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	"github.com/golkity/Calc_2.0/config"
)

type TestResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
}

var (
	testResult []TestResult
	results    sync.Mutex
)

func recordResult(name string, passed bool, d time.Duration) {
	results.Lock()
	defer results.Unlock()
	testResult = append(testResult, TestResult{Name: name, Passed: passed, Duration: d})
}

func TestMain(m *testing.M) {
	goose := `░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
░░░░ЗАПУСКАЕМ░ГУСЕЙ-РАЗВЕДЧИКОВ░░░░
░░░░░▄▀▀▀▄░░░▄▀▀▀▀▄░░░▄▀▀▀▄░░░░░
▄███▀░◐░░░▌░▐0░░░░0▌░▐░░░◐░▀███▄
░░░░▌░░░░░▐░▌░▐▀▀▌░▐░▌░░░░░▐░░░░
░░░░▐░░░░░▐░▌░▌▒▒▐░▐░▌░░░░░▌░░░░`
	fmt.Println(goose)

	exitCode := m.Run()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "\nTEST NAME\tSTATUS\tTIME (ms)")
	fmt.Fprintln(w, "---------\t------\t---------")
	var passedCount int
	for _, result := range testResult {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
		} else {
			passedCount++
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", result.Name, status, result.Duration.Milliseconds())
	}
	w.Flush()
	totalTests := len(testResult)
	overallRating := float64(passedCount) / float64(totalTests) * 100
	fmt.Printf("\nШкала выполнения тестов: [%s] (%.2f%%)\n", strings.Repeat("█", int(overallRating/5))+strings.Repeat("░", 20-int(overallRating/5)), overallRating)
	fmt.Printf("\nОбщая оценка: %d из %d тестов пройдено (%.2f%%)\n", passedCount, totalTests, overallRating)

	os.Exit(exitCode)
}

func TestCalc(t *testing.T) {
	suite, err := config.LoadYML[config.CalcTestSuite]("calc_test.yaml")
	if err != nil {
		t.Fatalf("Error loading YAML file: %v", err)
	}
	for _, tc := range suite.Tests {
		tc := tc
		t.Run(tc.Expression, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			defer func() {
				duration := time.Since(start)
				recordResult(tc.Expression, !t.Failed(), duration)
			}()
			result, err := Calc(tc.Expression)
			if err != nil {
				t.Errorf("Calc(%q) returned error: %v", tc.Expression, err)
				return
			}
			diff := result - tc.Expected
			if diff > 0.0001 || diff < -0.0001 {
				t.Errorf("Calc(%q) = %v, expected %v", tc.Expression, result, tc.Expected)
			}
		})
	}
}
