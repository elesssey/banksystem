package main

import (
	"database/sql"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"banksystem/internal/service"
	"banksystem/internal/state"
	"banksystem/internal/storage"
	"banksystem/internal/ui"
)

func main() {
	db, err := sql.Open("sqlite3", "./banking.db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	userStorage := storage.NewSQLiteUserStorage(db)

	authService := service.NewAuthService(userStorage)

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
	)
	navigationManager.Start()

	window.ShowAndRun()
}
