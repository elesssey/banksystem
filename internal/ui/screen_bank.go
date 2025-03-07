package ui

import (
	"banksystem/internal/model"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeBankPage(bank *model.Bank, user *model.User, account *model.User_Account) fyne.CanvasObject {
	heading := widget.NewLabelWithStyle(bank.Name, fyne.TextAlignLeading, fyne.TextStyle{})

	Label1 := canvas.NewText("Существующие счета:", color.Black)
	Label2 := canvas.NewText("Вы можете перести деньги другому пользователю используя дебетовый счет", color.RGBA{255, 0, 0, 1})
	debitLabel1 := widget.NewLabelWithStyle("Дебетовый счет", fyne.TextAlignCenter, fyne.TextStyle{})
	debitLabel2 := canvas.NewText("Номер счета:", color.Black)
	debitLabel3 := canvas.NewText("Баланс:", color.Black)
	debitNumber := widget.NewLabelWithStyle(account.Number, fyne.TextAlignCenter, fyne.TextStyle{})
	debitBalance := widget.NewLabelWithStyle(fmt.Sprintf("%f", account.Balance), fyne.TextAlignCenter, fyne.TextStyle{})
	savingLabel := widget.NewLabelWithStyle("Накопительный счет", fyne.TextAlignCenter, fyne.TextStyle{})
	creditLabel := widget.NewLabelWithStyle("Кредиторный счет", fyne.TextAlignCenter, fyne.TextStyle{})

	infButton := widget.NewButton("Посмотреть информацию", func() {})
	infButton.Resize(fyne.NewSize(200, 50))
	transferButton := widget.NewButton("СДЕЛАТЬ ПЕРЕВОД", func() {})

	imgUser := canvas.NewImageFromFile("D:/secondcurse/4sem/oop/banksystem/banksystem/internal/images/userLogo.png")
	imgUser.SetMinSize(fyne.NewSize(50, 50))
	imgEmail := canvas.NewImageFromFile("D:/secondcurse/4sem/oop/banksystem/banksystem/internal/images/email.png")
	imgEmail.SetMinSize(fyne.NewSize(50, 50))
	imgPlus := canvas.NewImageFromFile("D:/secondcurse/4sem/oop/banksystem/banksystem/internal/images/plus.png")
	imgPlus.SetMinSize(fyne.NewSize(150, 150))

	borderHeading := canvas.NewRectangle(color.Black)
	borderHeading.StrokeWidth = 2
	borderHeading.StrokeColor = color.Black
	borderHeading.FillColor = color.Transparent

	borderNameInfo := widget.NewLabelWithStyle(fmt.Sprintf("%s %s %s", user.Surname, user.Name, user.MiddleName), fyne.TextAlignCenter, fyne.TextStyle{})
	borderEmailInfo := widget.NewLabelWithStyle(fmt.Sprintf(user.Email), fyne.TextAlignCenter, fyne.TextStyle{})

	border1 := container.NewStack(borderHeading, borderNameInfo)
	border2 := container.NewStack(borderHeading, borderEmailInfo)

	form := container.NewVBox(
		container.NewHBox(imgUser, border1),
		container.NewHBox(imgEmail, border2),
	)

	debitBody := container.NewGridWithRows(5,
		debitLabel1,
		container.NewHBox(debitLabel2, debitNumber),
		container.NewHBox(debitLabel3, debitBalance),
		container.NewHBox(layout.NewSpacer(), infButton, layout.NewSpacer()),
		layout.NewSpacer(),
	)

	savingBody := container.NewVBox(
		savingLabel,
		container.NewHBox(layout.NewSpacer(), imgPlus, layout.NewSpacer()),
	)
	creditBody := container.NewVBox(
		creditLabel,
		container.NewHBox(layout.NewSpacer(), imgPlus, layout.NewSpacer()),
	)

	debit_account := container.NewStack(borderHeading, debitBody)
	saving_account := container.NewStack(borderHeading, savingBody)
	credit_account := container.NewStack(borderHeading, creditBody)
	accounts := container.NewGridWithColumns(3, debit_account, saving_account, credit_account)

	mainContainer := container.NewVBox(
		heading,
		layout.NewSpacer(),
		container.NewHBox(layout.NewSpacer(), form),
		layout.NewSpacer(),
		Label1,
		layout.NewSpacer(),
		accounts,
		layout.NewSpacer(),
		Label2,
		layout.NewSpacer(),
		transferButton,
		layout.NewSpacer(),
	)
	return mainContainer
}
