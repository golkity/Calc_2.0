package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golkity/Calc_2.0/config"
	"github.com/golkity/Calc_2.0/internal/http/handler"
	"github.com/golkity/Calc_2.0/internal/http/server"
	"github.com/golkity/Calc_2.0/pkg/logger"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"text/tabwriter"
	"time"
)

type TestResult struct {
	Name     string
	Passed   bool
	Duration time.Duration
}

var (
	testResults  []TestResult
	resultsMutex sync.Mutex
)

func recordResult(name string, passed bool, duration time.Duration) {
	resultsMutex.Lock()
	defer resultsMutex.Unlock()
	testResults = append(testResults, TestResult{Name: name, Passed: passed, Duration: duration})
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
	for _, result := range testResults {
		status := "PASS"
		if !result.Passed {
			status = "FAIL"
		} else {
			passedCount++
		}
		fmt.Fprintf(w, "%s\t%s\t%d\n", result.Name, status, result.Duration.Milliseconds())
	}
	w.Flush()
	totalTests := len(testResults)
	overallRating := float64(passedCount) / float64(totalTests) * 100
	fmt.Printf("\nОбщая оценка: %d из %d тестов пройдено (%.2f%%)\n", passedCount, totalTests, overallRating)

	os.Exit(exitCode)
}

func TestFromYAML(t *testing.T) {
	yamlTests, err := config.LoadYML[config.HandlerTestCases]("handler_test.yaml")
	if err != nil {
		t.Fatalf("Ошибка загрузки тестов из YAML: %v", err)
	}

	log := logger.New()
	mux := http.NewServeMux()
	server.RegRoutes(mux, log)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	var wg sync.WaitGroup
	for _, tc := range yamlTests.Tests {
		tc := tc
		wg.Add(1)
		t.Run(tc.Name, func(t *testing.T) {

			start := time.Now()
			defer func() {
				duration := time.Since(start)
				recordResult(tc.Name, !t.Failed(), duration)
				wg.Done()
			}()

			url := ts.URL + tc.Path
			var resp *http.Response
			switch strings.ToUpper(tc.Method) {
			case "POST":
				resp, err = http.Post(url, "application/json", strings.NewReader(tc.Body))
			case "GET":
				resp, err = http.Get(url)
			default:
				req, _ := http.NewRequest(tc.Method, url, strings.NewReader(tc.Body))
				resp, err = http.DefaultClient.Do(req)
			}

			if err != nil {
				t.Errorf("Ошибка запроса: %v", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.Status {
				data, _ := io.ReadAll(resp.Body)
				t.Errorf("Ожидался статус %d, получен %d, тело: %s", tc.Status, resp.StatusCode, data)
			}
		})
	}
	wg.Wait()
}

func TestCalculateDetailed(t *testing.T) {
	log := logger.New()
	h := handler.NewHandler(log)

	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/calculate", h.RegRequest)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	t.Run("CalculateDetailed: success 201", func(t *testing.T) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			recordResult("CalculateDetailed: success 201", !t.Failed(), duration)
		}()

		body := `{"expression":"2+2"}`
		resp, err := http.Post(ts.URL+"/api/v1/calculate", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Ошибка http post: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			data, _ := io.ReadAll(resp.Body)
			t.Fatalf("Ожидался статус 201, получен %d, тело: %s", resp.StatusCode, data)
		}
		var out map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			t.Fatalf("Ошибка декодирования ответа: %v", err)
		}
		if out["id"] == nil {
			t.Error("Ожидается наличие поля 'id' в ответе")
		}
	})

	t.Run("CalculateDetailed: unprocessable 422 invalid body", func(t *testing.T) {
		start := time.Now()
		defer func() {
			duration := time.Since(start)
			recordResult("CalculateDetailed: unprocessable 422 invalid body", !t.Failed(), duration)
		}()

		body := `{"bla":"?"}`
		resp, err := http.Post(ts.URL+"/api/v1/calculate", "application/json", bytes.NewBufferString(body))
		if err != nil {
			t.Fatalf("Ошибка запроса: %v", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusUnprocessableEntity {
			t.Errorf("Ожидался статус 422, получен %d", resp.StatusCode)
		}
	})
}
