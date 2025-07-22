package ui

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type TextBox struct {
	X, Y, Width, Height int
	Text               string
	Focused            bool
	CursorPos          int
	ShowCursor         bool
	lastCursorToggle   time.Time
	BackgroundColor    color.RGBA
	BorderColor        color.RGBA
	FocusedBorderColor color.RGBA
	TextColor          color.RGBA
	PlaceholderText    string
	PlaceholderColor   color.RGBA
	MaxLength          int
}

func NewTextBox(x, y, width, height int, placeholder string) *TextBox {
	return &TextBox{
		X:                  x,
		Y:                  y,
		Width:              width,
		Height:             height,
		Text:               "",
		PlaceholderText:    placeholder,
		CursorPos:          0,
		ShowCursor:         true,
		lastCursorToggle:   time.Now(),
		BackgroundColor:    color.RGBA{255, 255, 255, 255},
		BorderColor:        color.RGBA{108, 117, 125, 255},
		FocusedBorderColor: color.RGBA{0, 123, 255, 255},
		TextColor:          color.RGBA{33, 37, 41, 255},
		PlaceholderColor:   color.RGBA{108, 117, 125, 255},
		MaxLength:          100,
	}
}

func (tb *TextBox) Update() {
	// Handle mouse click for focus
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tb.Focused = tb.Contains(x, y)
	}

	if !tb.Focused {
		return
	}

	// Handle cursor blinking
	if time.Since(tb.lastCursorToggle) > 500*time.Millisecond {
		tb.ShowCursor = !tb.ShowCursor
		tb.lastCursorToggle = time.Now()
	}

	// Handle text input
	runes := ebiten.AppendInputChars(nil)
	for _, r := range runes {
		if r == '\n' || r == '\r' {
			// Enter key handling is done by caller
			continue
		}
		if r == '\b' {
			// Backspace
			if tb.CursorPos > 0 {
				tb.Text = tb.Text[:tb.CursorPos-1] + tb.Text[tb.CursorPos:]
				tb.CursorPos--
			}
		} else if len(tb.Text) < tb.MaxLength && r >= 32 {
			// Regular character input
			tb.Text = tb.Text[:tb.CursorPos] + string(r) + tb.Text[tb.CursorPos:]
			tb.CursorPos++
		}
	}

	// Handle arrow keys
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) && tb.CursorPos > 0 {
		tb.CursorPos--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) && tb.CursorPos < len(tb.Text) {
		tb.CursorPos++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyHome) {
		tb.CursorPos = 0
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnd) {
		tb.CursorPos = len(tb.Text)
	}

	// Handle Delete key
	if inpututil.IsKeyJustPressed(ebiten.KeyDelete) && tb.CursorPos < len(tb.Text) {
		tb.Text = tb.Text[:tb.CursorPos] + tb.Text[tb.CursorPos+1:]
	}

	// Handle Backspace key (additional handling for better responsiveness)
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && tb.CursorPos > 0 {
		tb.Text = tb.Text[:tb.CursorPos-1] + tb.Text[tb.CursorPos:]
		tb.CursorPos--
	}
}

func (tb *TextBox) Draw(screen *ebiten.Image) {
	// Draw background
	ebitenutil.DrawRect(screen, float64(tb.X), float64(tb.Y), float64(tb.Width), float64(tb.Height), tb.BackgroundColor)

	// Draw border
	borderColor := tb.BorderColor
	if tb.Focused {
		borderColor = tb.FocusedBorderColor
	}
	
	// Draw border with 2px width for focused state
	borderWidth := 1
	if tb.Focused {
		borderWidth = 2
	}
	
	for i := 0; i < borderWidth; i++ {
		ebitenutil.DrawRect(screen, float64(tb.X-i), float64(tb.Y-i), float64(tb.Width+2*i), 1, borderColor)
		ebitenutil.DrawRect(screen, float64(tb.X-i), float64(tb.Y-i), 1, float64(tb.Height+2*i), borderColor)
		ebitenutil.DrawRect(screen, float64(tb.X+tb.Width-1+i), float64(tb.Y-i), 1, float64(tb.Height+2*i), borderColor)
		ebitenutil.DrawRect(screen, float64(tb.X-i), float64(tb.Y+tb.Height-1+i), float64(tb.Width+2*i), 1, borderColor)
	}

	// Draw text or placeholder
	displayText := tb.Text
	textColor := tb.TextColor
	
	if displayText == "" && !tb.Focused {
		displayText = tb.PlaceholderText
		textColor = tb.PlaceholderColor
	}

	// Calculate text position (with padding)
	padding := 8
	textX := tb.X + padding
	textY := tb.Y + (tb.Height+text.BoundString(basicfont.Face7x13, "A").Max.Y)/2

	// Draw visible portion of text (simple horizontal scrolling)
	visibleText := displayText
	maxVisibleWidth := tb.Width - padding*2 - 10 // Reserve space for cursor
	
	if text.BoundString(basicfont.Face7x13, visibleText).Max.X > maxVisibleWidth {
		// Adjust visible text to keep cursor visible
		startPos := 0
		for startPos < len(displayText) {
			testText := displayText[startPos:]
			testWidth := text.BoundString(basicfont.Face7x13, testText).Max.X
			testCursorOffset := 0
			
			if tb.CursorPos >= startPos {
				relativeCursor := tb.CursorPos - startPos
				if relativeCursor <= len(testText) {
					testCursorOffset = text.BoundString(basicfont.Face7x13, testText[:relativeCursor]).Max.X
				}
			}
			
			if testWidth <= maxVisibleWidth && testCursorOffset <= maxVisibleWidth {
				break
			}
			startPos++
		}
		
		visibleText = displayText[startPos:]
		if text.BoundString(basicfont.Face7x13, visibleText).Max.X > maxVisibleWidth {
			// Truncate from the end
			for len(visibleText) > 0 {
				if text.BoundString(basicfont.Face7x13, visibleText).Max.X <= maxVisibleWidth {
					break
				}
				visibleText = visibleText[:len(visibleText)-1]
			}
		}
	}

	text.Draw(screen, visibleText, basicfont.Face7x13, textX, textY, textColor)

	// Draw cursor
	if tb.Focused && tb.ShowCursor && tb.CursorPos >= 0 {
		cursorX := textX
		if tb.CursorPos > 0 && tb.CursorPos <= len(tb.Text) {
			textBeforeCursor := tb.Text[:tb.CursorPos]
			// Adjust for visible portion
			if len(visibleText) < len(tb.Text) {
				startPos := len(tb.Text) - len(visibleText)
				if tb.CursorPos >= startPos {
					relativeCursor := tb.CursorPos - startPos
					if relativeCursor <= len(visibleText) {
						textBeforeCursor = visibleText[:relativeCursor]
					} else {
						textBeforeCursor = visibleText
					}
				} else {
					textBeforeCursor = ""
				}
			}
			cursorX += text.BoundString(basicfont.Face7x13, textBeforeCursor).Max.X
		}
		
		cursorY := tb.Y + 4
		cursorHeight := tb.Height - 8
		ebitenutil.DrawRect(screen, float64(cursorX), float64(cursorY), 1, float64(cursorHeight), tb.TextColor)
	}
}

func (tb *TextBox) Contains(x, y int) bool {
	return x >= tb.X && x <= tb.X+tb.Width && y >= tb.Y && y <= tb.Y+tb.Height
}

func (tb *TextBox) SetFocus(focused bool) {
	tb.Focused = focused
	if focused {
		tb.ShowCursor = true
		tb.lastCursorToggle = time.Now()
	}
}

func (tb *TextBox) GetText() string {
	return tb.Text
}

func (tb *TextBox) SetText(text string) {
	tb.Text = text
	tb.CursorPos = len(text)
}

func (tb *TextBox) Clear() {
	tb.Text = ""
	tb.CursorPos = 0
}

func (tb *TextBox) IsEnterPressed() bool {
	return tb.Focused && inpututil.IsKeyJustPressed(ebiten.KeyEnter)
}

func (tb *TextBox) IsEscapePressed() bool {
	return tb.Focused && inpututil.IsKeyJustPressed(ebiten.KeyEscape)
}