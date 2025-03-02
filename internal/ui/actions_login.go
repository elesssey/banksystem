package ui

import (
	"banksystem/internal/model"

	"fyne.io/fyne/v2/dialog"
)

func (n *NavigationManager) onLoginClick(email string) {
	if email == "" {
		n.showError("Введите Email", func() { n.navigateTo(ScreenLogin) })
		return
	}

	user, err := n.authService.Login(email)
	if err != nil {
		n.showError("Ошибка: "+err.Error(), func() { n.navigateTo(ScreenLogin) })
		return
	}

	n.handleSuccessfulLogin(user)
}

func (n *NavigationManager) handleSuccessfulLogin(user *model.User) {
	n.state.User.SetCurrentUser(user)

	dialog.ShowInformation("Info", "You are logged in!", n.window)
}
