package auth

import (
	"net/http"
	"strings"
)

func LoadAuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimSpace(authHeader[7:])
			if token != "" {
				ctx = ContextSet(ctx, "api-key", token)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
