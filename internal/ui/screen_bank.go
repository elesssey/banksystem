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

func MakeBankPage(onTransactionClick func(), bank *model.Bank, user *model.User, account *model.UserAccount) fyne.CanvasObject {
	heading := widget.NewLabelWithStyle(bank.Name, fyne.TextAlignLeading, fyne.TextStyle{})

	label1 := canvas.NewText("Существующие счета:", color.Black)
	label2 := canvas.NewText("Вы можете перести деньги другому пользователю используя дебетовый счет", color.RGBA{255, 0, 0, 255})
	debitLabel1 := widget.NewLabelWithStyle("Дебетовый счет", fyne.TextAlignCenter, fyne.TextStyle{})
	debitLabel2 := canvas.NewText("Номер счета:", color.Black)
	debitLabel3 := canvas.NewText("Баланс:", color.Black)
	debitNumber := widget.NewLabelWithStyle(account.Number, fyne.TextAlignCenter, fyne.TextStyle{})
	debitBalance := widget.NewLabelWithStyle(fmt.Sprintf("%.2f %s", account.Balance, account.Currency), fyne.TextAlignCenter, fyne.TextStyle{})
	savingLabel := widget.NewLabelWithStyle("Накопительный счет", fyne.TextAlignCenter, fyne.TextStyle{})
	creditLabel := widget.NewLabelWithStyle("Кредиторный счет", fyne.TextAlignCenter, fyne.TextStyle{})

	infButton := widget.NewButton("Посмотреть информацию", func() {})
	infButton.Resize(fyne.NewSize(200, 50))
	transferButton := widget.NewButton("СДЕЛАТЬ ПЕРЕВОД", func() { onTransactionClick() })

	imgUser := canvas.NewImageFromFile("./internal/images/userLogo.png")
	imgUser.SetMinSize(fyne.NewSize(50, 50))
	imgEmail := canvas.NewImageFromFile("./internal/images/email.png")
	imgEmail.SetMinSize(fyne.NewSize(50, 50))
	imgPlus := canvas.NewImageFromFile("./internal/images/plus.png")
	imgPlus.SetMinSize(fyne.NewSize(100, 100))

	borderNameInfo := widget.NewLabelWithStyle(fmt.Sprintf("%s %s %s", user.Surname, user.Name, user.MiddleName), fyne.TextAlignCenter, fyne.TextStyle{})
	borderEmailInfo := widget.NewLabelWithStyle(fmt.Sprintf(user.Email), fyne.TextAlignCenter, fyne.TextStyle{})

	border1 := container.NewStack(newRectangle(), borderNameInfo)
	border2 := container.NewStack(newRectangle(), borderEmailInfo)

	form := container.NewGridWithRows(2,
		container.New(layout.NewGridLayoutWithColumns(2), container.NewCenter(imgUser), border1),
		container.New(layout.NewGridLayoutWithColumns(2), container.NewCenter(imgEmail), border2),
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

	debitAccount := container.NewStack(newRectangle(), debitBody)
	savingAccount := container.NewStack(newRectangle(), savingBody)
	creditAccount := container.NewStack(newRectangle(), creditBody)
	accounts := container.NewGridWithColumns(3, debitAccount, savingAccount, creditAccount)

	mainContainer := container.NewVBox(
		heading,
		layout.NewSpacer(),
		container.NewGridWithColumns(3, layout.NewSpacer(), layout.NewSpacer(), form),
		layout.NewSpacer(),
		label1,
		layout.NewSpacer(),
		accounts,
		layout.NewSpacer(),
		label2,
		layout.NewSpacer(),
		transferButton,
		layout.NewSpacer(),
	)
	return mainContainer
}

func newRectangle() *canvas.Rectangle {
	rect := &canvas.Rectangle{}
	rect.StrokeWidth = 2
	rect.StrokeColor = color.Black
	rect.FillColor = color.Transparent
	return rect
}
