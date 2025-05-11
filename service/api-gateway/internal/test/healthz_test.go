package handler_test

import (
	"api-gateway/internal/handler"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockPing struct{ err error }

func (m mockPing) Ping(_ context.Context) error { return m.err }

func TestHealthz_OK(t *testing.T) {
	h := handler.Healthz(mockPing{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
}

func TestHealthz_KafkaDown(t *testing.T) {
	h := handler.Healthz(mockPing{err: errors.New("boom")})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)

	h(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", rec.Code)
	}
}
