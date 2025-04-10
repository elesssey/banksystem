package screens

import (
	"banksystem/internal/model"
	"banksystem/internal/ui/state"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeBankScreen(onCreditClick func(), onTransactionClick func(), banksState *state.BanksState, user *model.User) fyne.CanvasObject {
	heading := widget.NewLabelWithStyle(banksState.GetCurrentBank().Name, fyne.TextAlignLeading, fyne.TextStyle{})

	existingAccountsLabel := canvas.NewText("Существующие счета:", color.Black)
	transferInstructionLabel := canvas.NewText("Вы можете перести деньги другому пользователю используя дебетовый счет", color.RGBA{255, 0, 0, 255})
	debitAccountLabel := widget.NewLabelWithStyle("Дебетовый счет", fyne.TextAlignCenter, fyne.TextStyle{})
	accountNumberLabel := canvas.NewText("Номер счета:", color.Black)
	balanceLabel := canvas.NewText("Баланс:", color.Black)
	debitNumber := widget.NewLabelWithStyle(banksState.WorkingAccount.Number, fyne.TextAlignCenter, fyne.TextStyle{})
	debitText := fmt.Sprintf("%.2f %s", banksState.WorkingAccount.Balance, banksState.WorkingAccount.Currency)
	debitBalance := widget.NewLabelWithStyle(debitText, fyne.TextAlignCenter, fyne.TextStyle{})
	savingLabel := widget.NewLabelWithStyle("Накопительный счет", fyne.TextAlignCenter, fyne.TextStyle{})
	creditLabel := widget.NewLabelWithStyle("Кредиторный счет", fyne.TextAlignCenter, fyne.TextStyle{})

	infButton := widget.NewButton("Посмотреть информацию", func() {})
	infButton.Resize(fyne.NewSize(200, 50))
	transferButton := widget.NewButton("СДЕЛАТЬ ПЕРЕВОД", func() { onTransactionClick() })
	openButton := widget.NewButton("ОТКРЫТЬ СЧЕТ", func() { onCreditClick() })

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
		debitAccountLabel,
		container.NewHBox(accountNumberLabel, debitNumber),
		container.NewHBox(balanceLabel, debitBalance),
		container.NewHBox(layout.NewSpacer(), infButton, layout.NewSpacer()),
		layout.NewSpacer(),
	)

	savingBody := container.NewVBox(
		savingLabel,
		container.NewHBox(layout.NewSpacer(), imgPlus, layout.NewSpacer()),
	)
	creditBody := container.NewVBox(
		creditLabel,
		container.NewHBox(layout.NewSpacer(), openButton, layout.NewSpacer()),
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
		existingAccountsLabel,
		layout.NewSpacer(),
		accounts,
		layout.NewSpacer(),
		transferInstructionLabel,
		layout.NewSpacer(),
		transferButton,
		layout.NewSpacer(),
	)
	return mainContainer
}
