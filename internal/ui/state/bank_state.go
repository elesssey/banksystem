package state

import "banksystem/internal/model"

type BanksState struct {
	SelectedBankIndex int
	TransactionBankId int
	IsInitialized     bool
	Banks             [3]*model.Bank
}

func NewBanksState() *BanksState {
	return &BanksState{}
}
func (s *BanksState) GetBanksStateNames() []string {
	bankNames := []string{s.Banks[0].Name, s.Banks[1].Name, s.Banks[2].Name}
	return bankNames
}
func (s *BanksState) SetTransactionBankByName(bankName string) {

}

func (s *BanksState) SetBanks(bank1, bank2, bank3 *model.Bank) {
	s.Banks[0] = bank1
	s.Banks[1] = bank2
	s.Banks[2] = bank3
	s.IsInitialized = true
}
