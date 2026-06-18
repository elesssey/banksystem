package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"banksystem/internal/service"
	"banksystem/internal/storage"
	webui "banksystem/internal/web"
)

func main() {

	_ = godotenv.Load("../../.env")

	db, err := sql.Open("sqlite3", "../../banking.db")
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	userStorage := storage.NewSQLUserStorage(db)
	bankStorage := storage.NewSQLBankStorage(db)
	transactionStorage := storage.NewSQLTransactionStorage(db)
	creditStorage := storage.NewSQLTCreditStorage(db)
	pendingStorage := storage.NewPendingRegistrationStorage(db)

	pendingService := service.NewEmailService()
	authService := service.NewAuthService(userStorage, pendingStorage, pendingService)
	bankingService := service.NewBankingService(bankStorage, transactionStorage, creditStorage, userStorage, pendingStorage, pendingService)

	tmpl := template.Must(template.ParseGlob("../../web/templates/*.html"))

	natsConn, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("cannot connect to NATS: %v", err)
	}
	defer natsConn.Drain()

	producer := service.NewNATSTransactionProducer(natsConn)

	handler := webui.NewHandler(
		bankingService,
		authService,
		tmpl,
		producer,
	)

	r := chi.NewRouter()

	r.Get("/login", handler.LoginPage)
	r.Post("/login", handler.Login)

	r.Get("/register", handler.RegisterPage)
	r.Post("/register", handler.StartRegistration)

	r.Get("/verify-email", handler.VerifyEmailPage)
	r.Post("/verify-email", handler.VerifyEmail)

	r.Group(func(r chi.Router) {
		r.Use(handler.AuthMiddleware)
		r.Get("/transaction", handler.TransactionPage)
		r.Post("/transaction", handler.MakeTransaction)
		r.Get("/admin", handler.AdminPage)
		r.Get("/transaction/Confirm", handler.MakeTransactionConfirmation)
		r.Get("/transaction/Cancel", handler.MakeTransactionDeclanation)
		//r.Get("/transactions", handler.TransactionsPage)
		//r.Post("/transfer", handler.TransferMoney)
	})
	//r.Get("/transaction/Confirm", handler.MakeTransactionConfirmation)
	//r.Get("/transaction/Cancel", handler.MakeTransactionDeclanation)

	log.Println("web server started on :8180")
	err = http.ListenAndServe(":8180", r)
	if err != nil {
		log.Fatal(err)
	}
}
