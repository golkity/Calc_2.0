package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golkity/Calc_2.0/pkg/tokens"
)

type ctxKey string

const UserIDKey ctxKey = "userID"

func UserID(ctx context.Context) int64 {
	if v, ok := ctx.Value(UserIDKey).(int64); ok {
		return v
	}
	return 0
}

func Auth(tm *tokens.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			c, err := tm.Parse(raw)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, c.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
