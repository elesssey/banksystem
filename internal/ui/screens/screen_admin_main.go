package screens

import (
	"banksystem/internal/ui/state"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func MakeAdminMain(transactionState *state.TransactionState) fyne.CanvasObject {
	tabs := container.NewAppTabs(
		container.NewTabItem("ПЕРЕВОДЫ", MakeTransactionTab()),
		container.NewTabItem("ЗАРПЛАТНЫЕ ПРОЕКТЫ", MakeSalaryTab()),
		container.NewTabItem("КРЕДИТЫ", MakeCreditTab()),
		container.NewTabItem("РАССРОЧКИ", MakeInstallmentTab()),
	)
	return tabs
}

func MakeTransactionTab() fyne.CanvasObject {

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 5
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {
			return widget.NewLabel("utiiutitk")
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

			}

		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 170)
	// tableContainer := container.NewVBox(
	// 	container.NewHBox(table, layout.NewSpacer()),
	// 	layout.NewSpacer(),
	// )

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
