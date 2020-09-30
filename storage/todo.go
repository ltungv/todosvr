package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"
)

// TODO: giới thiệu go test

// Metadata chứa các thông tin phụ trợ.
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

// CreateOne thêm một tác vụ mới vào cơ sở dữ liệu và trả về tác vụ được thêm vào (kèm với thông tin đã được cập nhật).
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

// ReadAll trả về tất cả các tác vụ có trong cơ sở dữ liệu.
// TODO: giới hạn số lượng tác vụ được trả về
func (todo *Todo) ReadAll(ctx context.Context) ([]*Task, error) {
	panic("NOT IMPLEMENTED")
}

// UpdateOne cập nhật thông tin của một tác vụ có trong cơ sở dữ liệu và trả về tác vụ đó (kèm với thông tin đã được cập nhật).
func (todo *Todo) UpdateOne(ctx context.Context, id int64, task Task) (*Task, error) {
	panic("NOT IMPLEMENTED")
}

// DeleteOne loại bỏ một tác vụ khỏi cơ sở dữ liệu.
func (todo *Todo) DeleteOne(ctx context.Context, id int64) error {
	panic("NOT IMPLEMENTED")
}

// makeTx gíup đơn giản hoá việt sử dụng transaction. Hàm số này sẽ tạo một transaction,
// sử dụng function được truyền vào lên trasaction đó rồi kiểm tra lỗi có thể xảy ra.
// Với hàm số này, bạn không cần phải tự gọi `Rollback` mỗi khi có lỗi.
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
