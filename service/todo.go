package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/letung3105/todos-svr/middleware"
	"github.com/letung3105/todos-svr/storage"
	"rsc.io/quote/v3"
)

// Todo chứa các hàm số dùng để khởi tạo và chạy dịch vụ API.
// TODO: thêm database để  lưu trữ dữ liệu người dùng
type Todo struct {
	router  chi.Router
	storage *storage.Todo
}

// NewTodo khởi tạo dịch vụ API.
func NewTodo(storage *storage.Todo) *Todo {
	t := Todo{
		router:  chi.NewRouter(),
		storage: storage,
	}

	t.routes()
	return &t
}

// ServeHTTP dùng để gọi hàm số cung cấp bởi router nhằm thoả mãn interface http.Handler.
func (todo *Todo) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	todo.router.ServeHTTP(w, r)
}

// routes cài đặt các đường dẫn HTTP được hỗ trợ bởi dịch vụ.
// TODO: thêm handlers để xử lý thao tác CRUD
func (todo *Todo) routes() {
	todo.router.Use(
		middleware.JSONResponse(),
		middleware.LogRequest(log.New(os.Stdout, "", log.LstdFlags)),
	)

	// TODO: xử lý PUT, GET và DELETE cho đường dẫn /todo
	todo.router.Get("/hello", todo.GetHello())
	todo.router.Post("/todo", todo.CreateOneTask())
}

// CreateOneTask lấy thông tin của một tác vụ từ body của request rồi gửi cho storage.Todo để xứ lý.
func (todo *Todo) CreateOneTask() http.HandlerFunc {
	type errResponse struct {
		Err string `json:"error"`
	}

	// NOTE: chúng ta không cần kiểm tra lỗi trả về từ hàm số `Encode` vì:
	// + các struct không chứa các kiểu dữ liệu không được hỗ trợ bởi JSON
	// + các struct không phải là cyclic data structure
	// Đọc thêm tại https://golang.org/pkg/encoding/json/#Marshal
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := storage.Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errResponse{err.Error()})
			return
		}

		// TODO: kiểm tra thông tin của `task`
		// TODO: giải thích thư viện "context"

		newTask, err := todo.storage.CreateOne(r.Context(), task)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(errResponse{err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newTask)
	})
}

// GetHello trả về một handler có thể xử lý yêu cầu HTTP và trả về chuỗi ký tự "Hello, world." nằm trong một object JSON.
func (todo *Todo) GetHello() http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(response{quote.HelloV3()}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
}
