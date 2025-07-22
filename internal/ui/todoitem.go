package ui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"

	"todo-app/internal/models"
)

type TodoItem struct {
	Todo          *models.Todo
	X, Y          int
	Width         int
	Height        int
	Editing       bool
	EditText      string
	EditTextBox   *TextBox
	Checkbox      *Button
	DeleteBtn     *Button
	lastClickTime time.Time
	Hovered       bool
}

func NewTodoItem(todo *models.Todo, x, y, width, height int) *TodoItem {
	item := &TodoItem{
		Todo:   todo,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}

	// Create checkbox
	checkboxSize := 20
	item.Checkbox = NewButton(
		x+8, y+(height-checkboxSize)/2, 
		checkboxSize, checkboxSize,
		"", // No text, will draw custom checkbox
		func() {
			todo.Toggle()
		},
	)
	item.Checkbox.SetColors(
		color.RGBA{255, 255, 255, 255}, // White background
		color.RGBA{240, 240, 240, 255}, // Light gray hover
		color.RGBA{0, 0, 0, 255},       // Black text
	)

	// Create delete button
	deleteSize := 24
	item.DeleteBtn = NewButton(
		x+width-deleteSize-8, y+(height-deleteSize)/2,
		deleteSize, deleteSize,
		"×",
		func() {
			// Delete functionality will be handled by parent
		},
	)
	item.DeleteBtn.SetColors(
		color.RGBA{220, 53, 69, 255},   // Red background
		color.RGBA{200, 35, 51, 255},   // Darker red hover
		color.RGBA{255, 255, 255, 255}, // White text
	)

	// Create edit textbox (initially hidden)
	textboxX := x + 40
	textboxWidth := width - 80
	item.EditTextBox = NewTextBox(textboxX, y+4, textboxWidth, height-8, "Enter todo text")

	return item
}

func (ti *TodoItem) Update() {
	// Update position-dependent components if position changed
	ti.updateComponentPositions()

	// Handle mouse hover
	mouseX, mouseY := ebiten.CursorPosition()
	ti.Hovered = mouseX >= ti.X && mouseX <= ti.X+ti.Width && 
		       mouseY >= ti.Y && mouseY <= ti.Y+ti.Height

	if ti.Editing {
		ti.EditTextBox.Update()
		
		// Handle Enter key to save
		if ti.EditTextBox.IsEnterPressed() {
			newText := ti.EditTextBox.GetText()
			if len(newText) > 0 {
				ti.Todo.SetText(newText)
				ti.Editing = false
			}
		}
		
		// Handle Escape key to cancel
		if ti.EditTextBox.IsEscapePressed() {
			ti.Editing = false
		}
	} else {
		// Update checkbox and delete button only when not editing
		ti.Checkbox.Update()
		ti.DeleteBtn.Update()

		// Handle double-click to edit
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && ti.Hovered {
			currentTime := time.Now()
			if currentTime.Sub(ti.lastClickTime) < 300*time.Millisecond {
				// Double-click detected
				ti.startEditing()
			}
			ti.lastClickTime = currentTime
		}
	}
}

func (ti *TodoItem) startEditing() {
	ti.Editing = true
	ti.EditText = ti.Todo.Text
	ti.EditTextBox.SetText(ti.Todo.Text)
	ti.EditTextBox.SetFocus(true)
}

func (ti *TodoItem) Draw(screen *ebiten.Image) {
	// Draw background
	bgColor := color.RGBA{255, 255, 255, 255}
	if ti.Hovered {
		bgColor = color.RGBA{248, 249, 250, 255}
	}
	ebitenutil.DrawRect(screen, float64(ti.X), float64(ti.Y), float64(ti.Width), float64(ti.Height), bgColor)

	if ti.Editing {
		// Draw edit mode
		ti.EditTextBox.Draw(screen)
	} else {
		// Draw normal mode
		ti.drawCheckbox(screen)
		ti.drawTodoText(screen)
		ti.DeleteBtn.Draw(screen)
	}

	// Draw separator line
	separatorColor := color.RGBA{200, 200, 200, 255}
	ebitenutil.DrawRect(screen, float64(ti.X), float64(ti.Y+ti.Height-1), float64(ti.Width), 1, separatorColor)
}

func (ti *TodoItem) drawCheckbox(screen *ebiten.Image) {
	checkboxSize := 20
	x := ti.X + 8
	y := ti.Y + (ti.Height-checkboxSize)/2

	// Draw checkbox background
	bgColor := color.RGBA{255, 255, 255, 255}
	if ti.Checkbox.Hovered {
		bgColor = color.RGBA{240, 240, 240, 255}
	}
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(checkboxSize), float64(checkboxSize), bgColor)

	// Draw checkbox border
	borderColor := color.RGBA{108, 117, 125, 255}
	ebitenutil.DrawRect(screen, float64(x), float64(y), float64(checkboxSize), 1, borderColor)
	ebitenutil.DrawRect(screen, float64(x), float64(y), 1, float64(checkboxSize), borderColor)
	ebitenutil.DrawRect(screen, float64(x+checkboxSize-1), float64(y), 1, float64(checkboxSize), borderColor)
	ebitenutil.DrawRect(screen, float64(x), float64(y+checkboxSize-1), float64(checkboxSize), 1, borderColor)

	// Draw checkmark if completed
	if ti.Todo.Completed {
		checkColor := color.RGBA{40, 167, 69, 255} // Green
		// Simple checkmark using rectangles
		ebitenutil.DrawRect(screen, float64(x+4), float64(y+10), 3, 2, checkColor)
		ebitenutil.DrawRect(screen, float64(x+6), float64(y+12), 2, 2, checkColor)
		ebitenutil.DrawRect(screen, float64(x+8), float64(y+8), 2, 6, checkColor)
		ebitenutil.DrawRect(screen, float64(x+10), float64(y+6), 2, 4, checkColor)
		ebitenutil.DrawRect(screen, float64(x+12), float64(y+4), 2, 4, checkColor)
		ebitenutil.DrawRect(screen, float64(x+14), float64(y+6), 2, 2, checkColor)
	}
}

func (ti *TodoItem) drawTodoText(screen *ebiten.Image) {
	textX := ti.X + 40
	textY := ti.Y + (ti.Height+text.BoundString(basicfont.Face7x13, "A").Max.Y)/2

	textColor := color.RGBA{33, 37, 41, 255}
	displayText := ti.Todo.Text

	// Truncate text if too long
	maxWidth := ti.Width - 80 // Account for checkbox and delete button
	if text.BoundString(basicfont.Face7x13, displayText).Max.X > maxWidth {
		for len(displayText) > 0 {
			if text.BoundString(basicfont.Face7x13, displayText+"...").Max.X <= maxWidth {
				displayText += "..."
				break
			}
			displayText = displayText[:len(displayText)-1]
		}
	}

	// Apply strikethrough for completed tasks
	if ti.Todo.Completed {
		textColor = color.RGBA{108, 117, 125, 255} // Gray color for completed
	}

	text.Draw(screen, displayText, basicfont.Face7x13, textX, textY, textColor)

	// Draw strikethrough line for completed tasks
	if ti.Todo.Completed {
		textBounds := text.BoundString(basicfont.Face7x13, displayText)
		lineY := textY - textBounds.Max.Y/2
		lineWidth := textBounds.Max.X
		ebitenutil.DrawRect(screen, float64(textX), float64(lineY), float64(lineWidth), 1, textColor)
	}
}

func (ti *TodoItem) updateComponentPositions() {
	checkboxSize := 20
	ti.Checkbox.SetPosition(ti.X+8, ti.Y+(ti.Height-checkboxSize)/2)

	deleteSize := 24
	ti.DeleteBtn.SetPosition(ti.X+ti.Width-deleteSize-8, ti.Y+(ti.Height-deleteSize)/2)

	textboxX := ti.X + 40
	textboxWidth := ti.Width - 80
	ti.EditTextBox.X = textboxX
	ti.EditTextBox.Y = ti.Y + 4
	ti.EditTextBox.Width = textboxWidth
	ti.EditTextBox.Height = ti.Height - 8
}

func (ti *TodoItem) SetPosition(x, y int) {
	ti.X = x
	ti.Y = y
	ti.updateComponentPositions()
}

func (ti *TodoItem) SetWidth(width int) {
	ti.Width = width
	ti.updateComponentPositions()
}

func (ti *TodoItem) IsDeleteRequested() bool {
	// This should be called after Update() to check if delete was clicked
	return false // Will be handled by checking DeleteBtn state in parent
}

func (ti *TodoItem) GetDeleteButton() *Button {
	return ti.DeleteBtn
}

func (ti *TodoItem) Contains(x, y int) bool {
	return x >= ti.X && x <= ti.X+ti.Width && y >= ti.Y && y <= ti.Y+ti.Height
}