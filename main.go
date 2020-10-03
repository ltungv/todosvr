package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-sql-driver/mysql"
	"github.com/letung3105/todosvr/docs"
	"github.com/letung3105/todosvr/service"
	"github.com/letung3105/todosvr/storage"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	host     = flag.String("host", "localhost:3000", "Host address to bind to")
	dsnMySQL = flag.String("dsn_mysql", "@/todosvr", "The connection string of the database")
)

func init() {
}

// @contact.name Tung L. Vo
// @contact.url https://github.com/letung3105
// @contact.email letung3105@gmail.com
// @license.name MIT
// @license.url https://mit-license.org

func main() {
	flag.Parse()

	// Swagger documentations
	docs.SwaggerInfo.Title = "Todo API"
	docs.SwaggerInfo.Description = "Dịch vụ API đơn giản"
	docs.SwaggerInfo.Version = "0.1.0"
	docs.SwaggerInfo.Host = *host
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	// create an MySQL database object
	db, err := initMySQL(*dsnMySQL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Get("/swagger/*", httpSwagger.Handler(
		// doc.json được tạo từ package `docs`
		httpSwagger.URL("http://localhost:3000/swagger/doc.json"),
	))
	r.Mount("/", service.NewTodo(storage.NewTodo(db)))

	svr := http.Server{
		Addr:         *host,
		Handler:      r,
		TLSConfig:    nil,
		WriteTimeout: time.Second, // thời gian tối đa được dùng để ghi gói tin
		ReadTimeout:  time.Second, // thời gian tối đa được dùng để đọc gói tin
		IdleTimeout:  time.Second, // thời gian tối đa kết nối được phép ngưng hoạt động
	}

	// TODO: graceful shutdown
	log.Printf("Listen and server on %s", *host)
	if err := svr.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(err)
	}
}

func initMySQL(dsn string) (*sql.DB, error) {
	// TODO: cài đặt cơ sở dữ liệu thông qua cấu hình tuỳ chỉnh
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	cfg.ParseTime = true // NOTE: phải cài đặt để có thể Scan các biến có kiểu dữ liệu time.Time
	cfg.ReadTimeout = 200 * time.Millisecond
	cfg.WriteTimeout = 300 * time.Millisecond

	return sql.Open("mysql", cfg.FormatDSN())
}
