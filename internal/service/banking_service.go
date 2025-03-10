package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
	"log"
)

type BankingService interface {
	GetBanks() ([]*model.Bank, error)
	GetUserAccount(userId int, bankId int) (*model.UserAccount, error)
	CreateTransaction(tx *model.Transaction) error
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

func (s *bankingService) CreateTransaction(tx *model.Transaction) error {
	tx.SourceAccountType = model.AccountTypeUser
	tx.DestinationAccountType = model.AccountTypeUser
	tx.Status = model.TransactionStatusPending
	tx.Type = model.TransactionTypeTransfer

	log.Printf("CreateTransaction \n Amount %v \n SourceBankId %v \n SourceAccountId %v \n DestinationBankId %v \n DestinationAccountId %v \n InitiatedByUserId %v \n SourceAccountType %v \n DestinationAccountType %v \n Status %v \n Type %v \n DestinationAccountNumber %v \n",
		tx.Amount, tx.SourceBankId, tx.SourceAccountId, tx.DestinationBankId, tx.DestinationAccountId, tx.InitiatedByUserId, tx.SourceAccountType, tx.DestinationAccountType, tx.Status, tx.Type, tx.DestinationAccountNumber)

	// todo: check if source account has enough money

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
