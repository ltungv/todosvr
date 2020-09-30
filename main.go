package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/letung3105/todos-svr/service"
	"github.com/letung3105/todos-svr/storage"
)

var (
	host = ""
	port = 3000
	dsn  = ""
)

func init() {
	flag.StringVar(&dsn, "dsn", "@/todosvr", "The connection string of the database")
	flag.StringVar(&host, "host", "", "Host address to bind to")
	flag.IntVar(&port, "port", 3000, "Port number to bind to")
}

func main() {
	flag.Parse()

	// TODO: cài đặt cơ sở dữ liệu thông qua cấu hình tuỳ chỉnh
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		log.Fatal(err)
	}

	cfg.ParseTime = true
	// ^ NOTE: phải cài đặt để có thể Scan các biến có kiểu dữ liệu time.Time
	cfg.ReadTimeout = 200 * time.Millisecond
	cfg.WriteTimeout = 300 * time.Millisecond

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	svr := http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      service.NewTodo(storage.NewTodo(db)),
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
