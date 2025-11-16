package httpapi

import (
	"net/http"
	"strings"
)

func TokenAuthMiddleware(adminToken, userToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				tok := strings.TrimPrefix(auth, "Bearer ")
				if tok == adminToken && adminToken != "" {
					r.Header.Set("X-Admin", "1")
				} else if tok == userToken && userToken != "" {
					r.Header.Set("X-User", "1")
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
