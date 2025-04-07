package state

import "banksystem/internal/model"

type BanksState struct {
	SelectedBankIndex     int
	WorkingAccount        *model.UserAccount
	TransactionBankId     int
	IsInitialized         bool
	BanksList             [3]*model.Bank
	AdminTransactionsList []*model.Transaction
}

func NewBanksState() *BanksState {
	return &BanksState{}
}

func (s *BanksState) GetCurrentBank() *model.Bank {
	return s.BanksList[s.SelectedBankIndex]
}

func (s *BanksState) SetBanks(banks []*model.Bank) {
	for i := range 3 {
		s.BanksList[i] = banks[i]
	}
	s.IsInitialized = true
}

func (s *BanksState) FindBankNameById(id int) string {
	for _, bank := range s.BanksList {
		if bank.ID == id {
			return bank.Name
		}
	}
	return ""
}

func (s *BanksState) GetBankStateNames() []string {
	bankNames := []string{}
	for _, bank := range s.BanksList {
		bankNames = append(bankNames, bank.Name)
	}
	return bankNames
}
