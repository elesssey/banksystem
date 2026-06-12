package web

import "banksystem/internal/ui/state"

type StateStore struct {
	data map[int]*state.AppState
}

func NewStateStore() *StateStore {
	return &StateStore{
		data: make(map[int]*state.AppState),
	}
}

func (s *StateStore) Get(userID int) *state.AppState {
	st, ok := s.data[userID]
	if !ok {
		st = state.NewAppState()
		s.data[userID] = st
	}
	return st
}
