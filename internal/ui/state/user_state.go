package state

import (
	"banksystem/internal/model"
)

type UserState struct {
	currentUser         *model.User
	UserCreditList      []*model.Credit
	UserTransactionList []*model.Transaction
}

func NewUserState() *UserState {
	return &UserState{}
}

func (s *UserState) GetCurrentUser() *model.User {
	return s.currentUser
}

func (s *UserState) SetCurrentUser(user *model.User) {
	s.currentUser = user
}

func (s *UserState) ClearCurrentUser() {
	s.currentUser = nil
}
