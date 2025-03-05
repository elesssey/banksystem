package storage

import (
	"database/sql"
	"errors"

	"banksystem/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type UserStorage interface {
	FindByEmail(email string) (*model.User, error)
}

type sqlUserStorage struct {
	db *sql.DB
}

func NewSQLUserStorage(db *sql.DB) UserStorage {
	return &sqlUserStorage{
		db: db,
	}
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
