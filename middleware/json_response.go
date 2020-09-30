package middleware

import "net/http"

// JSONResponse chèn thông tin định dạng JSON vào tiêu đề của gói tin được gửi đi.
func JSONResponse() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	}
}
