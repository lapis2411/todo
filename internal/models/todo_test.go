package models

import (
	"testing"
)

func TestNewTodo(t *testing.T) {
	text := "Test todo"
	todo := NewTodo(text)

	if todo.Text != text {
		t.Errorf("Expected text %s, got %s", text, todo.Text)
	}

	if todo.Completed {
		t.Error("New todo should not be completed")
	}

	if todo.ID == "" {
		t.Error("Todo should have an ID")
	}

	if todo.CreatedAt.IsZero() {
		t.Error("Todo should have a creation time")
	}
}

func TestTodoToggle(t *testing.T) {
	todo := NewTodo("Test")
	
	if todo.Completed {
		t.Error("Todo should start uncompleted")
	}

	todo.Toggle()
	if !todo.Completed {
		t.Error("Todo should be completed after toggle")
	}

	todo.Toggle()
	if todo.Completed {
		t.Error("Todo should be uncompleted after second toggle")
	}
}

func TestTodoListAddTodo(t *testing.T) {
	todoList := TodoList{}
	text := "Test todo"

	todoList.AddTodo(text)

	if len(todoList.Todos) != 1 {
		t.Errorf("Expected 1 todo, got %d", len(todoList.Todos))
	}

	if todoList.Todos[0].Text != text {
		t.Errorf("Expected text %s, got %s", text, todoList.Todos[0].Text)
	}
}

func TestTodoListDeleteTodo(t *testing.T) {
	todoList := TodoList{}
	todoList.AddTodo("Test 1")
	todoList.AddTodo("Test 2")

	id := todoList.Todos[0].ID
	deleted := todoList.DeleteTodo(id)

	if !deleted {
		t.Error("Expected delete to return true")
	}

	if len(todoList.Todos) != 1 {
		t.Errorf("Expected 1 todo remaining, got %d", len(todoList.Todos))
	}

	if todoList.Todos[0].Text != "Test 2" {
		t.Errorf("Expected remaining todo to be 'Test 2', got %s", todoList.Todos[0].Text)
	}
}

func TestTodoListFindTodo(t *testing.T) {
	todoList := TodoList{}
	todoList.AddTodo("Test 1")
	todoList.AddTodo("Test 2")

	id := todoList.Todos[0].ID
	found := todoList.FindTodo(id)

	if found == nil {
		t.Error("Expected to find todo")
	}

	if found.Text != "Test 1" {
		t.Errorf("Expected found todo text to be 'Test 1', got %s", found.Text)
	}

	notFound := todoList.FindTodo("nonexistent")
	if notFound != nil {
		t.Error("Expected not to find nonexistent todo")
	}
}

func TestGetFilteredTodos(t *testing.T) {
	todoList := TodoList{}
	todoList.AddTodo("Active todo")
	todoList.AddTodo("Completed todo")
	
	// Mark second todo as completed
	todoList.Todos[1].Toggle()

	// Test FilterAll
	allTodos := todoList.GetFilteredTodos(FilterAll)
	if len(allTodos) != 2 {
		t.Errorf("FilterAll: expected 2 todos, got %d", len(allTodos))
	}

	// Test FilterActive
	activeTodos := todoList.GetFilteredTodos(FilterActive)
	if len(activeTodos) != 1 {
		t.Errorf("FilterActive: expected 1 todo, got %d", len(activeTodos))
	}
	if activeTodos[0].Text != "Active todo" {
		t.Errorf("FilterActive: expected 'Active todo', got %s", activeTodos[0].Text)
	}

	// Test FilterCompleted
	completedTodos := todoList.GetFilteredTodos(FilterCompleted)
	if len(completedTodos) != 1 {
		t.Errorf("FilterCompleted: expected 1 todo, got %d", len(completedTodos))
	}
	if completedTodos[0].Text != "Completed todo" {
		t.Errorf("FilterCompleted: expected 'Completed todo', got %s", completedTodos[0].Text)
	}
}