package test

import (
	"api-gateway/internal/handler"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubPinger struct {
	err error
}

func (s *stubPinger) Ping(ctx context.Context) error {
	return s.err
}

func TestHealthz_OK(t *testing.T) {
	h := handler.Healthz(&stubPinger{err: nil})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Healthz returned wrong status: got %d, want %d", rr.Code, http.StatusOK)
	}
}

func TestHealthz_KafkaDown(t *testing.T) {
	testErr := fmt.Errorf("kafka unreachable")
	h := handler.Healthz(&stubPinger{err: testErr})
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("Healthz returned wrong status: got %d, want %d", rr.Code, http.StatusServiceUnavailable)
	}

	expectedBody := "kafka down\n"
	if rr.Body.String() != expectedBody {
		t.Errorf("Healthz returned wrong body: got %q, want %q", rr.Body.String(), expectedBody)
	}
}
