package storage

import (
	"banksystem/internal/model"
	"database/sql"
)

type PendingRegistrationStorage interface {
	Save(pr *model.PendingRegistration) error
	FindByEmail(email string) (*model.PendingRegistration, error)
	DeleteByEmail(email string) error
}

type sqlPendingRegistrationStorage struct {
	db *sql.DB
}

func NewPendingRegistrationStorage(db *sql.DB) PendingRegistrationStorage {
	return &sqlPendingRegistrationStorage{db: db}
}

func (s *sqlPendingRegistrationStorage) Save(pr *model.PendingRegistration) error {
	_, err := s.db.Exec(`
        INSERT OR REPLACE INTO pending_registration
        (name, middlename, surname, password_hash, passport_series, passport_number, phone, email, bank_id, role, verification_code_hash, expires_at)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `,
		pr.Name, pr.Middlename, pr.Surname, pr.PasswordHash,
		pr.PassportSeries, pr.PassportNumber, pr.Phone, pr.Email,
		pr.BankID, pr.Role, pr.VerificationCodeHash, pr.ExpiresAt,
	)
	return err
}

func (s *sqlPendingRegistrationStorage) FindByEmail(email string) (*model.PendingRegistration, error) {
	row := s.db.QueryRow(`
        SELECT name, middlename, surname, password_hash, passport_series, passport_number, phone, email, bank_id, role, verification_code_hash, expires_at
        FROM pending_registration
        WHERE email = ?
    `, email)

	pr := &model.PendingRegistration{}
	err := row.Scan(
		&pr.Name, &pr.Middlename, &pr.Surname, &pr.PasswordHash,
		&pr.PassportSeries, &pr.PassportNumber, &pr.Phone, &pr.Email,
		&pr.BankID, &pr.Role, &pr.VerificationCodeHash, &pr.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

func (s *sqlPendingRegistrationStorage) DeleteByEmail(email string) error {
	_, err := s.db.Exec(`DELETE FROM pending_registration WHERE email = ?`, email)
	return err
}
