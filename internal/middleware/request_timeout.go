package middleware

import (
	"net/http"
	"time"
)

type TimeoutRequestMiddleware struct {
	timeout time.Duration
}

func NewTimeoutRequestMiddleware(d time.Duration) *TimeoutRequestMiddleware {
	return &TimeoutRequestMiddleware{
		timeout: d,
	}
}
func (m *TimeoutRequestMiddleware) RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.TimeoutHandler(next, m.timeout, "Request timeout").ServeHTTP(w, r)
	})
}
