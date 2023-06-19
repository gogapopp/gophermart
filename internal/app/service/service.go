package service

import (
	"github.com/gogapopp/gophermart/internal/app/storage"
	"github.com/gogapopp/gophermart/models"
)

type Auth interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Orders interface {
}

type Balance interface {
}

type Service struct {
	Auth
	Orders
	Balance
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Auth: NewAuthService(storage.Auth),
	}
}
