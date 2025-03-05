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
}

type authService struct {
	userStorage storage.UserStorage
}

func NewAuthService(userStorage storage.UserStorage) AuthService {
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
