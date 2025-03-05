package state

type AppState struct {
	User  *UserState
	Banks *BanksState
}

func NewAppState() *AppState {
	return &AppState{
		User:  NewUserState(),
		Banks: NewBanksState(),
	}
}
