package handler

import (
	"context"
	"net/http"
	"strings"
)

type ctxKey string

const userIDKey ctxKey = "uid"

func (h *Handler) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		if !strings.HasPrefix(bearer, "Bearer ") {
			http.Error(w, "missing token", 401)
			return
		}
		token := strings.TrimPrefix(bearer, "Bearer ")
		cl, err := h.tm.Parse(token)
		if err != nil {
			http.Error(w, "invalid token", 401)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, cl.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(userIDKey).(int64)
	return id, ok
}
