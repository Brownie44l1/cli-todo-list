package storage

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/Brownie44l1/cli-todo-list/todo"
)

type FileStore struct {
	filepath string
	//mu       sync.Mutex
	tasks    []*todo.Task
	nextID   int
}

func NewFileStore(path string) (*FileStore, error) {
	fs := &FileStore{
		filepath: path,
		tasks: []*todo.Task{},
		nextID: 1,
	}

	if err := fs.Load(); err != nil {
		return nil, err
	}

	return fs, nil
}

func (fs *FileStore) Load() error {
	data, err := os.ReadFile(fs.filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if len(data) == 0 {
		return nil
	}

	err = json.Unmarshal(data, &fs.tasks)
	if err != nil {
		return err
	}

	maxID := 0
	for _, t := range fs.tasks {
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
    data, err := json.MarshalIndent(fs.tasks, "", "	")

	if err != nil {
		return err
	}

	return os.WriteFile(fs.filepath, data, 0644)
}

func (fs *FileStore) Add(task todo.Task) error {
	task.Id = fs.nextID
	task.CreatedAt = time.Now()
	fs.nextID++
	fs.tasks = append(fs.tasks, &task)
	return fs.Save()
}

func (fs *FileStore) Get(id int) (*todo.Task, error) {
	for _, t := range fs.tasks {
		if t.Id == id {
			return t, nil
		}
	}
	return nil, errors.New("task not found")
}

func (fs *FileStore) List() ([]todo.Task, error) {
	tasks := make([]todo.Task, len(fs.tasks))
	for i, t := range fs.tasks {
		tasks[i] = *t
	}
	return tasks, nil
}

func (fs *FileStore) Complete(id int) error {
	for _, t := range fs.tasks {
		if t.Id == id {
			t.Completed = true
			return fs.Save()
		}
	}
	return errors.New("task not found")
}

func (fs *FileStore) Delete(id int) error {
	for i, t := range fs.tasks {
		if t.Id == id {
			fs.tasks = append(fs.tasks[:i], fs.tasks[i+1:]...)
			return fs.Save()
		}
	}
	return errors.New("task not found")
}

func (fs *FileStore) Close() error {
	// For file-based store, nothing to close, but we'll persist before exiting
	return fs.Save()
}