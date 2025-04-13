package storage

import (
	"database/sql"
	"errors"
	"log"

	"math/big"
	"math/rand"
	"time"

	"banksystem/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type UserStorage interface {
	FindByEmail(email string) (*model.User, error)
	FindById(id int) (*model.User, error)
	AddNewUserWithAccount(user *model.User, bankId int) error
	AddNewAccountToUser(bankId, userId int) (*model.UserAccount, error)
}

type sqlUserStorage struct {
	db *sql.DB
}

func NewSQLUserStorage(db *sql.DB) UserStorage {
	return &sqlUserStorage{
		db: db,
	}
}
func (s *sqlUserStorage) FindById(id int) (*model.User, error) {
	query := `SELECT id, password, name, middlename, surname, passport_series,
			  passport_number, phone, email, role
              FROM user WHERE id = ?`
	row := s.db.QueryRow(query, id)
	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Password,
		&user.Name,
		&user.MiddleName,
		&user.Surname,
		&user.PassportSeries,
		&user.PassportNumber,
		&user.Phone,
		&user.Email,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *sqlUserStorage) FindByEmail(email string) (*model.User, error) {
	query := `SELECT id, password, name, middlename, surname, passport_series,
			  passport_number, phone, email, role
              FROM user WHERE email = ?`

	row := s.db.QueryRow(query, email)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Password,
		&user.Name,
		&user.MiddleName,
		&user.Surname,
		&user.PassportSeries,
		&user.PassportNumber,
		&user.Phone,
		&user.Email,
		&user.Role,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *sqlUserStorage) AddNewUserWithAccount(user *model.User, bankId int) error {
	var numberStart string
	switch bankId {
	case 1:
		numberStart = "BY86AKBB30141"
	case 2:
		numberStart = "BY20BAPB30141"
	case 3:
		numberStart = "BY13PJCB30141"
	}

	mainNumber := numberStart + randomLargeNumberString()

	dbtx, err := s.db.Begin()
	if err != nil {
		return err
	}

	row, err := dbtx.Exec(
		`INSERT INTO user(
			name,
			middlename,
			surname,
			password,
			passport_series,
			passport_number,
			phone,
			email,
			role 
		)VALUES(?,?,?,?,?,?,?,?,?)`, user.Name, user.MiddleName, user.Surname, user.Password, user.PassportSeries, user.PassportNumber, user.Phone, user.Email, user.Role)

	if err != nil {
		dbtx.Rollback()
		return err
	}

	index, err := row.LastInsertId()
	userId := int(index)

	log.Printf("%s %d %d", mainNumber, userId, bankId)
	if err != nil {
		return err
	}
	_, err = dbtx.Exec(
		`INSERT INTO user_account(
			number,
			balance,
			currency,
			user_id,
			bank_id
		)VALUES(?,?,?,?,?)`, mainNumber, 0, "BYN", userId, bankId)

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

func randomLargeNumberString() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := big.NewInt(100000000000000)
	n := new(big.Int)
	n.Rand(r, max)

	return n.String()
}

func (s *sqlUserStorage) AddNewAccountToUser(bankId int, userId int) (*model.UserAccount, error) {
	var numberStart string
	switch bankId {
	case 1:
		numberStart = "BY86AKBB30141"
	case 2:
		numberStart = "BY20BAPB30141"
	case 3:
		numberStart = "BY13PJCB30141"
	}

	mainNumber := numberStart + randomLargeNumberString()

	dbtx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = dbtx.Exec(`INSERT INTO user_account(
		number,
		balance,
		currency,
		user_id,
		bank_id
	)VALUES(?,?,?,?,?)`, mainNumber, 0, "BYN", userId, bankId)

	if err != nil {
		dbtx.Rollback()
		return nil, err
	}

	row := dbtx.QueryRow(`SELECT id, number,balance,currency,user_id,bank_id,hold_balance FROM user_account WHERE user_id=? AND bank_id=?`, userId, bankId)

	userAccount := &model.UserAccount{}
	err = row.Scan(
		&userAccount.ID,
		&userAccount.Number,
		&userAccount.Balance,
		&userAccount.Currency,
		&userAccount.UserId,
		&userAccount.BankId,
		&userAccount.HoldBalance,
	)

	if err != nil {
		dbtx.Rollback()
		return nil, err
	}

	err = dbtx.Commit()
	if err != nil {
		return nil, err
	}

	return userAccount, nil
}
