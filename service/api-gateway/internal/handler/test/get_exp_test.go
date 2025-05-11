package test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"api-gateway/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/golkity/Calc_2.0/service/calc-orchenstrator/repository"
)

type stubRepo struct {
	data map[int64]repository.Expression
}

func (s *stubRepo) One(ctx context.Context, id, userID int64) (repository.Expression, error) {
	if expr, ok := s.data[id]; ok {
		return expr, nil
	}
	return repository.Expression{}, errors.New("not found")
}

func TestGetExpression(t *testing.T) {
	t.Parallel()

	stub := &stubRepo{data: make(map[int64]repository.Expression)}
	for i := int64(1); i <= 20; i += 2 {
		stub.data[i] = repository.Expression{
			ID:     i,
			UserID: 0,
			Raw:    fmt.Sprintf("expr-%d", i),
			Result: nil,
			Status: "done",
		}
	}

	h := handler.GetExpression(stub)

	type testDef struct {
		name     string
		idParam  string
		wantCode int
		wantExpr *repository.Expression
	}
	var tests []testDef
	for i := 1; i <= 20; i++ {
		idStr := strconv.Itoa(i)
		tc := testDef{
			name:     fmt.Sprintf("id=%d", i),
			idParam:  idStr,
			wantCode: http.StatusNotFound,
		}
		if i%2 == 1 {
			tc.wantCode = http.StatusOK
			e := stub.data[int64(i)]
			tc.wantExpr = &e
		}
		tests = append(tests, tc)
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.idParam)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			if rec.Code != tc.wantCode {
				t.Fatalf("для id=%s: ожидаю статус %d, получили %d",
					tc.idParam, tc.wantCode, rec.Code)
			}

			if tc.wantCode == http.StatusOK {
				var resp struct {
					Expression repository.Expression `json:"expression"`
				}
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("не удалось декодировать тело: %v", err)
				}
				if resp.Expression != *tc.wantExpr {
					t.Errorf("для id=%s: получили %+v, ожидаем %+v",
						tc.idParam, resp.Expression, *tc.wantExpr)
				}
			}
		})
	}
}
