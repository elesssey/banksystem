package state

type AppState struct {
	User         *UserState
	Banks        *BanksState
	Transactions *TransactionState
}

func NewAppState() *AppState {
	return &AppState{
		User:         NewUserState(),
		Banks:        NewBanksState(),
		Transactions: NewTransactionState(),
	}
}
