package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// NOTE: Using slog as its able to produce JSON logs which is good for production
func Logging(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.NewString() // lets add a identifier
			logger.Info("request started", "method", r.Method, "path", r.URL.Path, "request_id", requestID)
			next.ServeHTTP(w, r)
			logger.Info("request completed", "method", r.Method, "path", r.URL.Path, "request_id", requestID, "duration", time.Since(start))
		})
	}
}
