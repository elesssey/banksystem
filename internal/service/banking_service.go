package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
	"errors"
	"log"
)

type BankingService interface {
	GetBanks() ([]*model.Bank, error)
	GetUserAccount(userId int, bankId int) (*model.UserAccount, error)
	CreateTransaction(tx *model.Transaction) error
	CreateCredit(cr *model.Credit) error
	GetTransactions(bankId int) ([]*model.Transaction, error)
	GetCredits(bankid int) ([]*model.Credit, error)
	TransactionConfirmation(id int) error
	TransactionDeclination(id int) error
	CreditConfirmation(id int) error
	CreditDeclination(id int) error
}

type bankingService struct {
	bankStorage        storage.BankStorage
	transactionStorage storage.TransactionStorage
	creditStorage      storage.CreditStorage
}

func NewBankingService(bankStorage storage.BankStorage, transactionStorage storage.TransactionStorage, creditStorage storage.CreditStorage) BankingService {
	return &bankingService{
		bankStorage:        bankStorage,
		transactionStorage: transactionStorage,
		creditStorage:      creditStorage,
	}
}

func (s *bankingService) GetBanks() ([]*model.Bank, error) {
	return s.bankStorage.Fetch(3)
}

func (s *bankingService) GetTransactions(bankId int) ([]*model.Transaction, error) {
	return s.transactionStorage.FetchwithUsers(10, bankId)
}

func (s *bankingService) GetCredits(bankId int) ([]*model.Credit, error) {
	return s.creditStorage.FetchCreditwithUsers(10, bankId)
}

func (s *bankingService) GetUserAccount(userId, bankId int) (*model.UserAccount, error) {
	return s.bankStorage.FindUserAccount(userId, bankId)
}

func (s *bankingService) CreateTransaction(tx *model.Transaction) error {
	tx.SourceAccountType = model.AccountTypeUser
	tx.DestinationAccountType = model.AccountTypeUser
	tx.Status = model.TransactionStatusPending
	tx.Type = model.TransactionTypeTransfer
	sourseAccount, err := s.bankStorage.FindUserAccountByNumber(tx.SourceBankId, tx.SourseAccountNumber)
	if err != nil {
		return err
	}

	if sourseAccount.Balance < float64(tx.Amount) {
		return errors.New("недостаточно средств")
	}

	sourseAccount.HoldBalance = float64(tx.Amount)

	destinationAccount, err := s.bankStorage.FindUserAccountByNumber(tx.DestinationBankId, tx.DestinationAccountNumber)
	if err != nil {
		return err
	}
	tx.DestinationAccountId = destinationAccount.ID
	if err := s.bankStorage.CreateTransaction(tx); err != nil {
		return err
	}

	return nil
}

func (s *bankingService) TransactionConfirmation(id int) error {
	log.Printf("dasdasd")
	transaction, err := s.transactionStorage.FetchCurrentTransaction(id)
	if err != nil {
		return err
	}
	log.Printf("%d %d %d", transaction.SourceAccountId, transaction.SourceBankId, transaction.Amount)
	sourceAccount, err := s.bankStorage.FindUserAccountByAccountId(transaction.SourceBankId, transaction.SourceAccountId)
	if err != nil {
		return err
	}
	if transaction.Status != model.TransactionStatusPending {
		return errors.New("транзакция уже подтверждена")
	}

	log.Printf("%f %f", sourceAccount.Balance, sourceAccount.HoldBalance)
	if sourceAccount.HoldBalance < float64(transaction.Amount) {
		return errors.New("не достаточно средств у отправителя")
	}
	log.Printf("1111")
	err = s.transactionStorage.ConfirmTransaction(transaction)
	if err != nil {
		return err
	}

	return nil
}

func (s *bankingService) TransactionDeclination(id int) error {
	transaction, err := s.transactionStorage.FetchCurrentTransaction(id)
	if err != nil {
		return err
	}
	err = s.transactionStorage.DeclineTransaction(transaction)
	if err != nil {
		return err
	}
	return nil
}

func (s *bankingService) CreditConfirmation(id int) error {
	credit, err := s.creditStorage.FetchCurrentCredit(id)
	if err != nil {
		return err
	}

	sourceAccount, err := s.bankStorage.FindUserAccountByAccountId(credit.SourceBankId, credit.SourceAccountId)
	if err != nil {
		return err
	}
	if credit.Status != model.CreditStatusPending {
		return errors.New("транзакция уже подтверждена")
	}

	log.Printf("%f %f", sourceAccount.Balance, sourceAccount.HoldBalance)

	err = s.creditStorage.ConfirmCredit(credit)
	if err != nil {
		return err
	}
	return nil
}

func (s *bankingService) CreditDeclination(id int) error {
	credit, err := s.creditStorage.FetchCurrentCredit(id)
	if err != nil {
		return err
	}
	err = s.creditStorage.DeclineCredit(credit)
	if err != nil {
		return err
	}
	return nil
}

func (s *bankingService) CreateCredit(cr *model.Credit) error {
	cr.Status = model.CreditStatusPending
	sourseAccount, err := s.bankStorage.FindUserAccountByNumber(cr.SourceBankId, cr.SourseAccountNumber)
	if err != nil {
		return err
	}

	sourseAccount.HoldBalance = float64(cr.Amount)

	if err := s.bankStorage.CreateCredit(cr); err != nil {
		return err
	}

	return nil
}
