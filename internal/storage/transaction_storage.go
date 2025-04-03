package storage

import (
	"banksystem/internal/model"
	"database/sql"
	"fmt"
)

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
	FetchCurrentTransaction(id int) (*model.Transaction, error)
	ConfirmTransaction(transaction *model.Transaction) error
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

	transactionId := make([]any, 0, len(transactions)*2)
	for _, tx := range transactions {
		transactionId = append(transactionId, tx.SourceAccountId, tx.DestinationAccountId)
	}
	fmt.Println(transactionId)

	queryStart := `SELECT 
		user.id,
		user.name, 
		user.middlename, 
		user.surname, 
		user.passport_series,
		user.passport_number,
		user.phone,
		user.email,
		user_account.id 
	FROM user JOIN user_account ON user.id = user_account.user_id 
	WHERE user_account.id IN (`
	queryMiddle := ``
	queryEnd := `)`

	for i := range transactionId {
		if i < len(transactionId)-1 {
			queryMiddle += `?,`

		} else {
			queryMiddle += `?`
		}
	}
	erows, err := s.db.Query(queryStart+queryMiddle+queryEnd, transactionId...)
	if err != nil {
		return nil, fmt.Errorf("не получается достать аккаунты rows: %w", err)
	}

	for erows.Next() {
		var userWithAccountId struct {
			ID             int
			Name           string
			MiddleName     string
			Surname        string
			PassportSeries string
			PassportNumber string
			Phone          string
			Email          string
			AccountId      int
		}
		if err := erows.Scan(
			&userWithAccountId.ID,
			&userWithAccountId.Name,
			&userWithAccountId.MiddleName,
			&userWithAccountId.Surname,
			&userWithAccountId.PassportSeries,
			&userWithAccountId.PassportNumber,
			&userWithAccountId.Phone,
			&userWithAccountId.Email,
			&userWithAccountId.AccountId,
		); err != nil {
			return nil, fmt.Errorf("не получается достать пользователей с их аккаунтами scan: %w", err)
		}
		for _, transaction := range transactions {
			modelUser := model.User{
				ID:             userWithAccountId.ID,
				Name:           userWithAccountId.Name,
				MiddleName:     userWithAccountId.MiddleName,
				Surname:        userWithAccountId.Surname,
				PassportSeries: userWithAccountId.PassportSeries,
				PassportNumber: userWithAccountId.PassportNumber,
				Phone:          userWithAccountId.Phone,
				Email:          userWithAccountId.Email,
			}
			if transaction.DestinationAccountId == userWithAccountId.AccountId {
				transaction.DestinationAccountUser = &modelUser
			} else if transaction.SourceAccountId == userWithAccountId.AccountId {
				transaction.SourceAccountUser = &modelUser
			}
		}
	}
	if err := erows.Err(); err != nil {
		return nil, fmt.Errorf("не получается достать пользователй с их аккаунтами rows.Err(): %w", err)
	}

	return transactions, nil
}

func (s *sqlTransactionStorage) FetchCurrentTransaction(id int) (*model.Transaction, error) {
	query := `SELECT id,amount,currency,
		description,status,source_account_id,
		destination_account_id,source_account_type,
		destination_account_type,type,source_bank_id,
		destination_bank_id,initiated_by_user_id
              FROM system_transaction WHERE id = ?`

	row := s.db.QueryRow(query, id)

	transaction := &model.Transaction{}
	err := row.Scan(
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
	)
	if err != nil {
		return nil, fmt.Errorf("не получается достать транзакцию scan: %w", err)
	}
	return transaction, nil
}

func (s *sqlTransactionStorage) ConfirmTransaction(transaction *model.Transaction) error {
	dbtx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = dbtx.Exec(`UPDATE user_account SET hold_balance =  hold_balance - ? WHERE id =?`, transaction.Amount, transaction.SourceAccountId)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	_, err = dbtx.Exec(`UPDATE user_account SET balance = balance + ? WHERE id =?`, transaction.Amount, transaction.DestinationAccountId)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	err = dbtx.Commit()
	if err != nil {
		return err
	}
	return nil
}
