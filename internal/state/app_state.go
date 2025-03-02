package state

type AppState struct {
	User *UserState
}

func NewAppState() *AppState {
	return &AppState{
		User: NewUserState(),
	}
}
