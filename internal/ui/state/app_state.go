package state

type AppState struct {
	User        *UserState
	Banks       *BanksState
	Transaction *TransactionState
	Credit      *CreditState
}

func NewAppState() *AppState {
	return &AppState{
		User:  NewUserState(),
		Banks: NewBanksState(),
	}
}
