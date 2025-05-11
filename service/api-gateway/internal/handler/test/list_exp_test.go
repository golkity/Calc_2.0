package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"api-gateway/internal/handler"
	"api-gateway/internal/middleware"

	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

type stubLister struct {
	data     map[int64][]repository.Expression
	errorFor map[int64]error
}

func (s *stubLister) List(ctx context.Context, userID int64) ([]repository.Expression, error) {
	if err, ok := s.errorFor[userID]; ok {
		return nil, err
	}
	return s.data[userID], nil
}

func equalExprSlices(a, b []repository.Expression) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestListExpressions(t *testing.T) {
	t.Parallel()

	stub := &stubLister{
		data:     make(map[int64][]repository.Expression),
		errorFor: make(map[int64]error),
	}
	for uid := int64(1); uid <= 20; uid++ {
		switch {
		case uid%5 == 0:
			stub.errorFor[uid] = errors.New("db error")
		case uid%2 == 1:
			stub.data[uid] = []repository.Expression{
				{ID: uid*10 + 1, UserID: uid, Raw: fmt.Sprintf("raw-%d-A", uid), Result: nil, Status: "done"},
				{ID: uid*10 + 2, UserID: uid, Raw: fmt.Sprintf("raw-%d-B", uid), Result: nil, Status: "done"},
			}
		default:
			stub.data[uid] = []repository.Expression{}
		}
	}

	h := handler.ListExpressions(stub)

	type tcDef struct {
		uid      int64
		wantCode int
		want     []repository.Expression
	}
	var tests []tcDef
	for uid := int64(1); uid <= 20; uid++ {
		tt := tcDef{uid: uid, wantCode: http.StatusOK}
		if uid%5 == 0 {
			tt.wantCode = http.StatusInternalServerError
		} else {
			tt.want = stub.data[uid]
		}
		tests = append(tests, tt)
	}

	for _, tc := range tests {
		tc := tc
		name := fmt.Sprintf("userID=%d", tc.uid)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctx := context.WithValue(context.Background(), middleware.UserIDKey, tc.uid)
			req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Fatalf("uid=%d: expected status %d, got %d", tc.uid, tc.wantCode, rec.Code)
			}

			if tc.wantCode == http.StatusOK {
				var resp struct {
					Expressions []repository.Expression `json:"expressions"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("uid=%d: decode error: %v", tc.uid, err)
				}
				if !equalExprSlices(resp.Expressions, tc.want) {
					t.Errorf("uid=%d: got %+v, want %+v", tc.uid, resp.Expressions, tc.want)
				}
			}
		})
	}
}
