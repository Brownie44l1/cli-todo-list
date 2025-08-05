package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Todo struct {
	ID    int       `json:"id"`
	Title string    `json:"title"`
	Done  bool      `json:"done"`
	Added time.Time `json:"added"`
}

var todofile = "todo.json"

func saveTodos(todos []Todo) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal todos: %w", err)
	}

	return os.WriteFile(todofile, data, 0644)
}

func loadTodos() ([]Todo, error) {
	if _, err := os.Stat(todofile); os.IsNotExist(err) {
		return []Todo{}, nil
	}

	data, err := os.ReadFile(todofile)
	if err != nil {
		return nil, fmt.Errorf("failed to read Todo: %w", err)
	}

	var todos []Todo
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal todos: %w", err)
	}

	return todos, nil
}

func addTask() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("usage: %s add 'task description'", os.Args[0])
	}

	todos, err := loadTodos()
	if err != nil {
		return err
	}

	newID := 1
	if len(todos) > 0 {
		newID = todos[len(todos)-1].ID + 1
	}

	newTodo := Todo{
		ID:    newID,
		Title: strings.Join(os.Args[2:], " "),
		Done:  false,
		Added: time.Now(),
	}

	todos = append(todos, newTodo)
	return saveTodos(todos)
}

func deleteTask(id int) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			return saveTodos(todos)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func updateStatus(id int) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}

	for i := range todos {
		if todos[i].ID == id {
			todos[i].Done = true
			fmt.Printf("Marked task %s as done\n", todos[i].Title)
			return saveTodos(todos)
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}

func listTasks() error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}

	if len(todos) == 0 {
		fmt.Println("No tasks found.")
		return nil
	}

	fmt.Println("Your TODO list:")
	for _, todo := range todos {
		status := " "
		if todo.Done {
			status = "✓"
		}
		fmt.Printf("[%d] %s %s (added: %s)\n", todo.ID, status, todo.Title, todo.Added.Format("2006-01-02 15:04"))
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo <command> [arguments]")
		fmt.Println("Commands:")
		fmt.Println("  add <task> - Add a new task")
		fmt.Println("  list - List all tasks")
		fmt.Println("  done <id> - Mark task as done")
		fmt.Println("  delete <id> - Delete a task")
		os.Exit(1)
	}

	command := os.Args[1]
	var err error

	switch command {
	case "add":
		err = addTask()
	case "list":
		err = listTasks()
	case "done":
		if len(os.Args) < 3 {
			err = fmt.Errorf("usage: %s done <task-id>", os.Args[0])
			break
		}
		var id int
		_, err = fmt.Sscanf(os.Args[2], "%d", &id)
		if err != nil {
			err = fmt.Errorf("invalid task ID")
			break
		}
		err = updateStatus(id)
	case "delete":
		if len(os.Args) < 3 {
			err = fmt.Errorf("usage: %s delete <task-id>", os.Args[0])
			break
		}
		var id int
		_, err = fmt.Sscanf(os.Args[2], "%d", &id)
		if err != nil {
			err = fmt.Errorf("invalid task ID")
			break
		}
		err = deleteTask(id)
	default:
		err = fmt.Errorf("unknown command: %s", command)
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}