package model

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusCompleted TransactionStatus = "completed"
	TransactionStatusCancelled TransactionStatus = "cancelled"
)

type AccountType string

const (
	AccountTypeUser       AccountType = "user"
	AccountTypeEnterprise AccountType = "enterprise"
)

type TransactionType string

const (
	TransactionTypeTransfer TransactionType = "transfer"
	TransactionTypeSalary   TransactionType = "salary"
)

type Transaction struct {
	Id                       int
	Amount                   int
	Сurrency                 string
	Description              string
	Status                   TransactionStatus
	SourceAccountId          int
	DestinationAccountId     int
	SourceAccountType        AccountType
	DestinationAccountType   AccountType
	Type                     TransactionType
	SourceBankId             int
	DestinationBankId        int
	InitiatedByUserId        int
	DestinationAccountNumber string
}
