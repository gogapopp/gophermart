package service

import (
	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type Auth interface {
	CreateUser(user models.User) (int, error)
	GenerateToken(login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Orders interface {
	Create(userID int, order models.Order) (int, error)
	CheckUserOrder(userID int, order models.Order) error
	GetUserOrders(userID int) ([]models.Order, error)
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
		Auth:   NewAuthService(storage.Auth),
		Orders: NewUserOrders(storage.Orders),
	}
}
