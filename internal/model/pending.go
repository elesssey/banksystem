package model

import "time"

type PendingRegistration struct {
	Name                 string
	Middlename           string
	Surname              string
	PasswordHash         string
	PassportSeries       string
	PassportNumber       string
	Phone                string
	Email                string
	BankID               int
	Role                 Role
	VerificationCodeHash string
	ExpiresAt            time.Time
}
