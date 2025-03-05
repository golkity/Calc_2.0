package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"

	stuct "github.com/golkity/Calc_2.0/config/struct"
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
	bar := "[" + strings.Repeat("█", filled) + strings.Repeat("░", barLength-filled) + "]"
	fmt.Printf("\nШкала выполнения тестов: %s (%.2f%%)\n", bar, ratio*100)
	os.Exit(exitCode)
}

func TestLoadConfig(t *testing.T) {
	cases, err := stuct.LoadYML[stuct.ConfigTestCases]("config_test.yaml")
	if err != nil {
		t.Fatalf("Ошибка загрузки YAML-файла с тестами: %v", err)
	}

	for _, tc := range cases.Tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			start := time.Now()
			defer func() {
				duration := time.Since(start)
				recordResult(tc.Name, !t.Failed(), duration)
			}()

			var filePath string
			if tc.FileContent != nil {
				tmpFile, err := ioutil.TempFile("", "config_*.json")
				if err != nil {
					t.Fatalf("Ошибка создания временного файла: %v", err)
				}
				defer os.Remove(tmpFile.Name())

				_, err = tmpFile.Write([]byte(*tc.FileContent))
				if err != nil {
					t.Fatalf("Ошибка записи во временный файл: %v", err)
				}
				tmpFile.Close()
				filePath = tmpFile.Name()
			} else {
				filePath = "nonexistent_file.json"
			}

			config, err := LoadConfig(filePath)
			if tc.ExpectedError != "" {
				if err == nil {
					t.Errorf("Ожидалась ошибка %q, но ошибка не возникла", tc.ExpectedError)
				} else if !strings.Contains(err.Error(), tc.ExpectedError) {
					t.Errorf("Ожидалась ошибка %q, но получена: %v", tc.ExpectedError, err)
				}
				return
			}
			if err != nil {
				t.Errorf("Ошибка загрузки конфигурации: %v", err)
				return
			}
			if tc.Expected == nil {
				t.Errorf("Нет ожидаемых данных для сравнения")
				return
			}
			if config.ServerPort != tc.Expected.ServerPort {
				t.Errorf("Ожидался Port: %s, получено: %s", tc.Expected.ServerPort, config.ServerPort)
			}
			if config.LogFile != tc.Expected.LogFile {
				t.Errorf("Ожидался Log: %s, получено: %s", tc.Expected.LogFile, config.LogFile)
			}
		})
	}
}
