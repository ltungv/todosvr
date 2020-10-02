package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
	"github.com/letung3105/todos-svr/docs"
	"github.com/letung3105/todos-svr/service"
	"github.com/letung3105/todos-svr/storage"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	host = ""
	port = 3000
	dsn  = ""
)

func init() {
	flag.StringVar(&dsn, "dsn", "@/todosvr", "The connection string of the database")
	flag.StringVar(&host, "host", "localhost", "Host address to bind to")
	flag.IntVar(&port, "port", 3000, "Port number to bind to")
}

// @contact.name Tung L. Vo
// @contact.url https://github.com/letung3105
// @contact.email letung3105@gmail.com
// @license.name MIT
// @license.url https://mit-license.org

func main() {
	flag.Parse()

	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "Dịch vụ API đơn giản"
	docs.SwaggerInfo.Version = "0.1.0"
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", host, port)
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

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

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
		// ^ NOTE: đường dẫn tới tệp tin JSON chứa định nghĩa của API,
		// thông tin của tệp tin này được tạo bởi package docs
	))
	r.Mount("/", service.NewTodo(storage.NewTodo(db)))

	svr := http.Server{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Handler:      r,
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
