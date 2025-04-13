package screens

import (
	"banksystem/internal/model"
	"banksystem/internal/ui/state"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeBankSelectorScreen(onBankClick func(int), banksState *state.BanksState) fyne.CanvasObject {
	heading := canvas.NewText("Добро пожаловать", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	shortAboutBanks := widget.NewLabelWithStyle("Кратко о наших банках:", fyne.TextAlignLeading, fyne.TextStyle{})
	selectTheBank := widget.NewLabelWithStyle("Выберите банк в котором хотите работать", fyne.TextAlignLeading, fyne.TextStyle{})
	bankChange := widget.NewLabelWithStyle("В любое время вы можете вернуться на эту страницу и выбрать другой банк", fyne.TextAlignLeading, fyne.TextStyle{})

	form := container.New(layout.NewGridLayoutWithRows(2), selectTheBank, bankChange)

	bankCards := []fyne.CanvasObject{}
	for i, bank := range banksState.BanksList {
		bankCards = append(bankCards, MakeBankCard(onBankClick, bank, i))
	}

	return container.NewVBox(
		heading, shortAboutBanks,
		container.NewGridWithColumns(3, bankCards...),
		form,
	)
}

func MakeBankCard(onBankClick func(int), bank *model.Bank, index int) fyne.CanvasObject {
	borderHeading := canvas.NewRectangle(color.Black)
	borderHeading.StrokeWidth = 2
	borderHeading.StrokeColor = color.Black
	borderHeading.FillColor = color.Transparent
	headingContent := container.New(
		layout.NewCustomPaddedLayout(5, 5, 0, 0),
		widget.NewLabelWithStyle(bank.Name, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	heading := container.NewStack(borderHeading, headingContent)

	bankDescription := widget.NewLabelWithStyle(bank.Descrition, fyne.TextAlignCenter, fyne.TextStyle{})
	bankDescription.Wrapping = fyne.TextWrapWord
	bankEnterprises := widget.NewLabelWithStyle("Предприятия которые оперирует банк:", fyne.TextAlignCenter, fyne.TextStyle{})
	text := container.NewVBox(bankDescription, bankEnterprises)

	borderBody := canvas.NewRectangle(color.Black)
	borderBody.StrokeWidth = 2
	borderBody.StrokeColor = color.Black
	borderBody.FillColor = color.Transparent

	enterprises := container.NewGridWithColumns(1, createEnterpricesConvasList(bank.Enterprises)...)

	button := widget.NewButton("Select", func() {
		onBankClick(index)
	})
	button.Importance = widget.HighImportance
	buttonContainer := container.New(layout.NewCustomPaddedLayout(0, 5, 5, 5), button)

	rating := container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), container.NewCenter(createStarRating(bank.Rating)))
	body := container.NewStack(borderBody, container.NewVBox(text, enterprises, buttonContainer, rating))

	return container.NewVBox(heading, body)
}

func createEnterpricesConvasList(enterprises []*model.Enterprise) []fyne.CanvasObject {
	list := make([]fyne.CanvasObject, len(enterprises))
	for i, enterprise := range enterprises {
		list[i] = widget.NewLabelWithStyle(enterprise.Name, fyne.TextAlignCenter, fyne.TextStyle{})
	}
	return list
}

func createStarRating(count int) fyne.CanvasObject {
	stars := make([]fyne.CanvasObject, 0, 5)

	for i := range 5 {
		var char string
		if i < count {
			char = "★"
		} else {
			char = "☆"
		}
		star := canvas.NewText(char, color.NRGBA{R: 83, G: 82, B: 237, A: 255})
		star.TextSize = 24
		stars = append(stars, star)
	}

	return container.NewGridWithColumns(5, stars...)
}
