package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
)

type BankingService interface {
	GetBanks() ([]*model.Bank, error)
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
