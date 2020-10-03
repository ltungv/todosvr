package service

import (
	"log"
	"net/http"
	"time"
)

// Middleware tạo một handler mới dựa trên handler được truyền vào.
type Middleware func(http.Handler) http.Handler

// RequestLogger in ra thông tin của yêu cầu vừa được xử lý.
func RequestLogger(logger *log.Logger) Middleware {
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

// HeaderResponseContentTypeJSON chèn thông tin định dạng JSON vào tiêu đề của gói tin được gửi đi.
func HeaderResponseContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
