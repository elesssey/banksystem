package ui

import (
	"banksystem/internal/model"

	"fyne.io/fyne/v2/dialog"
)

func (n *NavigationManager) onRegistrateClick(user *model.User, bankId int) {
	if len(user.Password) < 6 {
		n.showError("Введите пароль хотя бы из шести символов", func() { n.navigateTo(ScreenRegistrate) })
		return
	}
	if user.Name == "" {
		n.showError("Введите свое имя", func() { n.navigateTo(ScreenRegistrate) })
		return
	}
	if user.Surname == "" {
		n.showError("Введите свою фамилию", func() { n.navigateTo(ScreenRegistrate) })
		return
	}
	if user.MiddleName == "" {
		n.showError("Введите свое отчество", func() { n.navigateTo(ScreenRegistrate) })
		return
	}
	if len(user.PassportSeries) != 2 {
		n.showError("Введите правильно серию своего паспорта (ДВЕ ЗАГЛАВНЫЕ БУКВЫ)", func() { n.navigateTo(ScreenRegistrate) })
		return
	}
	if len(user.PassportNumber) != 7 {
		n.showError("Введите правильно номер своего паспорта (СЕМЬ ЦИФР)", func() { n.navigateTo(ScreenRegistrate) })
		return
	}
	if user.Phone == "" || len(user.Phone) != 13 {
		n.showError("Введите свой телефон", func() { n.navigateTo(ScreenRegistrate) })
		return
	}

	if user.Email == "" {
		n.showError("Введите свою почту", func() { n.navigateTo(ScreenRegistrate) })
		return
	}

	err := n.authService.Registrate(user, bankId)
	if err != nil {
		return
	}
	n.handleSuccessfulRegistrate(user)
}

func (n *NavigationManager) handleSuccessfulRegistrate(user *model.User) {
	n.state.User.SetCurrentUser(user)
	n.navigateTo(ScreenBankSelector)
	dialog.ShowInformation("Информация:", "Вы зарегистрировались!", n.window)
}
