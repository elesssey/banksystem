package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"banksystem/internal/service"
	"banksystem/internal/storage"
)

type transactionMessage struct {
	ID    int    `json:"id"`
	Event string `json:"event"`
}

type checkOkMessage struct {
	TransactionID    int    `json:"transaction_id"`
	FeeTransactionID int    `json:"fee_transaction_id"`
	Check            string `json:"check"`
}

type feeRequestMessage struct {
	TransactionID            int    `json:"transaction_id"`
	SourceAccountNumber      string `json:"source_account_number"`
	DestinationAccountNumber string `json:"destination_account_number"`
	Amount                   int    `json:"amount"`
}

func main() {
	_ = godotenv.Load("../../.env")

	db, err := sql.Open("sqlite3", "../../banking.db")
	if err != nil {
		log.Fatalf("cannot open db: %v", err)
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

	_ = authService

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatalf("cannot connect to NATS: %v", err)
	}
	defer nc.Close()

	_, err = nc.Subscribe("transactions.created", func(m *nats.Msg) {
		var msg transactionMessage
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("cannot unmarshal message: %v", err)
			return
		}

		log.Printf("got transaction event: id=%d event=%s", msg.ID, msg.Event)

		transaction, err := bankingService.GetTransactionById(msg.ID)
		if err != nil {
			log.Printf("cannot get transaction from db (id=%d): %v", msg.ID, err)
			return
		}

		if transaction.Status != "pending" {
			log.Printf("transaction id=%d has unexpected status=%s, skip further processing", transaction.Id, transaction.Status)
			return
		}

		if transaction.Amount <= 0 {
			log.Printf("transaction id=%d has non-positive amount=%v, skip", transaction.Id, transaction.Amount)

			if err := bankingService.TransactionDeclination(transaction.Id); err != nil {
				log.Printf("cannot decline transaction id=%d: %v", transaction.Id, err)
			}
			return
		}

		userAccount, err := bankingService.GetUserAccount(transaction.InitiatedByUserId, transaction.SourceBankId)

		if err != nil {
			log.Fatalf("cannot get user account from db: %v", err)
			if err := bankingService.TransactionDeclination(transaction.Id); err != nil {
				log.Printf("cannot decline transaction id=%d: %v", transaction.Id, err)
			}
			return
		}

		if userAccount.HoldBalance != float64(transaction.Amount) {
			log.Printf(
				"mismatch hold_balance vs amount for tx id=%d: hold=%v amount=%v, skip",
				transaction.Id,
				userAccount.HoldBalance,
				transaction.Amount,
			)
			if err := bankingService.TransactionDeclination(transaction.Id); err != nil {
				log.Printf("cannot decline transaction id=%d: %v", transaction.Id, err)
			}
			return
		}

		okMsg := checkOkMessage{
			TransactionID:    transaction.Id,
			FeeTransactionID: 0,
			Check:            "validation",
		}

		data, err := json.Marshal(okMsg)
		if err != nil {
			log.Printf("cannot marshal checkOkMessage: %v", err)
			return
		}

		if err := nc.Publish("transactions.check_ok", data); err != nil {
			log.Printf("cannot publish check_ok for tx id=%d: %v", transaction.Id, err)
			return
		}

		feeReq := feeRequestMessage{
			TransactionID:            transaction.Id,
			SourceAccountNumber:      transaction.SourseAccountNumber,
			DestinationAccountNumber: transaction.DestinationAccountNumber,
			Amount:                   transaction.Amount,
		}

		data, err = json.Marshal(feeReq)
		if err != nil {
			log.Printf("cannot marshal feeRequestMessage: %v", err)
			return
		}

		if err := nc.Publish("transactions.fee_request", data); err != nil {
			log.Printf("cannot publish fee_request for tx id=%d: %v", transaction.Id, err)
			return
		}

		log.Printf(
			"published fee_request for tx id=%d: from=%s to=%s amount=%d",
			transaction.Id,
			transaction.SourseAccountNumber,
			transaction.DestinationAccountNumber,
			transaction.Amount,
		)

	})
	if err != nil {
		log.Fatalf("cannot subscribe: %v", err)
	}

	log.Println("transactions-consumer started, listening on transactions.created")
	select {}
}
