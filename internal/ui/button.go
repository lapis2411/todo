package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	X, Y, Width, Height int
	Text               string
	OnClick            func()
	Hovered            bool
	Pressed            bool
	Enabled            bool
	BackgroundColor    color.RGBA
	HoverColor         color.RGBA
	TextColor          color.RGBA
}

func NewButton(x, y, width, height int, text string, onClick func()) *Button {
	return &Button{
		X:               x,
		Y:               y,
		Width:           width,
		Height:          height,
		Text:            text,
		OnClick:         onClick,
		Enabled:         true,
		BackgroundColor: color.RGBA{0, 123, 255, 255},   // Primary blue
		HoverColor:      color.RGBA{0, 86, 179, 255},    // Darker blue
		TextColor:       color.RGBA{255, 255, 255, 255}, // White
	}
}

func (b *Button) Update() {
	if !b.Enabled {
		return
	}

	x, y := ebiten.CursorPosition()
	
	// Check if cursor is over button
	b.Hovered = x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height

	// Check for click
	if b.Hovered && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !b.Pressed {
			b.Pressed = true
		}
	} else if b.Pressed {
		// Button was pressed and now released
		if b.Hovered && b.OnClick != nil {
			b.OnClick()
		}
		b.Pressed = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	bgColor := b.BackgroundColor
	if !b.Enabled {
		bgColor = color.RGBA{108, 117, 125, 255} // Disabled gray
	} else if b.Hovered || b.Pressed {
		bgColor = b.HoverColor
	}

	// Draw background
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), float64(b.Width), float64(b.Height), bgColor)

	// Draw border
	borderColor := color.RGBA{0, 0, 0, 100}
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), float64(b.Width), 1, borderColor)
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y), 1, float64(b.Height), borderColor)
	ebitenutil.DrawRect(screen, float64(b.X+b.Width-1), float64(b.Y), 1, float64(b.Height), borderColor)
	ebitenutil.DrawRect(screen, float64(b.X), float64(b.Y+b.Height-1), float64(b.Width), 1, borderColor)

	// Draw text
	textColor := b.TextColor
	if !b.Enabled {
		textColor = color.RGBA{200, 200, 200, 255}
	}

	// Center text in button
	bounds := text.BoundString(basicfont.Face7x13, b.Text)
	textWidth := bounds.Max.X - bounds.Min.X
	textHeight := bounds.Max.Y - bounds.Min.Y
	
	textX := b.X + (b.Width-textWidth)/2
	textY := b.Y + (b.Height+textHeight)/2

	text.Draw(screen, b.Text, basicfont.Face7x13, textX, textY, textColor)
}

func (b *Button) SetEnabled(enabled bool) {
	b.Enabled = enabled
}

func (b *Button) Contains(x, y int) bool {
	return x >= b.X && x <= b.X+b.Width && y >= b.Y && y <= b.Y+b.Height
}

func (b *Button) SetPosition(x, y int) {
	b.X = x
	b.Y = y
}

func (b *Button) SetSize(width, height int) {
	b.Width = width
	b.Height = height
}

func (b *Button) SetText(text string) {
	b.Text = text
}

func (b *Button) SetColors(bg, hover, text color.RGBA) {
	b.BackgroundColor = bg
	b.HoverColor = hover
	b.TextColor = text
}