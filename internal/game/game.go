package game

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"github.com/lapis2411/todo/internal/models"
	"github.com/lapis2411/todo/internal/storage"
	"github.com/lapis2411/todo/internal/ui"
)

const (
	WindowWidth  = 800
	WindowHeight = 600
	HeaderHeight = 80
	FooterHeight = 60
)

type Game struct {
	todos         models.TodoList
	currentFilter models.FilterType
	storage       storage.Storage
	uiManager     *UIManager
	error         string
}

type UIManager struct {
	inputBox      *ui.TextBox
	addButton     *ui.Button
	filterButtons map[models.FilterType]*ui.Button
	todoItems     []*ui.TodoItem
	scrollOffset  int
	windowWidth   int
	windowHeight  int
}

func NewGame(storagePath string) (*Game, error) {
	fileStorage := storage.NewFileStorage(storagePath)
	
	game := &Game{
		todos:         models.TodoList{Todos: []models.Todo{}},
		currentFilter: models.FilterAll,
		storage:       fileStorage,
	}

	// Load existing todos
	todos, err := game.storage.LoadTodos()
	if err != nil {
		game.error = fmt.Sprintf("Failed to load todos: %v", err)
	} else {
		game.todos.Todos = todos
	}

	game.uiManager = game.createUIManager()
	game.updateTodoItems()

	return game, nil
}

func (g *Game) createUIManager() *UIManager {
	uiMgr := &UIManager{
		filterButtons: make(map[models.FilterType]*ui.Button),
		windowWidth:   WindowWidth,
		windowHeight:  WindowHeight,
	}

	// Create input textbox
	uiMgr.inputBox = ui.NewTextBox(20, 20, 500, 35, "Add a new todo...")

	// Create add button
	uiMgr.addButton = ui.NewButton(540, 20, 100, 35, "Add", func() {
		g.addTodo()
	})

	// Create filter buttons
	filterLabels := map[models.FilterType]string{
		models.FilterAll:       "All",
		models.FilterActive:    "Active",
		models.FilterCompleted: "Completed",
	}

	buttonWidth := 80
	startX := 20
	for filter, label := range filterLabels {
		button := ui.NewButton(
			startX+int(filter)*buttonWidth, WindowHeight-FooterHeight+20, 
			buttonWidth-5, 25,
			label,
			func(f models.FilterType) func() {
				return func() { g.setFilter(f) }
			}(filter),
		)
		
		// Set initial colors
		if filter == g.currentFilter {
			button.SetColors(
				color.RGBA{0, 123, 255, 255},   // Active blue
				color.RGBA{0, 86, 179, 255},    // Darker blue
				color.RGBA{255, 255, 255, 255}, // White text
			)
		} else {
			button.SetColors(
				color.RGBA{108, 117, 125, 255}, // Gray
				color.RGBA{90, 98, 104, 255},    // Darker gray
				color.RGBA{255, 255, 255, 255}, // White text
			)
		}
		
		uiMgr.filterButtons[filter] = button
	}

	return uiMgr
}

func (g *Game) addTodo() {
	text := strings.TrimSpace(g.uiManager.inputBox.GetText())
	if text == "" {
		g.error = "Todo text cannot be empty"
		return
	}

	g.todos.AddTodo(text)
	g.uiManager.inputBox.Clear()
	g.error = ""
	
	if err := g.saveTodos(); err != nil {
		g.error = fmt.Sprintf("Failed to save: %v", err)
	}
	
	g.updateTodoItems()
}

func (g *Game) deleteTodo(id string) {
	if g.todos.DeleteTodo(id) {
		g.error = ""
		if err := g.saveTodos(); err != nil {
			g.error = fmt.Sprintf("Failed to save: %v", err)
		}
		g.updateTodoItems()
	}
}

func (g *Game) toggleTodo(id string) {
	if todo := g.todos.FindTodo(id); todo != nil {
		todo.Toggle()
		g.error = ""
		if err := g.saveTodos(); err != nil {
			g.error = fmt.Sprintf("Failed to save: %v", err)
		}
		g.updateTodoItems()
	}
}

func (g *Game) editTodo(id, newText string) {
	if newText == "" {
		g.error = "Todo text cannot be empty"
		return
	}
	
	if todo := g.todos.FindTodo(id); todo != nil {
		todo.SetText(newText)
		g.error = ""
		if err := g.saveTodos(); err != nil {
			g.error = fmt.Sprintf("Failed to save: %v", err)
		}
		g.updateTodoItems()
	}
}

func (g *Game) setFilter(filter models.FilterType) {
	g.currentFilter = filter
	g.updateFilterButtons()
	g.updateTodoItems()
}

func (g *Game) updateFilterButtons() {
	for filter, button := range g.uiManager.filterButtons {
		if filter == g.currentFilter {
			button.SetColors(
				color.RGBA{0, 123, 255, 255},   // Active blue
				color.RGBA{0, 86, 179, 255},    // Darker blue
				color.RGBA{255, 255, 255, 255}, // White text
			)
		} else {
			button.SetColors(
				color.RGBA{108, 117, 125, 255}, // Gray
				color.RGBA{90, 98, 104, 255},    // Darker gray
				color.RGBA{255, 255, 255, 255}, // White text
			)
		}
	}
}

func (g *Game) updateTodoItems() {
	filteredTodos := g.todos.GetFilteredTodos(g.currentFilter)
	g.uiManager.todoItems = make([]*ui.TodoItem, 0, len(filteredTodos))

	itemHeight := 50
	startY := HeaderHeight + 10
	
	for i, todo := range filteredTodos {
		y := startY + i*itemHeight - g.uiManager.scrollOffset
		
		todoItem := ui.NewTodoItem(&filteredTodos[i], 20, y, g.uiManager.windowWidth-40, itemHeight)
		
		// Setup delete button callback
		todoItem.GetDeleteButton().OnClick = func(todoID string) func() {
			return func() {
				g.deleteTodo(todoID)
			}
		}(todo.ID)
		
		g.uiManager.todoItems = append(g.uiManager.todoItems, todoItem)
	}
}

func (g *Game) saveTodos() error {
	return g.storage.SaveTodos(g.todos.Todos)
}

func (g *Game) Update() error {
	// Handle adding todo with Enter key
	if g.uiManager.inputBox.IsEnterPressed() {
		g.addTodo()
	}

	// Update UI components
	g.uiManager.inputBox.Update()
	g.uiManager.addButton.Update()
	
	for _, button := range g.uiManager.filterButtons {
		button.Update()
	}

	// Update todo items
	for _, item := range g.uiManager.todoItems {
		item.Update()
	}

	// Handle scroll (simple implementation)
	if len(g.uiManager.todoItems) > 0 {
		_, dy := ebiten.Wheel()
		g.uiManager.scrollOffset -= int(dy * 20)
		if g.uiManager.scrollOffset < 0 {
			g.uiManager.scrollOffset = 0
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Clear screen with background color
	screen.Fill(color.RGBA{248, 249, 250, 255})

	// Draw header
	g.drawHeader(screen)
	
	// Draw main content area
	g.drawContent(screen)
	
	// Draw footer
	g.drawFooter(screen)

	// Draw error message if any
	if g.error != "" {
		g.drawError(screen)
	}
}

func (g *Game) drawHeader(screen *ebiten.Image) {
	// Draw header background
	headerColor := color.RGBA{255, 255, 255, 255}
	ebitenutil.DrawRect(screen, 0, 0, float64(g.uiManager.windowWidth), HeaderHeight, headerColor)

	// Draw title
	title := "Todo List"
	titleBounds := text.BoundString(basicfont.Face7x13, title)
	titleX := (g.uiManager.windowWidth - (titleBounds.Max.X - titleBounds.Min.X)) / 2
	titleY := 15
	text.Draw(screen, title, basicfont.Face7x13, titleX, titleY, color.RGBA{33, 37, 41, 255})

	// Draw input and add button
	g.uiManager.inputBox.Draw(screen)
	g.uiManager.addButton.Draw(screen)

	// Draw header border
	borderColor := color.RGBA{200, 200, 200, 255}
	ebitenutil.DrawRect(screen, 0, HeaderHeight-1, float64(g.uiManager.windowWidth), 1, borderColor)
}

func (g *Game) drawContent(screen *ebiten.Image) {
	// Draw todo items
	for _, item := range g.uiManager.todoItems {
		if item.Y >= HeaderHeight && item.Y < g.uiManager.windowHeight-FooterHeight {
			item.Draw(screen)
		}
	}

	// Draw empty state message if no todos
	filteredTodos := g.todos.GetFilteredTodos(g.currentFilter)
	if len(filteredTodos) == 0 {
		message := "No todos yet. Add one above!"
		if g.currentFilter == models.FilterActive {
			message = "No active todos!"
		} else if g.currentFilter == models.FilterCompleted {
			message = "No completed todos!"
		}

		messageBounds := text.BoundString(basicfont.Face7x13, message)
		messageX := (g.uiManager.windowWidth - (messageBounds.Max.X - messageBounds.Min.X)) / 2
		messageY := HeaderHeight + 50
		text.Draw(screen, message, basicfont.Face7x13, messageX, messageY, color.RGBA{108, 117, 125, 255})
	}
}

func (g *Game) drawFooter(screen *ebiten.Image) {
	footerY := float64(g.uiManager.windowHeight - FooterHeight)
	
	// Draw footer background
	footerColor := color.RGBA{255, 255, 255, 255}
	ebitenutil.DrawRect(screen, 0, footerY, float64(g.uiManager.windowWidth), FooterHeight, footerColor)

	// Draw footer border
	borderColor := color.RGBA{200, 200, 200, 255}
	ebitenutil.DrawRect(screen, 0, footerY, float64(g.uiManager.windowWidth), 1, borderColor)

	// Draw filter buttons
	for _, button := range g.uiManager.filterButtons {
		button.Draw(screen)
	}

	// Draw todo count
	filteredCount := len(g.todos.GetFilteredTodos(g.currentFilter))
	totalCount := len(g.todos.Todos)
	countText := fmt.Sprintf("%d of %d todos", filteredCount, totalCount)
	
	countX := g.uiManager.windowWidth - 150
	countY := int(footerY) + 35
	text.Draw(screen, countText, basicfont.Face7x13, countX, countY, color.RGBA{108, 117, 125, 255})
}

func (g *Game) drawError(screen *ebiten.Image) {
	// Draw error background
	errorY := HeaderHeight + 10
	errorHeight := 30
	ebitenutil.DrawRect(screen, 20, float64(errorY), float64(g.uiManager.windowWidth-40), float64(errorHeight), color.RGBA{248, 215, 218, 255})

	// Draw error border
	borderColor := color.RGBA{220, 53, 69, 255}
	ebitenutil.DrawRect(screen, 20, float64(errorY), float64(g.uiManager.windowWidth-40), 1, borderColor)
	ebitenutil.DrawRect(screen, 20, float64(errorY+errorHeight-1), float64(g.uiManager.windowWidth-40), 1, borderColor)

	// Draw error text
	errorX := 30
	errorTextY := errorY + 20
	text.Draw(screen, g.error, basicfont.Face7x13, errorX, errorTextY, color.RGBA{114, 28, 36, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.uiManager.windowWidth = outsideWidth
	g.uiManager.windowHeight = outsideHeight
	return outsideWidth, outsideHeight
}