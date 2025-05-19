package middleware

import (
	"net/http"
	"time"
)

// NOTE:
// I think the server itself also handled timeout so this might be a duplicate
// Lets investigate more
func Timeout(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, timeout, "Request timed out")
	}
}
