package state

import "banksystem/internal/model"

type BanksState struct {
	SelectedBankIndex int
	IsInitialized     bool
	Banks             [3]*model.Bank
}

func NewBanksState() *BanksState {
	return &BanksState{}
}

func (s *BanksState) SetBanks(bank1, bank2, bank3 *model.Bank) {
	s.Banks[0] = bank1
	s.Banks[1] = bank2
	s.Banks[2] = bank3
	s.IsInitialized = true
}
