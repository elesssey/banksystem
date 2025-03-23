package state

import (
	"banksystem/internal/model"
	"math"
)

type TransactionState struct {
	ToBanksList           [3]*model.Bank
	Sender                *model.User
	SenderAccount         *model.UserAccount
	SenderBank            *model.Bank
	ReceiverBank          *model.Bank
	ReceiverAccountNumber string
	Amount                float64
}

func NewTransactionState(
	toBanksList [3]*model.Bank,
	sender *model.User,
	senderAccount *model.UserAccount,
	senderBank *model.Bank,
) *TransactionState {
	return &TransactionState{
		ToBanksList:   toBanksList,
		Sender:        sender,
		SenderAccount: senderAccount,
		SenderBank:    senderBank,
	}
}

func (s *TransactionState) GetBanksStateNames() []string {
	bankNames := []string{}
	for _, bank := range s.ToBanksList {
		bankNames = append(bankNames, bank.Name)
	}
	return bankNames
}

func (s *TransactionState) SetTransactionBankByName(bankName string) {
	for _, bank := range s.ToBanksList {
		if bank.Name == bankName {
			s.ReceiverBank = bank
			return
		}
	}
}

func (s *TransactionState) BuildTransaction() *model.Transaction {
	return &model.Transaction{
		SourceBankId:             s.SenderBank.ID,
		SourceAccountId:          s.SenderAccount.ID,
		DestinationBankId:        s.ReceiverBank.ID,
		Amount:                   int(math.Floor(s.Amount)),
		Сurrency:                 s.SenderAccount.Currency,
		InitiatedByUserId:        s.Sender.ID,
		SourseAccountNumber:      s.SenderAccount.Number,
		DestinationAccountNumber: s.ReceiverAccountNumber,
	}
}
