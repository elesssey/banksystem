package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"

	"banksystem/internal/service"
	"banksystem/internal/storage"
)

type checkOkMessage struct {
	TransactionID    int    `json:"transaction_id"`
	FeeTransactionID int    `json:"fee_transaction_id"`
	Check            string `json:"check"`
}

type txState struct {
	ValidationOK bool
	FeeOK        bool
	BaseTxID     int
	FeeTxID      int
	Timer        *time.Timer
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
	defer nc.Drain()

	states := make(map[int]*txState)
	var mu sync.Mutex

	_, err = nc.Subscribe("transactions.check_ok", func(m *nats.Msg) {
		var msg checkOkMessage
		if err := json.Unmarshal(m.Data, &msg); err != nil {
			log.Printf("orchestrator: cannot unmarshal check_ok: %v", err)
			return
		}

		log.Printf("orchestrator: got check_ok base_tx_id=%d fee_tx_id=%d check=%s",
			msg.TransactionID, msg.FeeTransactionID, msg.Check)

		mu.Lock()
		defer mu.Unlock()

		state, exists := states[msg.TransactionID]
		if !exists {
			state = &txState{}
			states[msg.TransactionID] = state
		}

		switch msg.Check {
		case "validation":
			state.ValidationOK = true
			state.BaseTxID = msg.TransactionID

			if state.Timer == nil {
				txID := msg.TransactionID

				state.Timer = time.AfterFunc(5*time.Second, func() {
					mu.Lock()
					defer mu.Unlock()

					s, ok := states[txID]
					if !ok {
						return
					}

					if !s.FeeOK {
						log.Printf("orchestrator: fee not received in 5s for tx_id=%d, declining base+fee", txID)

						baseTx, err := bankingService.GetTransactionById(s.BaseTxID)
						if err != nil {
							log.Printf("orchestrator: cannot get base tx for decline: %v", err)
						} else if err := transactionStorage.DeclineTransaction(baseTx); err != nil {
							log.Printf("orchestrator: cannot decline base tx: %v", err)
						}

						if s.FeeTxID != 0 {
							feeTx, err := bankingService.GetTransactionById(s.FeeTxID)
							if err != nil {
								log.Printf("orchestrator: cannot get fee tx for decline: %v", err)
							} else if err := transactionStorage.DeclineFeeTransaction(feeTx); err != nil {
								log.Printf("orchestrator: cannot decline fee tx: %v", err)
							}
						}

						delete(states, txID)
					}
				})
			}

		case "fee":
			state.FeeOK = true
			state.FeeTxID = msg.FeeTransactionID

			if state.Timer != nil {
				state.Timer.Stop()
				state.Timer = nil
			}

			baseTx, err := bankingService.GetTransactionById(state.BaseTxID)
			if err != nil {
				log.Printf("orchestrator: cannot get base transaction: %v", err)
				return
			}

			feeTx, err := bankingService.GetTransactionById(state.FeeTxID)
			if err != nil {
				log.Printf("orchestrator: cannot get fee transaction: %v", err)
				return
			}

			log.Printf("orchestrator: before confirm base, tx id=%d status=%s", baseTx.Id, baseTx.Status)
			if err := transactionStorage.ConfirmTransaction(baseTx); err != nil {
				log.Printf("orchestrator: cannot confirm base transaction: %v", err)
				return
			}

			log.Printf("orchestrator: before confirm fee, tx id=%d status=%s", feeTx.Id, feeTx.Status)
			if err := transactionStorage.ConfirmFeeTransaction(feeTx); err != nil {
				log.Printf("orchestrator: cannot confirm fee transaction: %v", err)
				return
			}

			baseTx2, _ := bankingService.GetTransactionById(state.BaseTxID)
			feeTx2, _ := bankingService.GetTransactionById(state.FeeTxID)
			log.Printf("orchestrator: after confirm, base tx id=%d status=%s, fee tx id=%d status=%s",
				baseTx2.Id, baseTx2.Status, feeTx2.Id, feeTx2.Status)

			delete(states, msg.TransactionID)
			log.Printf("orchestrator: validation and fee ok for tx_id=%d", msg.TransactionID)
		}
	})
	if err != nil {
		log.Fatalf("cannot subscribe to transactions.check_ok: %v", err)
	}

	log.Println("transactions-orchestrator started, listening on transactions.check_ok")
	select {}
}
