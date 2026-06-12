package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidEmail = errors.New("invalid email")
var ErrInvalidCode = errors.New("invalid verification code")
var ErrCodeExpired = errors.New("verification code expired")

type AuthService interface {
	Login(email, password string) (*model.User, error)
	Registrate(user *model.User, bankId int) error
	CreateAccount(bankId, userID int) (*model.UserAccount, error)

	StartRegistration(user *model.User, bankId int) error
	ConfirmRegistration(email, code string) error
}

type authService struct {
	userStorage    storage.UserStorage
	pendingStorage storage.PendingRegistrationStorage
	pendingService EmailService
}

func NewAuthService(userStorage storage.UserStorage, pendingStorage storage.PendingRegistrationStorage, pendingService EmailService) *authService {
	return &authService{
		userStorage:    userStorage,
		pendingStorage: pendingStorage,
		pendingService: pendingService,
	}
}

func (s *authService) Login(email, password string) (*model.User, error) {
	if email == "" || password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userStorage.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	stored := user.Password

	if len(stored) == len(password) {
		if stored != password {
			return nil, ErrInvalidCredentials
		}
		return user, nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *authService) Registrate(user *model.User, bankId int) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashed)

	err = s.userStorage.AddNewUserWithAccount(user, bankId)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) CreateAccount(bankId int, userId int) (*model.UserAccount, error) {
	newAccount, err := s.userStorage.AddNewAccountToUser(bankId, userId)
	if err != nil {
		return nil, err
	}
	return newAccount, nil
}

func (s *authService) StartRegistration(user *model.User, bankId int) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	code, err := s.pendingService.GenerateVerificationCode()
	if err != nil {
		return err
	}

	hashedCode, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	pending := &model.PendingRegistration{
		Name:                 user.Name,
		Middlename:           user.MiddleName,
		Surname:              user.Surname,
		PasswordHash:         string(hashedPassword),
		PassportSeries:       user.PassportSeries,
		PassportNumber:       user.PassportNumber,
		Phone:                user.Phone,
		Email:                user.Email,
		BankID:               bankId,
		Role:                 "client",
		VerificationCodeHash: string(hashedCode),
		ExpiresAt:            time.Now().Add(10 * time.Minute),
	}

	if err := s.pendingStorage.Save(pending); err != nil {
		return err
	}

	return s.pendingService.SendVerificationCode(user.Email, code)
}

func (s *authService) ConfirmRegistration(email, code string) error {
	pending, err := s.pendingStorage.FindByEmail(email)
	if err != nil {
		return err
	}

	if time.Now().After(pending.ExpiresAt) {
		s.pendingStorage.DeleteByEmail(email)
		return ErrCodeExpired
	}

	if err := bcrypt.CompareHashAndPassword([]byte(pending.VerificationCodeHash), []byte(code)); err != nil {
		s.pendingStorage.DeleteByEmail(email)
		return ErrInvalidCode
	}

	user := &model.User{
		Name:           pending.Name,
		MiddleName:     pending.Middlename,
		Surname:        pending.Surname,
		Password:       pending.PasswordHash,
		PassportSeries: pending.PassportSeries,
		PassportNumber: pending.PassportNumber,
		Phone:          pending.Phone,
		Email:          pending.Email,
		Role:           pending.Role,
	}

	if err := s.userStorage.AddNewUserWithAccount(user, pending.BankID); err != nil {
		return err
	}

	return s.pendingStorage.DeleteByEmail(email)
}
