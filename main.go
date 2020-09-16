package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/letung3105/todos-svr/service"
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

	svr := http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      service.NewTodo(),
		TLSConfig:    nil,
		WriteTimeout: time.Second, // thời gian tối đa được dùng để ghi gói tin
		ReadTimeout:  time.Second, // thời gian tối đa được dùng để đọc gói tin
		IdleTimeout:  time.Second, // thời gian tối đa kết nối được phép ngưng hoạt động
	}

	// TODO: graceful shutdown
	if err := svr.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err)
	}
}
