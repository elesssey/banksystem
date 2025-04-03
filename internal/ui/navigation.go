package ui

import (
	"fyne.io/fyne/v2"

	"banksystem/internal/service"
	"banksystem/internal/ui/screens"
	"banksystem/internal/ui/state"
)

type ScreenID int

const (
	ScreenNone ScreenID = iota
	ScreenError
	ScreenLogin
	ScreenBankSelector
	ScreenBank
	ScreenTransaction
	ScreenAdminMain
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
		n.window.SetContent(screens.MakeLoginScreen(n.onLoginClick, n.handleSuccessfulLogin))
	case ScreenBankSelector:
		if err := n.initializeBankPageData(); err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenLogin) })
			return
		}
		if n.state.User.GetCurrentUser().Role == "admin" {
			n.window.SetContent(screens.MakeBankSelectorScreen(n.openAdminMain, n.state.Banks))
		} else {
			n.window.SetContent(screens.MakeBankSelectorScreen(n.openBankPage, n.state.Banks))
		}
	case ScreenBank:
		user := n.state.User.GetCurrentUser()
		userAccount, err := n.bankingService.GetUserAccount(user.ID, n.state.Banks.GetCurrentBank().ID)
		if err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenBankSelector) })
			return
		}
		n.state.Banks.WorkingAccount = userAccount
		n.window.SetContent(screens.MakeBankScreen(n.openTransactionPage, n.state.Banks, user))

	case ScreenTransaction:
		user := n.state.User.GetCurrentUser()
		n.state.Transaction = state.NewTransactionState(
			n.state.Banks.BanksList, user,
			n.state.Banks.WorkingAccount,
			n.state.Banks.GetCurrentBank(),
		)
		n.window.SetContent(screens.MakeTransactionScreen(n.createTransaction, n.onCreateTransactionError, n.state.Transaction))
	case ScreenAdminMain:
		if err := n.initializeAdminPageData(n.state.Banks.GetCurrentBank().ID); err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenBankSelector) })
			return
		}
		n.window.SetContent(screens.MakeAdminMain(n.state.Banks.AdminTransactionsList, n.state.Banks.FindBankNameById, n.adminConfirmationTransaction))
	}
}
