package storage

import (
	"os"
	"path/filepath"
	"testing"

	"todo-app/internal/models"
)

func TestFileStorageSaveAndLoad(t *testing.T) {
	// Create temporary file for testing
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_todos.json")

	storage := NewFileStorage(testFile)

	// Create test todos
	todos := []models.Todo{
		models.NewTodo("Test todo 1"),
		models.NewTodo("Test todo 2"),
	}
	todos[1].Toggle() // Mark second as completed

	// Save todos
	err := storage.SaveTodos(todos)
	if err != nil {
		t.Fatalf("Failed to save todos: %v", err)
	}

	// Check file exists
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("Todos file was not created")
	}

	// Load todos
	loadedTodos, err := storage.LoadTodos()
	if err != nil {
		t.Fatalf("Failed to load todos: %v", err)
	}

	// Verify loaded todos
	if len(loadedTodos) != 2 {
		t.Errorf("Expected 2 todos, got %d", len(loadedTodos))
	}

	if loadedTodos[0].Text != "Test todo 1" {
		t.Errorf("Expected first todo text 'Test todo 1', got %s", loadedTodos[0].Text)
	}

	if !loadedTodos[1].Completed {
		t.Error("Expected second todo to be completed")
	}
}

func TestFileStorageLoadNonexistentFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "nonexistent.json")

	storage := NewFileStorage(testFile)

	// Load from nonexistent file should return empty slice
	todos, err := storage.LoadTodos()
	if err != nil {
		t.Errorf("Loading nonexistent file should not error: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos from nonexistent file, got %d", len(todos))
	}
}

func TestFileStorageLoadEmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.json")

	// Create empty file
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}
	file.Close()

	storage := NewFileStorage(testFile)

	// Load from empty file should return empty slice
	todos, err := storage.LoadTodos()
	if err != nil {
		t.Errorf("Loading empty file should not error: %v", err)
	}

	if len(todos) != 0 {
		t.Errorf("Expected 0 todos from empty file, got %d", len(todos))
	}
}

func TestFileStorageClearTodos(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "clear_test.json")

	storage := NewFileStorage(testFile)

	// Save some todos first
	todos := []models.Todo{
		models.NewTodo("Test todo"),
	}
	err := storage.SaveTodos(todos)
	if err != nil {
		t.Fatalf("Failed to save todos: %v", err)
	}

	// Clear todos
	err = storage.ClearTodos()
	if err != nil {
		t.Errorf("Failed to clear todos: %v", err)
	}

	// Check file no longer exists
	if _, err := os.Stat(testFile); !os.IsNotExist(err) {
		t.Error("Todos file should have been deleted")
	}

	// Clear nonexistent file should not error
	err = storage.ClearTodos()
	if err != nil {
		t.Errorf("Clearing nonexistent file should not error: %v", err)
	}
}