package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
	"errors"
)

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidEmail = errors.New("invalid email")

type AuthService interface {
	Login(email, password string) (*model.User, error)
	Registrate(user *model.User, bankId int) error
	CreateAccount(bankId, userID int) (*model.UserAccount, error)
}

type authService struct {
	userStorage storage.UserStorage
}

func NewAuthService(userStorage storage.UserStorage) *authService {
	return &authService{
		userStorage: userStorage,
	}
}

func (s *authService) Login(email, password string) (*model.User, error) {
	if email == "" || password == "" { // todo: валидация почты
		return nil, ErrInvalidCredentials
	} else if model, err := s.userStorage.FindByEmail(email); err != nil {
		return nil, err
	} else {
		if model.Password == password { // конечно в реальном мире пароли хранятся в зашифрованном виде, но тема не про информационную безопасность
			return model, nil
		}
		return nil, ErrInvalidCredentials
	}
}

func (s *authService) Registrate(user *model.User, bankId int) error {
	err := s.userStorage.AddNewUserWithAccount(user, bankId)
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
