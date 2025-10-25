package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Brownie44l1/cli-todo-list/storage"
	"github.com/Brownie44l1/cli-todo-list/todo"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	store, err := storage.NewSQLiteStore("todo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()

	manager := todo.NewManager(store)

	fmt.Println("=== Todo List CLI ===")
	fmt.Println("Commands: add, list, get, complete, delete, help, exit")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := parseCommand(line)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "add":
			if len(parts) < 2 {
				fmt.Println("Usage: add <title> [description]")
				fmt.Println("Example: add \"Buy groceries\" \"Milk, eggs, bread\"")
				continue
			}
			
			title := parts[1]
			description := ""
			if len(parts) >= 3 {
				description = parts[2]
			}

			err := manager.AddTask(title, description)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("✓ Task added successfully")
			}

		case "list":
			err := manager.List()
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <id>")
				continue
			}
			
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid task ID")
				continue
			}

			task, err := manager.GetTask(id)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				status := "Incomplete"
				if task.Completed {
					status = "Complete"
				}
				fmt.Printf("\nTask #%d\n", task.Id)
				fmt.Printf("Title: %s\n", task.Title)
				fmt.Printf("Description: %s\n", task.Description)
				fmt.Printf("Status: %s\n", status)
				fmt.Printf("Created: %s\n\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
			}

		case "complete":
			if len(parts) < 2 {
				fmt.Println("Usage: complete <id>")
				continue
			}
			
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid task ID")
				continue
			}

			err = manager.Complete(id)
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <id>")
				continue
			}
			
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				fmt.Println("Error: invalid task ID")
				continue
			}

			err = manager.Delete(id)
			if err != nil {
				fmt.Println("Error:", err)
			}

		case "help":
			printHelp()

		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		default:
			fmt.Printf("Unknown command: %s\n", parts[0])
			fmt.Println("Type 'help' for available commands")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// parseCommand splits a command line into parts, handling quoted strings
func parseCommand(line string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	
	for i := 0; i < len(line); i++ {
		char := line[i]
		
		switch char {
		case '"':
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current.WriteByte(char)
			} else if current.Len() > 0 {
				parts = append(parts, current.String())
				current.Reset()
			}
		default:
			current.WriteByte(char)
		}
	}
	
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}
	
	return parts
}

func printHelp() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  add <title> [description]  - Add a new task")
	fmt.Println("                              Example: add \"Buy groceries\" \"Milk and bread\"")
	fmt.Println("  list                       - List all tasks")
	fmt.Println("  get <id>                   - Display a specific task")
	fmt.Println("  complete <id>              - Mark a task as completed")
	fmt.Println("  delete <id>                - Delete a task")
	fmt.Println("  help                       - Show this help message")
	fmt.Println("  exit                       - Exit the application")
	fmt.Println()
}