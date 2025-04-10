package state

import (
	"banksystem/internal/model"
	"math"
)

type CreditState struct {
	ToBanksList [3]*model.Bank
	User        *model.User
	UserAccount *model.UserAccount
	UserBank    *model.Bank
	Amount      float64
	Term        string
}

func NewCreditState(
	toBanksList [3]*model.Bank,
	user *model.User,
	userAccount *model.UserAccount,
	userBank *model.Bank,
) *CreditState {
	return &CreditState{
		ToBanksList: toBanksList,
		User:        user,
		UserAccount: userAccount,
		UserBank:    userBank,
	}
}

func (s *CreditState) BuildCredit() *model.Credit {
	return &model.Credit{
		SourceBankId:        s.UserBank.ID,
		SourceAccountId:     s.UserAccount.ID,
		Amount:              int(math.Floor(s.Amount)),
		Term:                s.Term,
		Сurrency:            s.UserAccount.Currency,
		InitiatedByUserId:   s.User.ID,
		SourseAccountNumber: s.UserAccount.Number,
	}
}
