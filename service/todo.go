package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/letung3105/todos-svr/middleware"
	"rsc.io/quote/v3"
)

// NewTodo khởi tạo dịch vụ API
func NewTodo() *Todo {
	t := Todo{router: chi.NewRouter()}
	t.routes()
	return &t
}

// Todo chứa các hàm số dùng để khởi tạo và chạy dịch vụ API
type Todo struct {
	router chi.Router
}

// ServeHTTP dùng để gọi hàm số cung cấp bởi router nhằm thoả mãn interface http.Handler
func (t *Todo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.router.ServeHTTP(w, r)
}

// GetHello trả về một handler có thể xử lý yêu cầu HTTP
// và trả về chuỗi ký tự "Hello, world." nằm trong một object JSON
func (t *Todo) GetHello() http.HandlerFunc {
	type response struct {
		Message string
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(response{Message: quote.HelloV3()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(b) // ghi thông tin cho người dùng
	})
}

// routes cài đặt các đường dẫn HTTP được hỗ trợ bởi dịch vụ
func (t *Todo) routes() {
	t.router.Use(
		middleware.JSONResponse(),
		middleware.LogRequest(log.New(os.Stdout, "", log.LstdFlags)),
	)
	t.router.Get("/hello", t.GetHello())
}
