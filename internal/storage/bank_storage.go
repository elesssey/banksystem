package storage

import (
	"banksystem/internal/model"
	"database/sql"
	"errors"
	"fmt"
)

type BankStorage interface {
	Fetch(limit int) ([]*model.Bank, error)
	FindUserAccount(user_id int) (*model.User_Account, error)
}

type sqlBankStorage struct {
	db *sql.DB
}

func NewSQLBankStorage(db *sql.DB) BankStorage {
	return &sqlBankStorage{
		db: db,
	}
}

func (s *sqlBankStorage) Fetch(limit int) ([]*model.Bank, error) {
	rows, err := s.db.Query("SELECT id, name, bic, address, description, rating FROM bank ORDER BY id LIMIT ?", limit)
	if err != nil {
		return nil, fmt.Errorf("не получается достать банки rows: %w", err)
	}
	defer rows.Close()

	var banks []*model.Bank
	for rows.Next() {
		bank := &model.Bank{}
		if err := rows.Scan(&bank.ID, &bank.Name, &bank.BIC, &bank.Address, &bank.Descrition, &bank.Rating); err != nil {
			return nil, fmt.Errorf("не получается достать банки scan: %w", err)
		}
		banks = append(banks, bank)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("не получается достать банки rows.Err(): %w", err)
	}

	if len(banks) != 3 {
		return nil, fmt.Errorf("неверное кол-во банков в бд: %w", err)
	}

	eRows, err := s.db.Query("SELECT id, name, unp, address, bank_id FROM enterprise WHERE bank_id IN (?, ?, ?)", banks[0].ID, banks[1].ID, banks[2].ID)
	if err != nil {
		return nil, err
	}
	defer eRows.Close()

	for eRows.Next() {
		enterprise := &model.Enterprise{}
		if err := eRows.Scan(&enterprise.ID, &enterprise.Name, &enterprise.UNP, &enterprise.Address, &enterprise.BankID); err != nil {
			return nil, err
		}

		for bi := range banks {
			if banks[bi].ID == enterprise.BankID {
				banks[bi].Enterprises = append(banks[bi].Enterprises, enterprise)
			}
		}
	}
	if err := eRows.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

func (s *sqlBankStorage) FindUserAccount(user_id int) (*model.User_Account, error) {
	query := `SELECT id, number, balance, currency, user_id, bank_id FROM user_account WHERE user_id = ?`
	row := s.db.QueryRow(query, user_id)
	user_account := &model.User_Account{}
	err := row.Scan(
		&user_account.ID,
		&user_account.Number,
		&user_account.Balance,
		&user_account.Currency,
		&user_account.User_id,
		&user_account.Bank_id,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user account not found")
	}

	if err != nil {
		return nil, err
	}

	return user_account, nil
}
