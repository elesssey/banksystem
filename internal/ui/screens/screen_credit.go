package screens

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"banksystem/internal/ui/state"
)

func MakeCreditScreen(createTransaction func() error, onError func(error), crState *state.CreditState) fyne.CanvasObject {
	heading := canvas.NewText("Оформление кредита", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	senderFullNameLabel := widget.NewLabelWithStyle(fmt.Sprintf("ОТ: %s %s %s", crState.User.Surname, crState.User.Name, crState.User.MiddleName), fyne.TextAlignCenter, fyne.TextStyle{})
	senderAccountLabel := widget.NewLabelWithStyle(fmt.Sprintf("Номер счета %s", crState.UserAccount.Number), fyne.TextAlignLeading, fyne.TextStyle{})
	balanceLabel := widget.NewLabelWithStyle(fmt.Sprintf("БАЛАНС: %.2f", crState.UserAccount.Balance), fyne.TextAlignLeading, fyne.TextStyle{})

	termLabel := widget.NewLabelWithStyle("Выберите срок(мес) для кредита -->", fyne.TextAlignLeading, fyne.TextStyle{})
	amountLabel := widget.NewLabelWithStyle("Введите сумму кредита -->", fyne.TextAlignLeading, fyne.TextStyle{})

	//accountEntry := widget.NewEntry()
	amountEntry := widget.NewEntry()

	button := widget.NewButton("ВЗЯТЬ КРЕДИТ", func() {
		if crState.Term == "" {
			onError(errors.New("вы не выбрали срок"))
			return
		}
		amount, err := strconv.ParseFloat(amountEntry.Text, 64)
		if err != nil {
			onError(err)
		}
		crState.Amount = amount
		err = createTransaction()
		if err != nil {
			onError(err)
		}
	})
	terms := []string{"3", "6", "12"}
	radio := widget.NewRadioGroup(terms, func(value string) {
		crState.Term = value
	})

	fromBody := container.NewGridWithRows(3, senderFullNameLabel, senderAccountLabel, balanceLabel)
	toBody := container.NewVBox(
		termLabel,
		radio,
		container.New(layout.NewFormLayout(), amountLabel, amountEntry),
	)

	border1 := container.NewStack(newRectangle(), fromBody)
	border2 := container.NewStack(newRectangle(), toBody)

	form := container.NewGridWithColumns(2, border1, border2)

	mainContainer := container.NewGridWithRows(4, layout.NewSpacer(), heading, form, button)

	return mainContainer
}
