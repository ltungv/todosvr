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
	// TODO: nhận thông tin từ giao diện dòng lệnh (command line interrface)
	host = ""
	port = 3000
)

func init() {
	flag.StringVar(&host, "host", "", "Host address to bind to")
	flag.IntVar(&port, "port", 3000, "Port number to bind to")
}

func main() {
	flag.Parse()

	// cài đặt đường dẫn cho ứng dụng (routes)
	r := chi.NewRouter()
	r.Use(
		middleware.JSONResponse(),
		middleware.LogRequest(log.New(os.Stdout, "", log.LstdFlags)),
	)
	r.Get("/hello", handler.GetHello())

	// tạo struct `http.Server`
	addr := fmt.Sprintf("%s:%d", host, port)
	svr := http.Server{
		Addr:         addr,
		Handler:      r,
		TLSConfig:    nil,
		WriteTimeout: time.Second, // thời gian tối đa được dùng để ghi gói tin
		ReadTimeout:  time.Second, // thời gian tối đa được dùng để đọc gói tin
		IdleTimeout:  time.Second, // thời gian tối đa kết nối được phép ngưng hoạt động
	}

	// TODO: graceful shutdown
	// bắt đầu chạy server
	log.Printf("Starting server on %s", addr)
	if err := svr.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
