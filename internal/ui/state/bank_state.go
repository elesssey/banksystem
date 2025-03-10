package state

import "banksystem/internal/model"

type BanksState struct {
	SelectedBankIndex int
	WorkingAccount    *model.UserAccount
	TransactionBankId int
	IsInitialized     bool
	BanksList         [3]*model.Bank
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
