package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

type TodoList struct {
	Todos []Todo `json:"todos"`
}

type FilterType int

const (
	FilterAll FilterType = iota
	FilterActive
	FilterCompleted
)

type AppState struct {
	Todos         []Todo
	CurrentFilter FilterType
	EditingID     string
	InputText     string
}

func NewTodo(text string) Todo {
	return Todo{
		ID:        uuid.New().String(),
		Text:      text,
		Completed: false,
		CreatedAt: time.Now(),
	}
}

func (t *Todo) Toggle() {
	t.Completed = !t.Completed
}

func (t *Todo) SetText(text string) {
	t.Text = text
}

func (tl *TodoList) AddTodo(text string) {
	todo := NewTodo(text)
	tl.Todos = append(tl.Todos, todo)
}

func (tl *TodoList) DeleteTodo(id string) bool {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			tl.Todos = append(tl.Todos[:i], tl.Todos[i+1:]...)
			return true
		}
	}
	return false
}

func (tl *TodoList) FindTodo(id string) *Todo {
	for i, todo := range tl.Todos {
		if todo.ID == id {
			return &tl.Todos[i]
		}
	}
	return nil
}

func (tl *TodoList) GetFilteredTodos(filter FilterType) []Todo {
	switch filter {
	case FilterActive:
		var active []Todo
		for _, todo := range tl.Todos {
			if !todo.Completed {
				active = append(active, todo)
			}
		}
		return active
	case FilterCompleted:
		var completed []Todo
		for _, todo := range tl.Todos {
			if todo.Completed {
				completed = append(completed, todo)
			}
		}
		return completed
	default:
		return tl.Todos
	}
}

type ErrorType int

const (
	ErrorValidation ErrorType = iota
	ErrorStorage
	ErrorUI
)

type AppError struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}