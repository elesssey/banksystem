package ui

import (
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"banksystem/internal/model"
	"banksystem/internal/ui/state"
)

func MakeTransactionPage(onClickTransaction func(tx model.Transaction) error, onError func(error), user *model.User, account *model.UserAccount, bankState *state.BanksState) fyne.CanvasObject {
	heading := canvas.NewText("ПЕРЕВОД", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	fLabel1 := widget.NewLabelWithStyle(fmt.Sprintf("ОТ: %s %s %s", user.Surname, user.Name, user.MiddleName), fyne.TextAlignCenter, fyne.TextStyle{})
	fLabel2 := widget.NewLabelWithStyle(fmt.Sprintf("Номер счета %s", account.Number), fyne.TextAlignLeading, fyne.TextStyle{})
	fLabel3 := widget.NewLabelWithStyle(fmt.Sprintf("БАЛАНС: %.2f", account.Balance), fyne.TextAlignLeading, fyne.TextStyle{})

	accountLabel := widget.NewLabelWithStyle("Введите  номер  счета  -->", fyne.TextAlignLeading, fyne.TextStyle{})
	amountLabel := widget.NewLabelWithStyle("Введите сумму перевода -->", fyne.TextAlignLeading, fyne.TextStyle{})

	accountEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	button := widget.NewButton("ПЕРЕВЕСТИ ДЕНЬГИ", func() {
		amount, err := strconv.Atoi(amountEntry.Text)
		if err != nil {
			onError(err)
		}
		err = onClickTransaction(model.Transaction{
			SourceAccountId:          account.ID,
			SourceBankId:             bankState.Banks[bankState.SelectedBankIndex].ID,
			DestinationBankId:        bankState.TransactionBankId,
			Amount:                   amount,
			Сurrency:                 account.Currency,
			DestinationAccountNumber: accountEntry.Text,
			InitiatedByUserId:        user.ID,
		})
		if err != nil {
			onError(err)
		}
	})

	radio := widget.NewRadioGroup(bankState.GetBanksStateNames(), func(value string) {
		bankState.SetTransactionBankByName(value)
	})

	fromBody := container.NewGridWithRows(3, fLabel1, fLabel2, fLabel3)
	toBody := container.NewVBox(
		radio,
		container.New(layout.NewFormLayout(), accountLabel, accountEntry),
		container.New(layout.NewFormLayout(), amountLabel, amountEntry),
	)

	border1 := container.NewStack(newRectangle(), fromBody)
	border2 := container.NewStack(newRectangle(), toBody)

	form := container.NewGridWithColumns(2, border1, border2)

	mainContainer := container.NewGridWithRows(4, layout.NewSpacer(), heading, form, button)

	return mainContainer
}
