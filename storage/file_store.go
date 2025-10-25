package storage

import (
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/Brownie44l1/cli-todo-list/todo"
)

type FileStore struct {
	filepath string
	mu       sync.RWMutex
	todos     map[int]*todo.Task
	nextID   int
}

func NewFileStore(path string) (*FileStore, error) {
	fs := &FileStore{
		filepath: path,
		todos: make(map[int]*todo.Task),
		nextID: 1,
	}

	if err := fs.Load(); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	return fs, nil
}

func (fs *FileStore) Load() error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	f, err := os.Open(fs.filepath) 
	if err != nil {
		return err
	}
	defer f.Close()

	var list []*todo.Task
	dec := json.NewDecoder(f)
	if err := dec.Decode(&list); err != nil && err != io.EOF {
		return err
	}

	maxID := 0
	for _, t := range list {
		fs.todos[t.Id] = t 
		if t.Id > maxID {
			maxID = t.Id
		}
	}

	fs.nextID = maxID + 1
	if fs.nextID == 0 {
		fs.nextID = 1
	}
	
	return nil
}

func (fs *FileStore) Save() error {
	fs.mu.RLock()
    list := make([]*todo.Task, 0, len(fs.todos))
    for _, t := range fs.todos {
        copy := *t
        list = append(list, &copy)
    }
    fs.mu.RUnlock()

	tmpPath := fs.filepath + ".tmp"
	f, err := os.Create(tmpPath)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "	")
	if err := enc.Encode(list); err != nil {
		f.Close()
		os.Remove(tmpPath)
		return err
	}
	if err := f.Close(); err != nil {
		os.Remove(tmpPath)
		return err
	}
	
	if err := os.Rename(tmpPath, fs.filepath); err != nil {
		os.Remove(tmpPath)
		return err
	}
	return nil
}

func (fs *FileStore) Add(task todo.Task) error {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	id := fs.nextID
	fs.nextID++
	fs.todos[id] = &todo.Task{Id: id, Title: title, Description: descreiption, }
}