package screens

import (
	"banksystem/internal/model"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeAdminMain(transactions []*model.Transaction) fyne.CanvasObject {
	//agreeButton := widget.NewButton("+", func() {})
	//disagreeButton := widget.NewButton("-", func() {})

	tabs := container.NewAppTabs(
		container.NewTabItem("ПЕРЕВОДЫ", MakeTransactionTab(transactions)),
		container.NewTabItem("ЗАРПЛАТНЫЕ ПРОЕКТЫ", MakeSalaryTab()),
		container.NewTabItem("КРЕДИТЫ", MakeCreditTab()),
		container.NewTabItem("РАССРОЧКИ", MakeInstallmentTab()),
	)
	return tabs
}

func MakeButtons(length int, button1 fyne.CanvasObject, button2 fyne.CanvasObject) []fyne.CanvasObject {
	buttons := container.NewHBox(button1, button2)
	rowsButtons := make([]fyne.CanvasObject, 0, 10)
	rowsButtons = append(rowsButtons, layout.NewSpacer())
	for i := 0; i < length; i++ {
		rowsButtons = append(rowsButtons, buttons)
	}
	for i := length; i < 10; i++ {
		rowsButtons = append(rowsButtons, layout.NewSpacer())
	}
	return rowsButtons
}

func MakeTransactionTab(transactions []*model.Transaction) fyne.CanvasObject {

	//agreeButton := widget.NewButton("+", func() {})
	//disagreeButton := widget.NewButton("-", func() {})

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 5
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		// Обновление содержимого
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.Alignment = fyne.TextAlignCenter

			if id.Row == 0 {
				switch id.Col {
				case 0:

					label.SetText("#")
				case 1:
					label.SetText("ФИО ОТПРАВИТЕЛЯ")
				case 2:
					label.SetText("СУММА")
				case 3:
					label.SetText("ФИО ПОЛУЧАТЕЛЯ")
				case 4:
					label.SetText("БАНК")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {
				for i /*, transaction */ := range transactions {
					if id.Row == i+1 {
						switch id.Col {
						case 0:
							label.SetText(strconv.Itoa(i + 1))

						case 1:
							//log.Printf("%s %s", transaction.Type, transaction.SourceAccountUser.Name)
							//label.SetText(fmt.Sprintf("%s %s %s", transaction.SourceAccountUser.Surname, transaction.SourceAccountUser.Name, transaction.SourceAccountUser.MiddleName))

						case 2:
						//	label.SetText(string(transaction.Amount))
						case 3:
							//label.SetText(transaction.DestinationAccountUser.Surname + transaction.DestinationAccountUser.Name + transaction.DestinationAccountUser.MiddleName)
						case 4:
							//label.SetText(tra)

						}
					}

				}
			}

		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 170)

	//mainButtons := container.NewGridWithRows(11, MakeButtons(len(transactions), agreeButton, disagreeButton)...)

	//mainContainer := container.NewHBox(table, mainButtons)
	return table
}

func MakeSalaryTab() fyne.CanvasObject {

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 5
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		// Обновление содержимого
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.Alignment = fyne.TextAlignCenter

			if id.Row == 0 {
				switch id.Col {
				case 0:

					label.SetText("#")
				case 1:
					label.SetText("ПРЕДПРИЯТИЕ")
				case 2:
					label.SetText("СУММА")
				case 3:
					label.SetText("ФИО ПОЛУЧАТЕЛЯ")
				case 4:
					label.SetText("БАНК")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {

			}

		},
	)
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 170)

	// mainContainer := container.NewHBox(
	// 	table,
	// 	layout.NewSpacer(),
	// )
	return table
}

func MakeCreditTab() fyne.CanvasObject {

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 5
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		// Обновление содержимого
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.Alignment = fyne.TextAlignCenter

			if id.Row == 0 {
				switch id.Col {
				case 0:

					label.SetText("#")
				case 1:
					label.SetText("ПОЛЬЗОВАТЕЛЬ")
				case 2:
					label.SetText("СУММА")
				case 3:
					label.SetText("СРОК(мес)")
				case 4:
					label.SetText("БАНК")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {

			}

		},
	)
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 170)

	// mainContainer := container.NewHBox(
	// 	table,
	// 	layout.NewSpacer(),
	// )
	return table
}

func MakeInstallmentTab() fyne.CanvasObject {

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 6
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {

			return widget.NewLabel("")
		},
		// Обновление содержимого
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			label.Alignment = fyne.TextAlignCenter

			if id.Row == 0 {
				switch id.Col {
				case 0:

					label.SetText("#")
				case 1:
					label.SetText("ПОЛЬЗОВАТЕЛЬ")
				case 2:
					label.SetText("СУММА(BYN)")
				case 3:
					label.SetText("СРОК(мес)")
				case 4:
					label.SetText("ЕЖЕМЕСЯЧНЫЙ ПЛАТЕЖ(BYN)")
				case 5:
					label.SetText("БАНК")
				}
				label.TextStyle = fyne.TextStyle{Bold: true}
			} else {

			}

		},
	)
	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 270)
	table.SetColumnWidth(5, 170)

	// mainContainer := container.NewHBox(
	// 	table,
	// 	layout.NewSpacer(),
	// )
	return table
}
