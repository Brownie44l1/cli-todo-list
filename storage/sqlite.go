package storage

import (
	"database/sql"
	"fmt"
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
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS todo_list (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT 0,
			created_at TEXT NOT NULL
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Add(task todo.Task) error {
	_, err := s.db.Exec(
		"INSERT INTO todo_list (title, description, completed, created_at) VALUES (?, ?, ?, ?)",
		task.Title,
		task.Description,
		task.Completed,
		task.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to add task: %w", err)
	}
	return nil
}

func (s *SQLiteStore) Get(id int) (*todo.Task, error) {
	var task todo.Task
	var createdAtStr string
	
	err := s.db.QueryRow(
		"SELECT id, title, description, completed, created_at FROM todo_list WHERE id = ?",
		id,
	).Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &createdAtStr)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task with id %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	
	task.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created_at: %w", err)
	}
	
	return &task, nil
}

func (s *SQLiteStore) List() ([]todo.Task, error) {
	rows, err := s.db.Query(
		"SELECT id, title, description, completed, created_at FROM todo_list ORDER BY id ASC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	var tasks []todo.Task
	for rows.Next() {
		var task todo.Task
		var createdAtStr string
		
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Completed, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		
		task.CreatedAt, err = time.Parse(time.RFC3339, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}
		
		tasks = append(tasks, task)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return tasks, nil
}

func (s *SQLiteStore) Complete(id int) error {
	result, err := s.db.Exec("UPDATE todo_list SET completed = 1 WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to complete task: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	
	return nil
}

func (s *SQLiteStore) Delete(id int) error {
	result, err := s.db.Exec("DELETE FROM todo_list WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rows == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}
	
	return nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}