package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"banksystem/internal/model"
	"banksystem/internal/ui/state"
)

func MakeTransactionPage( /*onClickTransaction func (string,string)*/ user *model.User, account *model.UserAccount, bankState *state.BanksState) fyne.CanvasObject {
	heading := canvas.NewText("ПЕРЕВОД", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	fLabel1 := widget.NewLabelWithStyle(fmt.Sprintf("ОТ: %s %s %s", user.Surname, user.Name, user.MiddleName), fyne.TextAlignCenter, fyne.TextStyle{})
	fLabel2 := widget.NewLabelWithStyle(fmt.Sprintf("Номер счета %s", account.Number), fyne.TextAlignLeading, fyne.TextStyle{})
	fLabel3 := widget.NewLabelWithStyle(fmt.Sprintf("БАЛАНС: %.2f", account.Balance), fyne.TextAlignLeading, fyne.TextStyle{})

	tLabel1 := widget.NewLabelWithStyle("Введите  номер  счета  -->", fyne.TextAlignLeading, fyne.TextStyle{})
	tLabel2 := widget.NewLabelWithStyle("Введите сумму перевода -->", fyne.TextAlignLeading, fyne.TextStyle{})

	tEntry1 := widget.NewEntry()
	tEntry2 := widget.NewEntry()

	button := widget.NewButton("ПЕРЕВЕСТИ ДЕНЬГИ", func() { /*onClickTransaction(tEntry1.Text,tEntry2.Text*/ })

	radio := widget.NewRadioGroup(bankState.GetBanksStateNames(), func(value string) {
		bankState.SetTransactionBankByName(value)
	})

	fromBody := container.NewGridWithRows(3, fLabel1, fLabel2, fLabel3)
	toBody := container.NewVBox(
		radio,
		container.New(layout.NewFormLayout(), tLabel1, tEntry1),
		container.New(layout.NewFormLayout(), tLabel2, tEntry2),
	)

	border1 := container.NewStack(newRectangle(), fromBody)
	border2 := container.NewStack(newRectangle(), toBody)

	form := container.NewGridWithColumns(2, border1, border2)

	mainContainer := container.NewGridWithRows(4, layout.NewSpacer(), heading, form, button)

	return mainContainer
}
