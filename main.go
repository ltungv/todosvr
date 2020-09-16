package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/letung3105/todos-svr/handler"
	"github.com/letung3105/todos-svr/middleware"
)

var (
	host = ""
	port = 3000
)

func init() {
	flag.StringVar(&host, "host", "", "Host address to bind to")
	flag.IntVar(&port, "port", 3000, "Port number to bind to")
}

func main() {
	flag.Parse()

	// TODO: graceful shutdown (thêm hàm số vào cấu trúc `Service` để ngưng server)
	s := Service{}
	if err := s.Serve(fmt.Sprintf("%s:%d", host, port)); err != nil {
		log.Println(err)
	}
}

// Service chứa các hàm số dùng để khởi tạo và chạy dịch vụ API
type Service struct{}

// Serve khởi tạo và chạy server trên địa chỉ được nhận
func (s *Service) Serve(addr string) error {
	svr := http.Server{
		Addr:         addr,
		Handler:      s.routes(),
		TLSConfig:    nil,
		WriteTimeout: time.Second, // thời gian tối đa được dùng để ghi gói tin
		ReadTimeout:  time.Second, // thời gian tối đa được dùng để đọc gói tin
		IdleTimeout:  time.Second, // thời gian tối đa kết nối được phép ngưng hoạt động
	}

	log.Printf("Starting server on %s", addr)
	return svr.ListenAndServe()
}

// routes cài đặt các đường dẫn HTTP được hỗ trợ bởi dịch vụ
func (s *Service) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(
		middleware.JSONResponse(),
		middleware.LogRequest(log.New(os.Stdout, "", log.LstdFlags)),
	)
	r.Get("/hello", handler.GetHello())
	return r
}
