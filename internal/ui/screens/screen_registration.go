package screens

import (
	"banksystem/internal/model"
	"banksystem/internal/ui/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeRegistratePage(onRegistrateTry func(*model.User, int), banksState *state.BanksState) fyne.CanvasObject {
	heading := widget.NewLabelWithStyle("РЕГИСТРАЦИЯ", fyne.TextAlignCenter, fyne.TextStyle{})
	label := widget.NewLabelWithStyle("Заполните все поля для создания нового аккаунта", fyne.TextAlignLeading, fyne.TextStyle{})
	label2 := widget.NewLabelWithStyle("Выберите банк для создания начального счета:", fyne.TextAlignLeading, fyne.TextStyle{})

	nameLabel := widget.NewLabelWithStyle("Имя:                        ", fyne.TextAlignCenter, fyne.TextStyle{})
	nameEntry := widget.NewEntry()
	surnameLabel := widget.NewLabelWithStyle("Фамилия:              ", fyne.TextAlignCenter, fyne.TextStyle{})
	surnameEntry := widget.NewEntry()
	thirdNameLabel := widget.NewLabelWithStyle("Отчество:              ", fyne.TextAlignCenter, fyne.TextStyle{})
	thirdNameEntry := widget.NewEntry()
	passwordLabel := widget.NewLabelWithStyle("Пароль:                 ", fyne.TextAlignCenter, fyne.TextStyle{})
	passwordEntry := widget.NewEntry()
	passportSeriesLabel := widget.NewLabelWithStyle("Серия паспорта: ", fyne.TextAlignCenter, fyne.TextStyle{})
	seriesEntry := widget.NewEntry()
	passportNumberLabel := widget.NewLabelWithStyle("Номер паспорта:", fyne.TextAlignCenter, fyne.TextStyle{})
	numberEntry := widget.NewEntry()
	phoneLabel := widget.NewLabelWithStyle("Телефон:               ", fyne.TextAlignCenter, fyne.TextStyle{})
	phoneEntry := widget.NewEntry()
	emailLabel := widget.NewLabelWithStyle("Почта:                    ", fyne.TextAlignCenter, fyne.TextStyle{})
	emailEntry := widget.NewEntry()

	var bankId int

	registrateButton := widget.NewButton("ЗАРЕГИСТРИРОВАТЬСЯ", func() {
		newUser := &model.User{
			Name:           nameEntry.Text,
			Surname:        surnameEntry.Text,
			MiddleName:     thirdNameEntry.Text,
			Password:       passwordEntry.Text,
			PassportSeries: seriesEntry.Text,
			PassportNumber: numberEntry.Text,
			Phone:          phoneEntry.Text,
			Email:          emailEntry.Text,
			Role:           "client",
		}
		onRegistrateTry(newUser, bankId)
	})

	radio := widget.NewRadioGroup(banksState.GetBankStateNames(), func(value string) {
		for _, bank := range banksState.BanksList {
			if bank.Name == value {
				bankId = bank.ID
			}
		}
	})

	form := container.New(layout.NewCustomPaddedLayout(20, 20, 0, 1000), container.NewVBox(
		container.New(layout.NewFormLayout(), nameLabel, nameEntry),
		container.New(layout.NewFormLayout(), surnameLabel, surnameEntry),
		container.New(layout.NewFormLayout(), thirdNameLabel, thirdNameEntry),
		container.New(layout.NewFormLayout(), passwordLabel, passwordEntry),
		container.New(layout.NewFormLayout(), passportSeriesLabel, seriesEntry),
		container.New(layout.NewFormLayout(), passportNumberLabel, numberEntry),
		container.New(layout.NewFormLayout(), phoneLabel, phoneEntry),
		container.New(layout.NewFormLayout(), emailLabel, emailEntry),
	))
	mainContainer := container.NewVBox(
		heading,
		label,
		form,
		label2,
		radio,
		registrateButton,
	)
	return mainContainer
}
