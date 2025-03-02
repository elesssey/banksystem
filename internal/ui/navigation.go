package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"banksystem/internal/service"
	"banksystem/internal/state"
)

type ScreenID int

const (
	ScreenNone ScreenID = iota
	ScreenError
	ScreenLogin
	ScreenBankSelector
)

type NavigationManager struct {
	app         fyne.App
	window      fyne.Window
	state       *state.AppState
	authService service.AuthService

	currentScreen ScreenID
}

func NewNavigationManager(
	app fyne.App,
	window fyne.Window,
	state *state.AppState,
	authService service.AuthService,
) *NavigationManager {
	return &NavigationManager{
		app:           app,
		window:        window,
		state:         state,
		authService:   authService,
		currentScreen: ScreenNone,
	}
}

func (n *NavigationManager) Start() {
	n.navigateTo(ScreenLogin)
}

func (n *NavigationManager) navigateTo(screenID ScreenID) {
	if n.currentScreen == screenID {
		return
	}

	n.currentScreen = screenID

	switch screenID {
	case ScreenLogin:
		n.window.SetContent(MakeLoginScreen(n.onLoginClick, n.handleSuccessfulLogin))
	case ScreenBankSelector:
		n.window.SetContent(MakeBankSelectorScreen())
	}
}

func (n *NavigationManager) showError(message string, onOk func()) {
	n.currentScreen = ScreenError
	n.window.SetContent(container.NewVBox(
		widget.NewLabel(message),
		widget.NewButton("OK", onOk),
	))
}
