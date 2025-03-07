package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"banksystem/internal/service"
	"banksystem/internal/ui/state"
)

type ScreenID int

const (
	ScreenNone ScreenID = iota
	ScreenError
	ScreenLogin
	ScreenBankSelector
	ScreenBank
)

type NavigationManager struct {
	app            fyne.App
	window         fyne.Window
	state          *state.AppState
	authService    service.AuthService
	bankingService service.BankingService

	currentScreen ScreenID
}

func NewNavigationManager(
	app fyne.App,
	window fyne.Window,
	state *state.AppState,
	authService service.AuthService,
	bankingServicce service.BankingService,
) *NavigationManager {
	return &NavigationManager{
		app:            app,
		window:         window,
		state:          state,
		authService:    authService,
		bankingService: bankingServicce,
		currentScreen:  ScreenNone,
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
		if err := n.initializeBankPageData(); err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenLogin) })
			return
		}
		n.window.SetContent(MakeBankSelectorScreen(n.openBankPage, n.state.Banks.Banks[0], n.state.Banks.Banks[1], n.state.Banks.Banks[2]))
	case ScreenBank:
		user := n.state.User.GetCurrentUser()
		user_account, err := n.bankingService.GetUserAccount(user.ID)
		if err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenBankSelector) })
			return
		}
		n.window.SetContent(MakeBankPage(n.state.Banks.Banks[n.state.Banks.SelectedBankIndex], user, user_account))
	}
}

func (n *NavigationManager) showError(message string, onOk func()) {
	n.currentScreen = ScreenError
	n.window.SetContent(container.NewVBox(
		widget.NewLabel(message),
		widget.NewButton("OK", onOk),
	))
}
