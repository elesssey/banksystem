package screens

import (
	"banksystem/internal/model"
	"fmt"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func MakeAdminMain(credits []*model.Credit, transactions []*model.Transaction, findById func(int) string, trConfirm func(*model.Transaction) error, trDecline func(*model.Transaction) error, crConfirm func(*model.Credit) error, crDecline func(*model.Credit) error) fyne.CanvasObject {

	tabs := container.NewAppTabs(
		container.NewTabItem("ПЕРЕВОДЫ", MakeTransactionTab(transactions, findById, trConfirm, trDecline)),
		container.NewTabItem("ЗАРПЛАТНЫЕ ПРОЕКТЫ", MakeSalaryTab()),
		container.NewTabItem("КРЕДИТЫ", MakeCreditTab(credits, findById, crConfirm, crDecline)),
		container.NewTabItem("РАССРОЧКИ", MakeInstallmentTab()),
	)
	return tabs
}

func MakeTransactionTab(transactions []*model.Transaction, findById func(int) string, trConfirm func(*model.Transaction) error, trDecline func(*model.Transaction) error) fyne.CanvasObject {

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 7
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {
			cellContainer := container.NewCenter(container.New(layout.NewCustomPaddedLayout(20, 20, 20, 20)))
			return cellContainer
		},
		// Обновление содержимого
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			myContainer := cell.(*fyne.Container)

			if id.Row == 0 {
				switch id.Col {
				case 0:
					myContainer.Add(widget.NewLabel("#"))
				case 1:
					myContainer.Add(widget.NewLabel("ФИО ОТПРАВИТЕЛЯ"))
				case 2:
					myContainer.Add(widget.NewLabel("СУММА"))
				case 3:
					myContainer.Add(widget.NewLabel("ФИО ПОЛУЧАТЕЛЯ"))
				case 4:
					myContainer.Add(widget.NewLabel("БАНК"))
				case 5:
					myContainer.Add(widget.NewLabel("ПОДТВЕРЖДЕНИЕ"))

				}

			} else {
				for i, transaction := range transactions {
					if id.Row == i+1 {
						switch id.Col {
						case 0:
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%d", i+1)))
						case 1:
							log.Printf("%v", transaction)
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%s %s %s", transaction.SourceAccountUser.Surname, transaction.SourceAccountUser.Name, transaction.SourceAccountUser.MiddleName)))
						case 2:
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%d", transaction.Amount)))
						case 3:
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%s %s %s", transaction.DestinationAccountUser.Surname, transaction.DestinationAccountUser.Name, transaction.DestinationAccountUser.MiddleName)))
						case 4:
							bank := findById(transaction.DestinationBankId)
							myContainer.Add(widget.NewLabel(bank))
						case 5:
							agreeButton := widget.NewButton("ПОДТВЕРДИТЬ", func() { trConfirm(transaction) })

							disagreeButton := widget.NewButton("ОТКЛОНИТЬ", func() { trDecline(transaction) })

							separator := widget.NewLabel(" | ")
							linkContainer := container.NewHBox(agreeButton, separator, disagreeButton)
							myContainer.Add(linkContainer)
						}
					}

				}
			}

		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 200)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 250)
	table.SetColumnWidth(4, 170)
	table.SetColumnWidth(5, 300)
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

func MakeCreditTab(credits []*model.Credit, findById func(int) string, crConfirm func(*model.Credit) error, crDecline func(*model.Credit) error) fyne.CanvasObject {

	table := widget.NewTable(
		// Размер таблицы
		func() (int, int) {
			return 11, 7
		},
		// Шаблон ячейки
		func() fyne.CanvasObject {
			cellContainer := container.NewCenter(container.New(layout.NewCustomPaddedLayout(20, 20, 20, 20)))
			return cellContainer
		},
		// Обновление содержимого
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			myContainer := cell.(*fyne.Container)

			if id.Row == 0 {
				switch id.Col {
				case 0:
					myContainer.Add(widget.NewLabel("#"))
				case 1:
					myContainer.Add(widget.NewLabel("ПОЛЬЗОВАТЕЛЬ"))
				case 2:
					myContainer.Add(widget.NewLabel("СУММА"))
				case 3:
					myContainer.Add(widget.NewLabel("СРОК(мес)"))
				case 4:
					myContainer.Add(widget.NewLabel("БАНК"))
				case 5:
					myContainer.Add(widget.NewLabel("ПОДТВЕРЖДЕНИЕ"))

				}

			} else {
				for i, credit := range credits {
					if id.Row == i+1 {
						switch id.Col {
						case 0:
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%d", i+1)))
						case 1:
							log.Printf("%v", credit)
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%s %s %s", credit.SourceAccountUser.Surname, credit.SourceAccountUser.Name, credit.SourceAccountUser.MiddleName)))
						case 2:
							myContainer.Add(widget.NewLabel(fmt.Sprintf("%d", credit.Amount)))
						case 3:
							myContainer.Add(widget.NewLabel(fmt.Sprintf(credit.Term)))
						case 4:
							bank := findById(credit.SourceBankId)
							myContainer.Add(widget.NewLabel(bank))
						case 5:
							agreeButton := widget.NewButton("ПОДТВЕРДИТЬ", func() { crConfirm(credit) })

							disagreeButton := widget.NewButton("ОТКЛОНИТЬ", func() { crDecline(credit) })

							separator := widget.NewLabel(" | ")
							linkContainer := container.NewHBox(agreeButton, separator, disagreeButton)
							myContainer.Add(linkContainer)
						}
					}

				}
			}

		},
	)

	table.SetColumnWidth(0, 50)
	table.SetColumnWidth(1, 300)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 200)
	table.SetColumnWidth(5, 300)

	//mainButtons := container.NewGridWithRows(11, MakeButtons(len(transactions), agreeButton, disagreeButton)...)

	//mainContainer := container.NewHBox(table, mainButtons)
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
