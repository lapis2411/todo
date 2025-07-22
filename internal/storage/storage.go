package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/lapis2411/todo/internal/models"
)

type Storage interface {
	SaveTodos(todos []models.Todo) error
	LoadTodos() ([]models.Todo, error)
	ClearTodos() error
}

type FileStorage struct {
	filepath string
}

func NewFileStorage(filepath string) *FileStorage {
	return &FileStorage{
		filepath: filepath,
	}
}

func (fs *FileStorage) SaveTodos(todos []models.Todo) error {
	todoList := models.TodoList{Todos: todos}
	
	data, err := json.MarshalIndent(todoList, "", "  ")
	if err != nil {
		return &models.AppError{
			Type:    models.ErrorStorage,
			Message: "Failed to marshal todos to JSON",
			Err:     err,
		}
	}

	dir := filepath.Dir(fs.filepath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &models.AppError{
			Type:    models.ErrorStorage,
			Message: "Failed to create data directory",
			Err:     err,
		}
	}

	if err := os.WriteFile(fs.filepath, data, 0644); err != nil {
		return &models.AppError{
			Type:    models.ErrorStorage,
			Message: "Failed to write todos to file",
			Err:     err,
		}
	}

	return nil
}

func (fs *FileStorage) LoadTodos() ([]models.Todo, error) {
	if _, err := os.Stat(fs.filepath); os.IsNotExist(err) {
		return []models.Todo{}, nil
	}

	data, err := os.ReadFile(fs.filepath)
	if err != nil {
		return nil, &models.AppError{
			Type:    models.ErrorStorage,
			Message: "Failed to read todos file",
			Err:     err,
		}
	}

	if len(data) == 0 {
		return []models.Todo{}, nil
	}

	var todoList models.TodoList
	if err := json.Unmarshal(data, &todoList); err != nil {
		return nil, &models.AppError{
			Type:    models.ErrorStorage,
			Message: "Failed to parse todos from JSON",
			Err:     err,
		}
	}

	return todoList.Todos, nil
}

func (fs *FileStorage) ClearTodos() error {
	if _, err := os.Stat(fs.filepath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(fs.filepath); err != nil {
		return &models.AppError{
			Type:    models.ErrorStorage,
			Message: "Failed to remove todos file",
			Err:     err,
		}
	}

	return nil
}

func (fs *FileStorage) GetFilePath() string {
	return fs.filepath
}