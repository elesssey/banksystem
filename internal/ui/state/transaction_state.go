package state

type TransactionState struct {
	Amount float64
	Number string
	BankId int
}

func NewTransactionState() *TransactionState {
	return &TransactionState{}
}
