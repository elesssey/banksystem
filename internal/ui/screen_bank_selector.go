package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var enterpriceList1 = []fyne.CanvasObject{
	widget.NewLabelWithStyle("Enterprise 1", fyne.TextAlignCenter, fyne.TextStyle{}),
	widget.NewLabelWithStyle("Enterprise 2", fyne.TextAlignCenter, fyne.TextStyle{}),
	widget.NewLabelWithStyle("Enterprise 3", fyne.TextAlignCenter, fyne.TextStyle{}),
	widget.NewLabelWithStyle("Enterprise 4", fyne.TextAlignCenter, fyne.TextStyle{}),
	widget.NewLabelWithStyle("Enterprise 5", fyne.TextAlignCenter, fyne.TextStyle{}),
}
var enterpriceList2 = []fyne.CanvasObject{
	widget.NewLabelWithStyle("Enterprise 1", fyne.TextAlignCenter, fyne.TextStyle{}),
	widget.NewLabelWithStyle("Enterprise 2", fyne.TextAlignCenter, fyne.TextStyle{}),
	widget.NewLabelWithStyle("Enterprise 3", fyne.TextAlignCenter, fyne.TextStyle{}),
}

func MakeBankSelectorScreen() fyne.CanvasObject {
	heading := canvas.NewText("Welcome to our bank system", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	Label1 := widget.NewLabelWithStyle("Short information about our banks:", fyne.TextAlignLeading, fyne.TextStyle{})
	Label2 := widget.NewLabelWithStyle("Select one bank, if yoy want to work with him", fyne.TextAlignLeading, fyne.TextStyle{})
	Label3 := widget.NewLabelWithStyle("After that, you can return and use the services of other banks", fyne.TextAlignLeading, fyne.TextStyle{})

	form := container.New(layout.NewGridLayoutWithRows(2),
		Label2,
		Label3,
	)
	bankCards := container.NewGridWithColumns(3, MakeBankCard(), MakeBankCard(), MakeBankCard())
	return container.NewVBox(heading, Label1, bankCards, form)
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

func MakeBankCard() fyne.CanvasObject {
	borderHeading := canvas.NewRectangle(color.Black)
	borderHeading.StrokeWidth = 2
	borderHeading.StrokeColor = color.Black
	borderHeading.FillColor = color.Transparent
	headingContent := container.New(
		layout.NewCustomPaddedLayout(5, 5, 0, 0),
		widget.NewLabelWithStyle("Альфа банк", fyne.TextAlignCenter, fyne.TextStyle{}),
	)
	heading := container.NewStack(borderHeading, headingContent)

	l1 := widget.NewLabelWithStyle("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.", fyne.TextAlignCenter, fyne.TextStyle{})
	l1.Wrapping = fyne.TextWrapWord
	l2 := widget.NewLabelWithStyle("Enterprises which cooperate with bank:", fyne.TextAlignCenter, fyne.TextStyle{})
	text := container.NewVBox(l1, l2)

	borderBody := canvas.NewRectangle(color.Black)
	borderBody.StrokeWidth = 2
	borderBody.StrokeColor = color.Black
	borderBody.FillColor = color.Transparent

	enterpriceList1 := container.NewGridWithRows(len(enterpriceList1), enterpriceList1...)
	enterpriceList2 := container.NewGridWithRows(len(enterpriceList2), enterpriceList2...)
	enterprises := container.NewGridWithColumns(2, enterpriceList1, enterpriceList2)

	button := widget.NewButton("Select", func() {})
	button.Importance = widget.HighImportance
	buttonContainer := container.New(layout.NewCustomPaddedLayout(0, 5, 5, 5), button)

	rating := container.New(layout.NewCustomPaddedLayout(0, 0, 10, 10), container.NewCenter(createStarRating(4)))
	body := container.NewStack(borderBody, container.NewVBox(text, enterprises, buttonContainer, rating))

	mainContainer := container.NewVBox(heading, body)
	return mainContainer
}
