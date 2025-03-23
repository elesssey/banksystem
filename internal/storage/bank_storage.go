package storage

import (
	"banksystem/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var ErrAccountNotFound = errors.New("account not found")

type BankStorage interface {
	Fetch(limit int) ([]*model.Bank, error)
	FindUserAccount(user_id int, bankId int) (*model.UserAccount, error)
	FindUserAccountByNumber(bankId int, number string) (*model.UserAccount, error)
	CreateTransaction(tx *model.Transaction) error
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

func (s *sqlBankStorage) FindUserAccount(userId int, bankId int) (*model.UserAccount, error) {
	query := `SELECT id, number, balance, currency, user_id, bank_id FROM user_account WHERE user_id = ? AND bank_id = ?`
	row := s.db.QueryRow(query, userId, bankId)
	userAccount := &model.UserAccount{}
	err := row.Scan(
		&userAccount.ID,
		&userAccount.Number,
		&userAccount.Balance,
		&userAccount.Currency,
		&userAccount.UserId,
		&userAccount.BankId,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return userAccount, nil
}

func (s *sqlBankStorage) FindUserAccountByNumber(bankId int, number string) (*model.UserAccount, error) {
	query := `SELECT id, number, balance, currency, user_id, bank_id FROM user_account WHERE number = ? AND bank_id = ?`
	row := s.db.QueryRow(query, number, bankId)
	userAccount := &model.UserAccount{}
	err := row.Scan(
		&userAccount.ID,
		&userAccount.Number,
		&userAccount.Balance,
		&userAccount.Currency,
		&userAccount.UserId,
		&userAccount.BankId,
	)

	if err == sql.ErrNoRows {
		return nil, ErrAccountNotFound
	} else if err != nil {
		return nil, err
	}

	return userAccount, nil
}

func (s *sqlBankStorage) CreateTransaction(tx *model.Transaction) error {
	log.Printf("CreateTransaction \n Amount %v \n SourceBankId %v \n SourceAccountId %v \n DestinationBankId %v \n DestinationAccountId %v \n InitiatedByUserId %v \n SourceAccountType %v \n DestinationAccountType %v \n Status %v \n Type %v \n DestinationAccountNumber %v \n",
		tx.Amount, tx.SourceBankId, tx.SourceAccountId, tx.DestinationBankId, tx.DestinationAccountId, tx.InitiatedByUserId, tx.SourceAccountType, tx.DestinationAccountType, tx.Status, tx.Type, tx.DestinationAccountNumber)
	// elisey todo: insert into transaction table

	sql := `INSERT INTO system_transaction(amount,currency,description,status,source_account_id,destination_account_id,source_account_type,destination_account_type,type,source_bank_id,destination_bank_id,initiated_by_user_id)
	VALUES (?,?,?,?,?,?,?,?,?,?,?,?)`
	_, err := s.db.Exec(sql, tx.Amount, tx.Сurrency, tx.Description, tx.Status, tx.SourceAccountId, tx.DestinationAccountId, tx.SourceAccountType, tx.DestinationAccountType, tx.Type, tx.SourceBankId, tx.DestinationBankId, tx.InitiatedByUserId)
	if err != nil {
		return err
	}
	return nil

}
