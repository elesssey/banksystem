package main

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"banksystem/internal/model"
	"banksystem/internal/service"
	"banksystem/internal/storage"
)

type feeRequestMessage struct {
	TransactionID            int    `json:"transaction_id"`
	SourceAccountNumber      string `json:"source_account_number"`
	DestinationAccountNumber string `json:"destination_account_number"`
	Amount                   int    `json:"amount"`
}

type checkOkMessage struct {
	TransactionID    int    `json:"transaction_id"`
	FeeTransactionID int    `json:"fee_transaction_id"`
	Check            string `json:"check"`
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

	_, err = nc.Subscribe("transactions.fee_request", func(m *nats.Msg) {
		var req feeRequestMessage
		if err := json.Unmarshal(m.Data, &req); err != nil {
			log.Printf("cannot unmarshal fee_request: %v", err)
			return
		}

		log.Printf(
			"fee-service: got fee_request tx_id=%d from=%s to=%s amount=%d",
			req.TransactionID,
			req.SourceAccountNumber,
			req.DestinationAccountNumber,
			req.Amount,
		)

		var fee int
		var bankId int
		var managerId int

		transaction, err := bankingService.GetTransactionById(req.TransactionID)

		if err != nil {
			log.Fatalf("cannot get  transaction by id: %v", err)
			return
		}

		if transaction.SourceBankId == transaction.DestinationBankId {
			fee = 0
		} else {
			switch transaction.DestinationBankId {
			case 1:
				fee = transaction.Amount * 5 / 100
				bankId = 1
				managerId = 101
			case 2:
				fee = transaction.Amount * 6 / 100
				bankId = 2
				managerId = 103
			case 3:
				fee = transaction.Amount * 7 / 100
				bankId = 3
				managerId = 105
			}
		}

		bankAccount, err := bankingService.GetUserAccount(managerId, bankId)
		if err != nil {
			log.Printf("fee-service: cannot get bank account: %v", err)
			return
		}

		admin, err := bankingService.GetUserById(bankAccount.UserId)
		if err != nil {
			log.Printf("fee-service: cannot get bank admin: %v", err)
			return
		}

		sourceAccount, err := bankingService.GetUserAccount(transaction.InitiatedByUserId, transaction.SourceBankId)
		if err != nil {
			log.Printf("fee-service: cannot get source account user: %v", err)
			return
		}

		if sourceAccount.Balance < float64(fee+transaction.Amount) {
			log.Printf("fee-service: cannot make transaction, not enough money: %v", err)
			return
		}

		newFeeTransaction := &model.Transaction{
			Amount:                   fee,
			Сurrency:                 transaction.Сurrency,
			Status:                   "pending",
			SourceAccountId:          transaction.SourceAccountId,
			DestinationAccountId:     bankAccount.ID,
			SourceAccountType:        "user",
			DestinationAccountType:   "user",
			Type:                     "transfer",
			SourceBankId:             transaction.SourceBankId,
			DestinationBankId:        transaction.DestinationBankId,
			InitiatedByUserId:        transaction.InitiatedByUserId,
			SourseAccountNumber:      transaction.SourseAccountNumber,
			DestinationAccountNumber: bankAccount.Number,
			SourceAccountUser:        transaction.SourceAccountUser,
			DestinationAccountUser:   admin,
		}

		err = bankStorage.CreateFeeTransaction(newFeeTransaction)
		if err != nil {
			log.Printf("fee-service: cannot create fee transaction: %v", err)
			return
		}

		log.Printf(
			"fee-service: calculated fee for tx_id=%d: amount=%d fee=%d",
			newFeeTransaction.Id,
			req.Amount,
			fee,
		)

		okMsg := checkOkMessage{
			TransactionID:    req.TransactionID,
			FeeTransactionID: newFeeTransaction.Id,
			Check:            "fee",
		}

		data, err := json.Marshal(okMsg)
		if err != nil {
			log.Printf("fee-service: cannot marshal checkOkMessage: %v", err)
			return
		}

		if err := nc.Publish("transactions.check_ok", data); err != nil {
			log.Printf("fee-service: cannot publish fee check_ok for tx id=%d: %v", req.TransactionID, err)
			return
		}

		log.Printf("fee-service: published fee check_ok for tx id=%d", req.TransactionID)
	})
	if err != nil {
		log.Fatalf("cannot subscribe to transactions.fee_request: %v", err)
	}

	log.Println("fee-service started, listening on transactions.fee_request")
	select {}
}
