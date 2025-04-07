package screens

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"banksystem/internal/ui/state"
)

func MakeTransactionScreen(createTransaction func() error, onError func(error), txState *state.TransactionState) fyne.CanvasObject {
	heading := canvas.NewText("ПЕРЕВОД", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	senderFullNameLabel := widget.NewLabelWithStyle(fmt.Sprintf("ОТ: %s %s %s", txState.Sender.Surname, txState.Sender.Name, txState.Sender.MiddleName), fyne.TextAlignCenter, fyne.TextStyle{})
	senderAccountLabel := widget.NewLabelWithStyle(fmt.Sprintf("Номер счета %s", txState.SenderAccount.Number), fyne.TextAlignLeading, fyne.TextStyle{})
	balanceLabel := widget.NewLabelWithStyle(fmt.Sprintf("БАЛАНС: %.2f", txState.SenderAccount.Balance), fyne.TextAlignLeading, fyne.TextStyle{})

	accountLabel := widget.NewLabelWithStyle("Введите  номер  счета  -->", fyne.TextAlignLeading, fyne.TextStyle{})
	amountLabel := widget.NewLabelWithStyle("Введите сумму перевода -->", fyne.TextAlignLeading, fyne.TextStyle{})

	accountEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	button := widget.NewButton("ПЕРЕВЕСТИ ДЕНЬГИ", func() {
		if txState.ReceiverBank == nil {
			onError(errors.New("вы не выбрали банк получателя"))
			return
		}
		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			onError(err)
		}
		txState.Amount = amount
		log.Println("Account number: ", accountEntry.Text)
		txState.ReceiverAccountNumber = accountEntry.Text
		err = createTransaction()
		if err != nil {
			onError(err)
		}
	})

	radio := widget.NewRadioGroup(txState.GetBanksStateNames(), func(value string) {
		txState.SetTransactionBankByName(value)
	})

	fromBody := container.NewGridWithRows(3, senderFullNameLabel, senderAccountLabel, balanceLabel)
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
