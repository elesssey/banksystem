package model

type CreditStatus string

const (
	CreditStatusPending   CreditStatus = "pending"
	CreditStatusCompleted CreditStatus = "completed"
	CreditStatusCancelled CreditStatus = "cancelled"
)

type Credit struct {
	Id                  int
	Amount              int
	Term                string
	Сurrency            string
	Status              CreditStatus
	SourceAccountId     int
	SourceBankId        int
	InitiatedByUserId   int
	SourseAccountNumber string
	SourceAccountUser   *User
}
