package ui

import (
	"banksystem/internal/model"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeBankSelectorScreen(bank1, bank2, bank3 *model.Bank) fyne.CanvasObject {
	heading := canvas.NewText("Добро пожаловать", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	Label1 := widget.NewLabelWithStyle("Кратко о наших банках:", fyne.TextAlignLeading, fyne.TextStyle{})
	Label2 := widget.NewLabelWithStyle("Выберите банк в котором хотите работать", fyne.TextAlignLeading, fyne.TextStyle{})
	Label3 := widget.NewLabelWithStyle("В любое время вы можете вернуться на эту страницу и выбрать другой банк", fyne.TextAlignLeading, fyne.TextStyle{})

	form := container.New(layout.NewGridLayoutWithRows(2),
		Label2,
		Label3,
	)

	bankCards := container.NewGridWithColumns(3, MakeBankCard(bank1), MakeBankCard(bank2), MakeBankCard(bank3))
	return container.NewVBox(heading, Label1, bankCards, form)
}

func MakeBankCard(bank *model.Bank) fyne.CanvasObject {
	borderHeading := canvas.NewRectangle(color.Black)
	borderHeading.StrokeWidth = 2
	borderHeading.StrokeColor = color.Black
	borderHeading.FillColor = color.Transparent
	headingContent := container.New(
		layout.NewCustomPaddedLayout(5, 5, 0, 0),
		widget.NewLabelWithStyle(bank.Name, fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	heading := container.NewStack(borderHeading, headingContent)

	l1 := widget.NewLabelWithStyle(bank.Descrition, fyne.TextAlignCenter, fyne.TextStyle{})
	l1.Wrapping = fyne.TextWrapWord
	l2 := widget.NewLabelWithStyle("Предприятия которые оперирует банк:", fyne.TextAlignCenter, fyne.TextStyle{})
	text := container.NewVBox(l1, l2)

	borderBody := canvas.NewRectangle(color.Black)
	borderBody.StrokeWidth = 2
	borderBody.StrokeColor = color.Black
	borderBody.FillColor = color.Transparent

	enterpriceList1 := createEnterpricesConvasList(bank.Enterprises)
	enterprises := container.NewGridWithColumns(1, enterpriceList1...)

	button := widget.NewButton("Select", func() {})
	button.Importance = widget.HighImportance
	buttonContainer := container.New(layout.NewCustomPaddedLayout(0, 5, 5, 5), button)

	rating := container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), container.NewCenter(createStarRating(bank.Rating)))
	body := container.NewStack(borderBody, container.NewVBox(text, enterprises, buttonContainer, rating))

	mainContainer := container.NewVBox(heading, body)
	return mainContainer
}

func createEnterpricesConvasList(enterprises []*model.Enterprise) []fyne.CanvasObject {
	list := make([]fyne.CanvasObject, len(enterprises))
	for i, enterprise := range enterprises {
		list[i] = widget.NewLabelWithStyle(enterprise.Name, fyne.TextAlignCenter, fyne.TextStyle{})
	}
	return list
}

func createStarRating(count int) fyne.CanvasObject {
	stars := make([]fyne.CanvasObject, 5)

	for i := 0; i < 5; i++ {
		var star fyne.CanvasObject
		if i < count {
			star = canvas.NewText("★", color.NRGBA{R: 83, G: 82, B: 237, A: 255})
		} else {
			star = canvas.NewText("☆", color.NRGBA{R: 83, G: 82, B: 237, A: 255})
		}
		star.(*canvas.Text).TextSize = 24
		stars[i] = star
	}

	rating := container.NewGridWithColumns(5, stars[0], stars[1], stars[2], stars[3], stars[4])
	return rating
}
