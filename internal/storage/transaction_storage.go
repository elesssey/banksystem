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
