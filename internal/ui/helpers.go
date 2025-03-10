package ui

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (n *NavigationManager) showError(message string, onOk func()) {
	n.currentScreen = ScreenError
	n.window.SetContent(container.NewVBox(
		widget.NewLabel(message),
		widget.NewButton("OK", onOk),
	))
}
