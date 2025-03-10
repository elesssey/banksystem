package state

type AppState struct {
	User        *UserState
	Banks       *BanksState
	Transaction *TransactionState
}

func NewAppState() *AppState {
	return &AppState{
		User:  NewUserState(),
		Banks: NewBanksState(),
	}
}
