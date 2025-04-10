package main

import (
	"database/sql"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"banksystem/internal/service"
	"banksystem/internal/storage"
	"banksystem/internal/ui"
	"banksystem/internal/ui/state"
)

func main() {
	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Ошибка при открытии файла логов:", err)
	}
	defer logFile.Close()

	// Перенаправляем вывод логов в файл
	log.SetOutput(logFile)

	db, err := sql.Open("sqlite3", "./banking.db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	userStorage := storage.NewSQLUserStorage(db)
	bankStorage := storage.NewSQLBankStorage(db)
	transactionStorage := storage.NewSQLTransactionStorage(db)
	creditStorage := storage.NewSQLTCreditStorage(db)

	authService := service.NewAuthService(userStorage)
	bankingService := service.NewBankingService(bankStorage, transactionStorage, creditStorage)

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
