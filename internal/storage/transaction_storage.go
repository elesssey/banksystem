package storage

import (
	"banksystem/internal/model"
	"database/sql"
	"fmt"
)

type UserWithAccount struct {
	ID             int
	Password       string
	Name           string
	MiddleName     string
	Surname        string
	PassportSeries string
	PassportNumber string
	Phone          string
	Email          string
	AccountNumber  string
}

const fetchQuery = `
	SELECT 
		id,
		amount,
		currency,
		description,
		status,
		source_account_id,
		destination_account_id,
		source_account_type,
		destination_account_type,
		type,
		source_bank_id,
		destination_bank_id,
		initiated_by_user_id
	FROM system_transaction WHERE source_bank_id = ? ORDER BY id DESC LIMIT ?
`

type TransactionStorage interface {
	Fetch(limit int, bankId int) ([]*model.Transaction, error)
	FetchwithUsers(limit int, bankId int) ([]*model.Transaction, error)
}

type sqlTransactionStorage struct {
	db *sql.DB
}

func NewSQLTransactionStorage(db *sql.DB) TransactionStorage {
	return &sqlTransactionStorage{
		db: db,
	}
}

func (s *sqlTransactionStorage) Fetch(limit int, bankId int) ([]*model.Transaction, error) {
	rows, err := s.db.Query(fetchQuery, bankId, limit)
	if err != nil {
		return nil, fmt.Errorf("не получается достать транзакции rows: %w", err)
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		transaction := &model.Transaction{}
		if err := rows.Scan(
			&transaction.Id,
			&transaction.Amount,
			&transaction.Сurrency,
			&transaction.Description,
			&transaction.Status,
			&transaction.SourceAccountId,
			&transaction.DestinationAccountId,
			&transaction.SourceAccountType,
			&transaction.DestinationAccountType,
			&transaction.Type,
			&transaction.SourceBankId,
			&transaction.DestinationBankId,
			&transaction.InitiatedByUserId,
		); err != nil {
			return nil, fmt.Errorf("не получается достать транзакции scan: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("не получается достать транзакции rows.Err(): %w", err)
	}
	return transactions, nil
}

func (s *sqlTransactionStorage) FetchwithUsers(limit int, bankId int) ([]*model.Transaction, error) {
	transactions, err := s.Fetch(limit, bankId)
	if err != nil {
		return nil, fmt.Errorf("не получается достать транзакции rows: %w", err)
	}
	transactionNumbers := make([]any, 0, len(transactions)*2)

	for _, tx := range transactions {
		transactionNumbers = append(transactionNumbers, tx.SourseAccountNumber, tx.DestinationAccountNumber)
	}

	queryStart := `SELECT 
		user.id,
		user.name, 
		user.middlename, 
		user.surname, 
		user.password,
		user.passport_series,
		user.passport_number,
		user.phone,
		user.email,
		user.role,	
		user_account.number 
	FROM user JOIN user_account ON user.id = user_account.user_id 
	WHERE user_account.number IN (`
	queryMiddle := ``
	queryEnd := `)`

	for i := range transactionNumbers {
		if i < len(transactionNumbers)-1 {
			queryMiddle += `?,`

		} else {
			queryMiddle = `?`
		}
	}
	erows, err := s.db.Query(queryStart+queryMiddle+queryEnd, transactionNumbers...)
	if err != nil {
		return nil, fmt.Errorf("не получается достать аккаунты rows: %w", err)
	}
	//запрос будет с inner join надо соединить таблицу аккаунтов юзеров и моделей юзеров по user.Id и фильтровать по аккаунту

	for erows.Next() {
		user := &UserWithAccount{}
		if err := erows.Scan(
			&user.ID,
			&user.Name,
			&user.MiddleName,
			&user.Surname,
			&user.PassportSeries,
			&user.PassportNumber,
			&user.Phone,
			&user.Email,
			&user.AccountNumber,
		); err != nil {
			return nil, fmt.Errorf("не получается достать пользователей с их аккаунтами scan: %w", err)
		}
		for _, transaction := range transactions {
			modelUser := model.User{
				ID:             user.ID,
				Name:           user.Name,
				MiddleName:     user.MiddleName,
				Surname:        user.Surname,
				PassportSeries: user.PassportSeries,
				PassportNumber: user.PassportNumber,
				Phone:          user.Phone,
				Email:          user.Email,
			}
			if transaction.DestinationAccountNumber == user.AccountNumber {
				transaction.DestinationAccountUser = &modelUser
			} else if transaction.SourseAccountNumber == user.AccountNumber {
				transaction.SourceAccountUser = &modelUser
			}
		}
	}
	if err := erows.Err(); err != nil {
		return nil, fmt.Errorf("не получается достать пользователй с их аккаунтами rows.Err(): %w", err)
	}

	return transactions, nil
}
