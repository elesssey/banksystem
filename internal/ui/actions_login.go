package ui

import (
	"banksystem/internal/model"

	"fyne.io/fyne/v2/dialog"
)

func (n *NavigationManager) onLoginClick(email string, password string) {
	if email == "" {
		n.showError("Введите Email", func() { n.navigateTo(ScreenLogin) })
		return
	}

	user, err := n.authService.Login(email, password)
	if err != nil {
		n.showError("Ошибка: "+err.Error(), func() { n.navigateTo(ScreenLogin) })
		return
	}

	n.handleSuccessfulLogin(user)
}

func (n *NavigationManager) handleSuccessfulLogin(user *model.User) {
	n.state.User.SetCurrentUser(user)
	n.navigateTo(ScreenBankSelector)
	dialog.ShowInformation("Информация!", "Вы вошли в аккаунт!", n.window)
}
