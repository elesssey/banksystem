package screens

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeLoginScreen(onLoginTry func(string, string), onRegistrationStartClick func()) fyne.CanvasObject {
	heading := canvas.NewText("BANK-SYSTEM", color.Black)
	heading.TextSize = 30
	heading.Alignment = fyne.TextAlignCenter
	heading.TextStyle.Bold = true

	emailLabel := widget.NewLabelWithStyle("EMAIL:         ", fyne.TextAlignTrailing, fyne.TextStyle{})
	emailEntry := widget.NewEntry()

	passwordLabel := widget.NewLabelWithStyle("PASSWORD:", fyne.TextAlignTrailing, fyne.TextStyle{})
	passwordEntry := widget.NewPasswordEntry()

	signInButton := widget.NewButton("SIGN IN", func() { onLoginTry(emailEntry.Text, passwordEntry.Text) })
	signInButton.Importance = widget.HighImportance

	orLabel := canvas.NewText("OR", color.Black)
	orLabel.Alignment = fyne.TextAlignCenter
	orLabel.TextSize = 20
	orLabel.TextStyle.Bold = true

	orContainer := container.NewCenter(
		container.New(layout.NewHBoxLayout(),
			container.New(layout.NewPaddedLayout(), orLabel),
		),
	)

	signUpText := widget.NewLabel("don't have an account?")
	signUpButton := widget.NewButton("SIGN UP", func() { onRegistrationStartClick() })
	signUpButton.Importance = widget.MediumImportance

	signUpContainer := container.New(layout.NewGridLayout(2),
		signUpText,
		signUpButton,
	)
	signUpWrapper := container.New(layout.NewCenterLayout(), signUpContainer)

	form := container.NewVBox(
		container.New(layout.NewFormLayout(), emailLabel, emailEntry),
		container.New(layout.NewFormLayout(), passwordLabel, passwordEntry),
	)

	content := container.New(layout.NewCustomPaddedLayout(100, 100, 300, 300), container.NewVBox(
		layout.NewSpacer(),
		container.New(layout.NewCenterLayout(), heading),
		layout.NewSpacer(),
		form,
		container.New(layout.NewStackLayout(), signInButton),
		layout.NewSpacer(),
		orContainer,
		layout.NewSpacer(),
		signUpWrapper,
		layout.NewSpacer(),
	))

	return container.New(layout.NewGridLayout(1), content)
}
