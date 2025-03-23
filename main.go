package main

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"banksystem/internal/service"
	"banksystem/internal/storage"
	"banksystem/internal/ui"
	"banksystem/internal/ui/state"
)

func main() {
	db, err := sql.Open("sqlite3", "./banking.db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	userStorage := storage.NewSQLUserStorage(db)
	bankStorage := storage.NewSQLBankStorage(db)
	transactionStorage := storage.NewSQLTransactionStorage(db)

	authService := service.NewAuthService(userStorage)
	bankingService := service.NewBankingService(bankStorage, transactionStorage)

	appState := state.NewAppState()

	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	window := a.NewWindow("Система межбанковских сообщений")
	window.Resize(fyne.NewSize(1024, 800))

	navigationManager := ui.NewNavigationManager(
		a,
		window,
		appState,
		authService,
		bankingService,
	)
	navigationManager.Start()

	window.ShowAndRun()
}
