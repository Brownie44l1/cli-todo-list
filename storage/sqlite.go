package storage

import (
	"database/sql"
	"time"

	"github.com/Brownie44l1/cli-todo-list/todo"
	_ "github.com/mattn/go-sqlite3"
)


type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todo_list (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			description TEXT,
			completed BOOLEAN,
			created_at TEXT
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(task todo.Task) error {
	_, err := s.db.Exec(
		"INSERT INTO todo_list (title, description, completed , created_at) VALUES (?,?,?,?)", task.Title, task.Description, task.Completed, task.CreatedAt.Format(time.RFC3339))
	return err
}

func (s *SQLiteStore) Get(id int) (*todo.Task, error) {
	var task todo.Task
	var createdAtStr string
	err := s.db.QueryRow("SELECT id, title, description, completed, created_at FROM todo_list WHERE id=?", id).
		Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &createdAtStr)
	if err != nil {
		return nil, err
	}
	task.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
	return &task, nil
}

func (s *SQLiteStore) List() ([]todo.Task, error) {
	rows, err := s.db.Query("SELECT id, title, description, completed, created_at FROM todo_list")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []todo.Task
	for rows.Next() {
		var task todo.Task
		var createdAtStr string
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &createdAtStr)
		if err != nil {
			return nil, err
		}
		task.CreatedAt, _ = time.Parse(time.RFC3339, createdAtStr)
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *SQLiteStore) Complete(id int) error {
	_, err := s.db.Exec("UPDATE todo_list SET completed=1 WHERE id=?", id)
	return err
}

func (s *SQLiteStore) Delete(id int) error {
	_, err := s.db.	Exec("DELETE FROM todo_list WHERE id=?", id)
	return err
}