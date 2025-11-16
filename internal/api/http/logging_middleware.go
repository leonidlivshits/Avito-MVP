package httpapi

import (
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/leonidlivshits/Avito-MVP/internal/log"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	size   int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func LoggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.New().String()
			}
			ctx := WithRequestID(r.Context(), reqID)

			sr := &statusRecorder{ResponseWriter: w, status: 0}

			next.ServeHTTP(sr, r.WithContext(ctx))

			duration := time.Since(start)
			role := GetRole(ctx)
			userID := GetUserID(ctx)
			remote := clientIP(r)

			logger.Infow("http_request",
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", sr.status,
				"duration_ms", duration.Milliseconds(),
				"remote", remote,
				"user_id", userID,
				"role", string(role),
				"size", sr.size,
			)
		})
	}
}

func clientIP(r *http.Request) string {
	if xf := r.Header.Get("X-Forwarded-For"); xf != "" {
		return xf
	}
	if xr := r.Header.Get("X-Real-IP"); xr != "" {
		return xr
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
