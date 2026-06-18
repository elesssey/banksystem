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
	GetUserById(userId int) (*model.User, error)
	CreateTransaction(tx *model.Transaction) error
	CreationTransaction(transaction *model.Transaction) error
	CreateCredit(cr *model.Credit) error
	GetTransactions(bankId int) ([]*model.Transaction, error)
	GetTransactionById(transactionId int) (*model.Transaction, error)
	GetAllTransactions() ([]*model.Transaction, error)
	GetCredits(bankid int) ([]*model.Credit, error)
	TransactionConfirmation(id int) error
	TransactionDeclination(id int) error
	CreditConfirmation(id int) error
	CreditDeclination(id int) error
	FreezeAccount(id int) error
	UnFreezeAccount(id int) error
	GetAnyUserAccount(userId int) (*model.UserAccount, error)
	GetAccountByNumber(number string) (*model.UserAccount, error)
	MarkValidationOk(transactionId int) error
	MarkFeeOk(transactionId int) error
	GetChecks(transactionId int) (*model.TransactionChecks, error)
	//ConfirmTransaction(amount int, number string) error
}

type bankingService struct {
	bankStorage        storage.BankStorage
	transactionStorage storage.TransactionStorage
	creditStorage      storage.CreditStorage
	userStorage        storage.UserStorage
	pendingStorage     storage.PendingRegistrationStorage
	pendingService     EmailService
}

func NewBankingService(bankStorage storage.BankStorage, transactionStorage storage.TransactionStorage, creditStorage storage.CreditStorage, userStorage storage.UserStorage, pendingStorage storage.PendingRegistrationStorage, pendingService EmailService) BankingService {
	return &bankingService{
		bankStorage:        bankStorage,
		transactionStorage: transactionStorage,
		creditStorage:      creditStorage,
		userStorage:        userStorage,
		pendingStorage:     pendingStorage,
		pendingService:     pendingService,
	}
}

func (s *bankingService) GetBanks() ([]*model.Bank, error) {
	return s.bankStorage.Fetch(3)
}

func (s *bankingService) GetTransactions(bankId int) ([]*model.Transaction, error) {
	return s.transactionStorage.FetchwithUsers(10, bankId)
}
func (s *bankingService) GetTransactionById(transactionId int) (*model.Transaction, error) {
	return s.transactionStorage.FetchTransaction(transactionId)
}

func (s *bankingService) GetAllTransactions() ([]*model.Transaction, error) {
	return s.transactionStorage.FetchAllTransaction()
}

func (s *bankingService) GetUserById(userId int) (*model.User, error) {
	return s.userStorage.FindById(userId)
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

func (s *bankingService) FreezeAccount(id int) error {
	err := s.bankStorage.FreezeAccount(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *bankingService) UnFreezeAccount(id int) error {
	err := s.bankStorage.UnFreezeAccount(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *bankingService) GetAnyUserAccount(userID int) (*model.UserAccount, error) {
	return s.userStorage.FindAnyAccountByUserID(userID)
}

func (s *bankingService) GetAccountByNumber(number string) (*model.UserAccount, error) {
	return s.userStorage.FindAccountByNumber(number)
}

func (s *bankingService) CreationTransaction(transaction *model.Transaction) error {
	err := s.CreateTransaction(transaction)
	if err != nil {
		return err
	}

	err = s.pendingService.SendVerification("cchheltyyy81@gmail.com", transaction)

	return nil
}

func (s *bankingService) MarkValidationOk(transactionId int) error {
	return s.transactionStorage.SetValidationTransaction(transactionId)
}

func (s *bankingService) MarkFeeOk(transactionId int) error {
	return s.transactionStorage.SetFeeOk(transactionId)
}

func (s *bankingService) GetChecks(transactionId int) (*model.TransactionChecks, error) {
	return s.transactionStorage.FetchChecks(transactionId)
}

//func (s *bankingService) ConfirmTransaction(amount int, number string) error {}
