package service

import (
	"context"

	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type Auth interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GenerateToken(ctx context.Context, login, password string) (string, error)
	ParseToken(token string) (int, error)
}

type Orders interface {
	Create(ctx context.Context, userID int, order models.Order) (int, error)
	CheckUserOrder(ctx context.Context, userID int, order models.Order) error
	GetUserOrders(ctx context.Context, userID int) ([]models.Order, error)
}

type Balance interface {
	UpdateUserBalance(ctx context.Context, userID int, accrual float64) error
	GetUserBalance(ctx context.Context, userID int) (models.Balance, error)
}

type Withdrawals interface {
	UserWithdraw(ctx context.Context, userID int, withdraw models.Withdraw) error
	GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error)
}

type Service struct {
	Auth
	Orders
	Balance
	Withdrawals
}

// NewService возвращает указатель на Service со встроенными интерфейсами
func NewService(storage *storage.Storage) *Service {
	return &Service{
		Auth:        NewAuthService(storage.Auth),
		Orders:      NewUserOrdersService(storage.Orders),
		Balance:     NewUserBalanceService(storage.Balance),
		Withdrawals: NewWithdrawalsService(storage.Withdrawals),
	}
}
