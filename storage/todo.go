package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// TODO: giới thiệu go test

// Metadata chứa các thông tin được chèn thêm để mô tả các thông tin được chứa trong cơ sỏ dữ liệu.
type Metadata struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Task miêu tả thông tin của một tác vụ.
type Task struct {
	Metadata
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Done     bool   `json:"done"`
}

// txFunc định nghĩa một hàm số nhận một transaction để sử dụng,
// và trả về lỗi có thể xảy ra trong quá trình chạy hàm số đó.
type txFunc func(*sql.Tx) error

// Todo được dùng để đóng gói (encapsulate) các câu lệnh dùng để truy cập cơ sở dữ liệu,
// thay vì sử dụng các câu lệnh SQL trực tiếp tại các hàm số xử lý đường dẫn URL.
// service.Todo có thể định nghĩa một interface bao gồm các methods của storage.Todo,
// và sử dụng storage.Todo thông qua interface đó, thay vì sử dụng storage.Todo trục tiếp.
type Todo struct {
	db *sql.DB
}

// NewTodo tạo một object mới dùng để quản lý cơ sở dữ liệu cho các tác vụ.
func NewTodo(db *sql.DB) *Todo {
	return &Todo{
		db: db,
	}
}

// CreateOne tạo một tác vụ mới và thêm nó vào cơ sở dữ liệu.
// Thông tin của tác vụ được thêm vào sẽ được trả về kèm với lỗi có thể xảy ra.
func (todo *Todo) CreateOne(ctx context.Context, task Task) (*Task, error) {
	txErr := todo.makeTx(ctx, func(tx *sql.Tx) error {
		res, err := tx.Exec(
			"INSERT INTO tasks(title, subtitle, done) VALUES (?, ?, ?)",
			task.Title, task.Subtitle, task.Done,
		)

		if err != nil {
			return err
		}

		// kiểu tra xem lệnh INSERT có thực hiện đúng không
		if rowsAffected, err := res.RowsAffected(); err != nil {
			return err
		} else if rowsAffected != 1 {
			return errors.New("number of affected rows is not 1")
		}

		// lấy id được thêm vào bởi cơ sở dữ liệu
		task.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}

		// lấy các thông tin khác được thêm vào bởi cơ sở dữ liệu
		row := tx.QueryRow("SELECT created_at, updated_at FROM tasks WHERE id = ?", task.ID)
		return row.Scan(&task.CreatedAt, &task.UpdatedAt)
	})

	if txErr != nil {
		return nil, txErr
	}

	return &task, nil
}

func (todo *Todo) makeTx(ctx context.Context, f txFunc) error {
	tx, err := todo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := f(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Wrapf(rollbackErr, "could not rollback from error: %v", err)
		}
		return err
	}

	return tx.Commit()
}
