package model

type UserAccount struct {
	ID          int
	Number      string
	Balance     float64
	Currency    string
	UserId      int
	BankId      int
	HoldBalance float64
}
