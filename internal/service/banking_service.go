package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
)

type BankingService interface {
	GetBanks() ([]*model.Bank, error)
	GetUserAccount(userId int, bankId int) (*model.UserAccount, error)
	CreateTransaction(tx model.Transaction) error
}

type bankingService struct {
	bankStorage storage.BankStorage
}

func NewBankingService(bankStorage storage.BankStorage) BankingService {
	return &bankingService{
		bankStorage: bankStorage,
	}
}

func (s *bankingService) GetBanks() ([]*model.Bank, error) {
	return s.bankStorage.Fetch(3)
}

func (s *bankingService) GetUserAccount(userId, bankId int) (*model.UserAccount, error) {
	return s.bankStorage.FindUserAccount(userId, bankId)
}

func (s *bankingService) CreateTransaction(tx model.Transaction) error {
	tx.SourceAccountType = model.AccountTypeUser
	tx.DestinationAccountType = model.AccountTypeUser
	tx.Status = model.TransactionStatusPending
	tx.Type = model.TransactionTypeTransfer

	destinationAccount, err := s.bankStorage.FindUserAccountByNumber(tx.DestinationBankId, tx.DestinationAccountNumber)
	if err != nil {
		return err
	}
	tx.DestinationAccountId = destinationAccount.ID

	if err := s.bankStorage.CreateTransaction(tx); err != nil {
		return err
	}

	return nil
}
