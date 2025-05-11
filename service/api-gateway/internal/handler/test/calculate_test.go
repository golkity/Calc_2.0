package test

import (
	"api-gateway/internal/handler"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golkity/Calc_2.0/testsuite"
)

type badCase struct {
	Name string `yaml:"name"`
	Body string `yaml:"body"`
	Want int    `yaml:"want"`
}

func TestMain(m *testing.M) { testsuite.WrapMain(m) }

func TestCalculate_BadRequests(t *testing.T) {
	suite, err := testsuite.Load[badCase]("testdata/calculate_bad.yaml")
	if err != nil {
		t.Fatal(err)
	}

	h := handler.Calculate(nil, nil)

	testsuite.Run(t, suite.Tests, func(tc badCase, t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.Body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != tc.Want {
			t.Fatalf("got %d, want %d (body: %q)", rec.Code, tc.Want, rec.Body.String())
		}
	})
}
