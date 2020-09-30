package middleware

import (
	"log"
	"net/http"
	"time"
)

// LogRequest in ra thông tin của yêu cầu vừa được xử lý.
func LogRequest(logger *log.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			begin := time.Now()
			next.ServeHTTP(w, r)
			logger.Printf(
				"Handle %s %s from %s in %v",
				r.Method, r.URL, r.RemoteAddr, time.Since(begin),
			)
		})
	}
}
