## Go TODO CLI App

A simple and efficient command-line TODO application written in Go.  
Easily manage your tasks locally with support for both **SQLite** and **JSON file** storage backends.

---

### 🚀 Features

- [x] Add new tasks  
- [x] List all tasks (with creation timestamps)  
- [x] Mark tasks as completed (✓)  
- [x] Delete tasks by ID  
- [x] Persistent storage (choose between SQLite or JSON)  
- [x] Sequential task numbering — always compact (no gaps after deletion)  
- [x] Auto-loads tasks on startup and auto-saves on exit  

---

### 🧩 Storage Options

You can easily switch between backends in `cmd/main.go`:

```go
// Use JSON file storage
store, err := storage.NewFileStore("./data/todo.json")

// OR use SQLite storage
// store, err := storage.NewSQLiteStore("./data/todo.db")
