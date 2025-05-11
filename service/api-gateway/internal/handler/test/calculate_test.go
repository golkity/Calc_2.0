package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"api-gateway/internal/handler"
)

func TestCalculate_BadRequests(t *testing.T) {
	t.Parallel()
	h := handler.Calculate(nil, nil)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"empty body", ``, http.StatusBadRequest},
		{"unterminated JSON object", `{`, http.StatusBadRequest},
		{"unexpected closing", `{"expression":1+1}`, http.StatusBadRequest},
		{"wrong field", `{"expr":"1+1"}`, http.StatusBadRequest},
		{"missing field", `{}`, http.StatusBadRequest},
		{"empty expression", `{"expression":""}`, http.StatusBadRequest},
		{"leading comma", `{, "expression":"1+1"}`, http.StatusBadRequest},
		{"trailing comma", `{"expression":"1+2",}`, http.StatusBadRequest},
		{"numeric expression", `{"expression":123}`, http.StatusBadRequest},
		{"array expression", `{"expression":["1+1"]}`, http.StatusBadRequest},
		{"boolean expression", `{"expression":true}`, http.StatusBadRequest},
		{"object expression", `{"expression":{"expr":"1+1"}}`, http.StatusBadRequest},
		{"whitespace only", `{"expression":"   "}`, http.StatusInternalServerError},
		{"syntax error plus", `{"expression":"1+"}`, http.StatusInternalServerError},
		{"syntax error minus", `{"expression":"2-"}`, http.StatusInternalServerError},
		{"double plus", `{"expression":"1++2"}`, http.StatusInternalServerError},
		{"unbalanced open paren", `{"expression":"(1+2"}`, http.StatusInternalServerError},
		{"invalid chars", `{"expression":"abc"}`, http.StatusInternalServerError},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			if rec.Code != tc.want {
				t.Fatalf("%q: expected status %d, got %d; body=%q",
					tc.name, tc.want, rec.Code, rec.Body.String())
			}
		})
	}
}
