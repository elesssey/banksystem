package storage

import (
	"banksystem/internal/model"
	"database/sql"
	"fmt"
)

const fetchCreditQuery = `
	SELECT 
		id,
		amount,
		term,
		currency,
		status,
		source_account_id,
		source_bank_id,
		initiated_by_user_id
	FROM system_credit WHERE source_bank_id = ? ORDER BY id DESC LIMIT ?
`

type CreditStorage interface {
	FetchCredit(limit int, bankId int) ([]*model.Credit, error)
	FetchCreditwithUsers(limit int, bankId int) ([]*model.Credit, error)
	FetchCurrentCredit(id int) (*model.Credit, error)
	ConfirmCredit(credit *model.Credit) error
	DeclineCredit(credit *model.Credit) error
}

type sqlCreditStorage struct {
	db *sql.DB
}

func NewSQLTCreditStorage(db *sql.DB) CreditStorage {
	return &sqlCreditStorage{
		db: db,
	}
}

func (s *sqlCreditStorage) FetchCredit(limit int, bankId int) ([]*model.Credit, error) {
	rows, err := s.db.Query(fetchCreditQuery, bankId, limit)
	if err != nil {
		return nil, fmt.Errorf("не получается достать кредиты rows: %w", err)
	}
	defer rows.Close()

	var credits []*model.Credit
	for rows.Next() {
		credit := &model.Credit{}
		if err := rows.Scan(
			&credit.Id,
			&credit.Amount,
			&credit.Term,
			&credit.Сurrency,
			&credit.Status,
			&credit.SourceAccountId,
			&credit.SourceBankId,
			&credit.InitiatedByUserId,
		); err != nil {
			return nil, fmt.Errorf("не получается достать кредиты scan: %w", err)
		}
		credits = append(credits, credit)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("не получается достать кредиты rows.Err(): %w", err)
	}
	return credits, nil
}

func (s *sqlCreditStorage) FetchCreditwithUsers(limit int, bankId int) ([]*model.Credit, error) {
	credits, err := s.FetchCredit(limit, bankId)
	if err != nil {
		return nil, fmt.Errorf("не получается достать кредиты rows: %w", err)
	}

	creditId := make([]any, 0, len(credits))
	for _, cr := range credits {
		creditId = append(creditId, cr.SourceAccountId)
	}
	fmt.Println(creditId)

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

	for i := range creditId {
		if i < len(creditId)-1 {
			queryMiddle += `?,`

		} else {
			queryMiddle += `?`
		}
	}
	erows, err := s.db.Query(queryStart+queryMiddle+queryEnd, creditId...)
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
		for _, credit := range credits {
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
			if credit.SourceAccountId == userWithAccountId.AccountId {
				credit.SourceAccountUser = &modelUser
			}
		}
	}
	if err := erows.Err(); err != nil {
		return nil, fmt.Errorf("не получается достать пользователй с их аккаунтами rows.Err(): %w", err)
	}

	return credits, nil
}

func (s *sqlCreditStorage) FetchCurrentCredit(id int) (*model.Credit, error) {
	query := `SELECT id,amount,currency,term,
	status,source_account_id,source_bank_id,
	initiated_by_user_id
    FROM system_credit WHERE id = ?`

	row := s.db.QueryRow(query, id)

	credit := &model.Credit{}
	err := row.Scan(
		&credit.Id,
		&credit.Amount,
		&credit.Term,
		&credit.Сurrency,
		&credit.Status,
		&credit.SourceAccountId,
		&credit.SourceBankId,
		&credit.InitiatedByUserId,
	)
	if err != nil {
		return nil, fmt.Errorf("не получается достать транзакцию scan: %w", err)
	}
	return credit, nil
}

func (s *sqlCreditStorage) ConfirmCredit(credit *model.Credit) error {
	dbtx, err := s.db.Begin()
	if err != nil {
		return err
	}
	_, err = dbtx.Exec(`UPDATE user_account SET hold_balance =  hold_balance - ? WHERE id =?`, credit.Amount, credit.SourceAccountId)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	_, err = dbtx.Exec(`UPDATE user_account SET balance = balance + ? WHERE id =?`, credit.Amount, credit.SourceAccountId)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	_, err = dbtx.Exec(`UPDATE system_credit SET status = ? WHERE id =?`, "completed", credit.Id)
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
func (s *sqlCreditStorage) DeclineCredit(credit *model.Credit) error {
	dbtx, err := s.db.Begin()
	if err != nil {
		return err
	}

	_, err = dbtx.Exec(`UPDATE user_account SET hold_balance =  hold_balance - ? WHERE id =?`, credit.Amount, credit.SourceAccountId)
	if err != nil {
		dbtx.Rollback()
		return err
	}

	_, err = dbtx.Exec(`UPDATE system_credit SET status = ? WHERE id =?`, "cancelled", credit.Id)
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
