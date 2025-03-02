package service

import (
	"banksystem/internal/model"
	"banksystem/internal/storage"
)

type AuthService interface {
	Login(username string) (*model.User, error)
}

type authService struct {
	userStorage storage.UserStorage
}

func NewAuthService(userStorage storage.UserStorage) AuthService {
	return &authService{
		userStorage: userStorage,
	}
}

func (s *authService) Login(username string) (*model.User, error) {
	return s.userStorage.FindByUsername(username)
}
