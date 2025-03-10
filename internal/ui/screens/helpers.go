package screens

import (
	"image/color"

	"fyne.io/fyne/v2/canvas"
)

func newRectangle() *canvas.Rectangle {
	rect := &canvas.Rectangle{}
	rect.StrokeWidth = 2
	rect.StrokeColor = color.Black
	rect.FillColor = color.Transparent
	return rect
}
