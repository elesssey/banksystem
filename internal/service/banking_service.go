package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
	"errors"
	"log"
)

type BankingService interface {
	GetBanks() ([]*model.Bank, error)
	GetUserAccount(userId int, bankId int) (*model.UserAccount, error)
	CreateTransaction(tx *model.Transaction) error
	GetTransactions(bankId int) ([]*model.Transaction, error)
}

type bankingService struct {
	bankStorage        storage.BankStorage
	transactionStorage storage.TransactionStorage
}

func NewBankingService(bankStorage storage.BankStorage, transactionStorage storage.TransactionStorage) BankingService {
	return &bankingService{
		bankStorage:        bankStorage,
		transactionStorage: transactionStorage,
	}
}

func (s *bankingService) GetBanks() ([]*model.Bank, error) {
	return s.bankStorage.Fetch(3)
}

func (s *bankingService) GetTransactions(bankId int) ([]*model.Transaction, error) {
	return s.transactionStorage.Fetch(10, bankId)
}

func (s *bankingService) GetUserAccount(userId, bankId int) (*model.UserAccount, error) {
	return s.bankStorage.FindUserAccount(userId, bankId)
}

func (s *bankingService) CreateTransaction(tx *model.Transaction) error {
	tx.SourceAccountType = model.AccountTypeUser
	tx.DestinationAccountType = model.AccountTypeUser
	tx.Status = model.TransactionStatusPending
	tx.Type = model.TransactionTypeTransfer
	log.Printf("asdasdasdasdasdasd")
	sourseAccount, err := s.bankStorage.FindUserAccountByNumber(tx.SourceBankId, tx.SourseAccountNumber)
	if err != nil {
		return err
	}

	if sourseAccount.Balance*100 < float64(tx.Amount) {
		return errors.New("недостаточно средств")
	}

	destinationAccount, err := s.bankStorage.FindUserAccountByNumber(tx.DestinationBankId, tx.DestinationAccountNumber)
	if err != nil {
		return err
	}
	tx.DestinationAccountId = destinationAccount.ID
	log.Printf("2222222222222222")
	if err := s.bankStorage.CreateTransaction(tx); err != nil {
		return err
	}

	return nil
}
