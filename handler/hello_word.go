package handler

import (
	"encoding/json"
	"net/http"

	"rsc.io/quote/v3"
)

// GetHello xử lý yêu cầu và trả về chuỗi ký tự "Hello, world."
func GetHello() http.HandlerFunc {

	type response struct {
		Message string
	}

	// http.HandlerFunc là một kiểu dữ liệu được định nghĩa dựa trên
	// func(http.ResponseWriter, r *http.Request) và có một method
	// ServeHTTP(http.ResponseWriter, r *http.Request)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(response{Message: quote.HelloV3()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b) // ghi thông tin cho người dùng
	})
}
