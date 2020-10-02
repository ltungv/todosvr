package service

import (
	"context"
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
	notImpl := func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SimpleMsgResponse{"NOT IMPLEMENTED"})
	}

	todo.router.Use(
		middleware.JSONResponse(),
		middleware.LogRequest(log.New(os.Stdout, "", log.LstdFlags)),
	)
	todo.router.Get("/hello", todo.GetHello())
	todo.router.Post("/todo", todo.CreateOneTask())
	todo.router.Get("/todo", notImpl)
	todo.router.Route("/todo/{id}", func(r chi.Router) {
		r.Use(TaskIDCtx())
		r.Get("/", notImpl)
		r.Put("/", notImpl)
		r.Delete("/", notImpl)
	})
}

// CreateOneTask lấy thông tin của một tác vụ từ body của request rồi gửi cho storage.Todo để xứ lý.
// @Summary Create one task
// @Description tạo một tác vụ mới
// @Accept json
// @Produce json
// @Param task body storage.Task true "tác vụ được thêm vào"
// @Failure 500 {object} service.ErrResponse
// @Failure 400 {object} service.ErrResponse
// @Success 200 {array} storage.Task
// @Router /todo [post]
func (todo *Todo) CreateOneTask() http.HandlerFunc {
	// NOTE: chúng ta không cần kiểm tra lỗi trả về từ hàm số `Encode` vì:
	// + các struct không chứa các kiểu dữ liệu không được hỗ trợ bởi JSON
	// + các struct không phải là cyclic data structure
	// Đọc thêm tại https://golang.org/pkg/encoding/json/#Marshal
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task := storage.Task{}
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrResponse{err.Error()})
			return
		}
		// TODO: kiểm tra thông tin của `task`

		newTask, err := todo.storage.CreateOne(r.Context(), task)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrResponse{err.Error()})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(newTask)
	})
}

// GetHello trả về một handler có thể xử lý yêu cầu HTTP và trả về chuỗi ký tự "Hello, world."
// @Summary Hello world
// @Description trả về chuỗi kí tự "Hello World"
// @Produce json
// @Success 200 {object} service.SimpleMsgResponse
// @Router /todo [post]
func (todo *Todo) GetHello() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(SimpleMsgResponse{quote.HelloV3()})
	})
}

// ctxKey là kiểu dữ liệu được dùng để chỉ đến các giá trị được lưu trong request context,
// do bất cứ package nào cũng có thể sử dụng request context nên chúng ta phải có một kiểu
// dữ liệu riêng cho mỗi package để tránh việc các package ghi đè lên thông tin của nhau
type ctxKey string

// TaskIDCtxKey tên của giá trị nằm trong request context
var TaskIDCtxKey = ctxKey("TaskIDCtxKey")

// TaskIDCtx lấy id của task từ đường dẫn url
func TaskIDCtx() middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), TaskIDCtxKey, chi.URLParam(r, "id"))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SimpleMsgResponse chứa một tin nhắn đơn giản được trả về từ server
type SimpleMsgResponse struct {
	Message string `json:"message"`
}

// ErrResponse chứa chuỗi kí tự miêu tả lỗi đã xảy ra
type ErrResponse struct {
	Err string `json:"error"`
}
