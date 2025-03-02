package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func MakeBankSelectorScreen() fyne.CanvasObject {
	banksContainer := container.NewGridWithColumns(3)
	mainContainer := container.NewGridWithRows(4, banksContainer)
	return mainContainer
}
