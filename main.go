package main

import (
	"log"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/lapis2411/todo/internal/game"
)

const (
	WindowTitle  = "Todo List App"
	WindowWidth  = 800
	WindowHeight = 600
	DataFile     = "data/todos.json"
)

func main() {
	// Set window properties
	ebiten.SetWindowSize(WindowWidth, WindowHeight)
	ebiten.SetWindowTitle(WindowTitle)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// Get absolute path for data file
	dataPath, err := filepath.Abs(DataFile)
	if err != nil {
		log.Fatalf("Failed to get absolute path for data file: %v", err)
	}

	// Create and initialize game
	g, err := game.NewGame(dataPath)
	if err != nil {
		log.Fatalf("Failed to create game: %v", err)
	}

	// Run game
	if err := ebiten.RunGame(g); err != nil {
		log.Fatalf("Game failed: %v", err)
	}
}