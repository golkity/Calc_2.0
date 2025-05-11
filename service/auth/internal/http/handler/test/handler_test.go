package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"auth_service/internal/http/handler"
)

func TestRegister_BadRequests(t *testing.T) {
	t.Parallel()

	h := handler.New(nil, nil, nil)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"empty body", ``, http.StatusBadRequest},
		{"broken JSON", `{`, http.StatusBadRequest},
		{"closing brace only", `}`, http.StatusBadRequest},
		{"only email key", `{"email":`, http.StatusBadRequest},
		{"only password key", `{"password":`, http.StatusBadRequest},
		{"unterminated email value", `{"email":"a"`, http.StatusBadRequest},
		{"unterminated password value", `{"password":"p"`, http.StatusBadRequest},
		{"numeric email", `{"email":123}`, http.StatusBadRequest},
		{"numeric password", `{"password":123}`, http.StatusBadRequest},
		{"missing quotes email", `{email:"a"}`, http.StatusBadRequest},
		{"missing quotes password", `{password:"p"}`, http.StatusBadRequest},
		{"email without password", `{"email":"a","password"}`, http.StatusBadRequest},
		{"comma instead of colon", `{"email","password":"p"}`, http.StatusBadRequest},
		{"array payload", `[]`, http.StatusBadRequest},
		{"boolean payload", `true`, http.StatusBadRequest},
		{"string payload", `"string"`, http.StatusBadRequest},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			h.Register(rec, req)
			if rec.Code != tc.want {
				t.Errorf("Register %q: got %d, want %d", tc.name, rec.Code, tc.want)
			}
		})
	}
}

func TestLogin_BadRequests(t *testing.T) {
	t.Parallel()

	h := handler.New(nil, nil, nil)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"empty body", ``, http.StatusBadRequest},
		{"broken JSON", `{`, http.StatusBadRequest},
		{"closing brace only", `}`, http.StatusBadRequest},
		{"only email key", `{"email":`, http.StatusBadRequest},
		{"only password key", `{"password":`, http.StatusBadRequest},
		{"unterminated email value", `{"email":"a"`, http.StatusBadRequest},
		{"unterminated password value", `{"password":"p"`, http.StatusBadRequest},
		{"numeric email", `{"email":123}`, http.StatusBadRequest},
		{"numeric password", `{"password":123}`, http.StatusBadRequest},
		{"missing quotes email", `{email:"a"}`, http.StatusBadRequest},
		{"missing quotes password", `{password:"p"}`, http.StatusBadRequest},
		{"email without password", `{"email":"a","password"}`, http.StatusBadRequest},
		{"comma instead of colon", `{"email","password":"p"}`, http.StatusBadRequest},
		{"array payload", `[]`, http.StatusBadRequest},
		{"boolean payload", `false`, http.StatusBadRequest},
		{"string payload", `"foo"`, http.StatusBadRequest},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			h.Login(rec, req)
			if rec.Code != tc.want {
				t.Errorf("Login %q: got %d, want %d", tc.name, rec.Code, tc.want)
			}
		})
	}
}

func TestRefresh_BadRequests(t *testing.T) {
	t.Parallel()

	h := handler.New(nil, nil, nil)

	tests := []struct {
		name string
		body string
		want int
	}{
		{"empty body", ``, http.StatusBadRequest},
		{"broken JSON", `{`, http.StatusBadRequest},
		{"closing brace only", `}`, http.StatusBadRequest},
		{"missing field", `{}`, http.StatusUnauthorized},
		{"wrong field", `{"refresh":""}`, http.StatusUnauthorized},
		{"null token", `{"refresh_token":null}`, http.StatusUnauthorized},
		{"numeric token", `{"refresh_token":123}`, http.StatusBadRequest},
		{"array payload", `["x"]`, http.StatusBadRequest},
		{"boolean payload", `true`, http.StatusBadRequest},
		{"string payload", `"bla"`, http.StatusBadRequest},
		{"unterminated token", `{"refresh_token":"a"`, http.StatusBadRequest},
		{"trailing comma", `{"refresh_token":"a",}`, http.StatusBadRequest},
		{"comma only", `,`, http.StatusBadRequest},
		{"newline only", `"\n"`, http.StatusBadRequest},
		{"unicode only", `"😀"`, http.StatusBadRequest},
		{"empty string", `{"refresh_token":""}`, http.StatusUnauthorized},
		{"invalid json key", `{"invalid":"x"}`, http.StatusUnauthorized},
		{"nested object", `{"refresh_token":{}}`, http.StatusBadRequest},
		{"nested array", `{"refresh_token":[1]}`, http.StatusBadRequest},
		{"spaces only", `"   "`, http.StatusBadRequest},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			h.Refresh(rec, req)
			if rec.Code != tc.want {
				t.Errorf("Refresh %q: got %d, want %d", tc.name, rec.Code, tc.want)
			}
		})
	}
}
