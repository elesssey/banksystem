package storage

import (
	"database/sql"
	"errors"

	"banksystem/internal/model"

	_ "github.com/mattn/go-sqlite3"
)

type UserStorage interface {
	FindByUsername(username string) (*model.User, error)
}

type sQLiteUserStorage struct {
	db *sql.DB
}

func NewSQLiteUserStorage(db *sql.DB) UserStorage {
	return &sQLiteUserStorage{
		db: db,
	}
}

func (s *sQLiteUserStorage) FindByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, name, surname, passport_series,
			  passport_number, phone, email, role 
              FROM users WHERE username = ?`

	row := s.db.QueryRow(query, username)

	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Name,
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
