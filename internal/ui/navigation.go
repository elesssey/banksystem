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
	ScreenRegistrate
	ScreenCredit
	ScreenWatchLogs
)

type NavigationManager struct {
	app            fyne.App
	window         fyne.Window
	state          *state.AppState
	authService    service.AuthService
	bankingService service.BankingService

	currentScreen  ScreenID
	previousScreen ScreenID
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
		previousScreen: ScreenNone,
	}
}

func (n *NavigationManager) Start() {
	n.navigateTo(ScreenLogin)
}

func (n *NavigationManager) navigateTo(screenID ScreenID) {
	if n.currentScreen == ScreenNone {
		n.previousScreen = ScreenNone
	} else {
		n.previousScreen = n.currentScreen
	}
	if n.currentScreen == screenID {
		return
	}

	n.currentScreen = screenID

	switch screenID {
	case ScreenLogin:
		n.window.SetContent(screens.MakeLoginScreen(n.onLoginClick, n.RegistrationStart))
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
		n.window.SetContent(screens.MakeBankScreen(n.openLogsPage, n.openCreditPage, n.openTransactionPage, n.state.Banks, user))

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
		n.window.SetContent(screens.MakeAdminMain(n.state.Banks.AdminCreditList,
			n.state.Banks.AdminTransactionsList,
			n.state.Banks.FindBankNameById,
			n.adminConfirmationTransaction,
			n.adminDeclineTransaction,
			n.adminConfirmationCredit,
			n.adminDeclineCredit,
			n.backToPreviousPage))
	case ScreenRegistrate:
		if err := n.initializeBankPageData(); err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenLogin) })
			return
		}
		n.window.SetContent(screens.MakeRegistratePage(n.onRegistrateClick, n.state.Banks, n.backToPreviousPage))
	case ScreenCredit:
		user := n.state.User.GetCurrentUser()
		n.state.Credit = state.NewCreditState(
			n.state.Banks.BanksList, user,
			n.state.Banks.WorkingAccount,
			n.state.Banks.GetCurrentBank(),
		)
		n.window.SetContent(screens.MakeCreditScreen(n.createCredit, n.onCreateCreditError, n.state.Credit))
	case ScreenWatchLogs:
		n.state.User.UserTransactionList = nil
		n.state.User.UserCreditList = nil
		user := n.state.User.GetCurrentUser()
		if err := n.initializeUserPageData(n.state.Banks.GetCurrentBank().ID, user.ID); err != nil {
			n.showError(err.Error(), func() { n.navigateTo(ScreenLogin) })
			return
		}
		n.window.SetContent(screens.MakeUserPageData(n.state.User.UserTransactionList, n.state.User.UserCreditList, n.state.Banks.FindBankNameById, n.backToPreviousPage))
	}
}
