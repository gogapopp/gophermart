package storage

import (
	"context"
	"database/sql"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type Auth interface {
	CreateUser(ctx context.Context, user models.User) (int, error)
	GetUser(ctx context.Context, login, password string) (models.User, error)
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

type Storage struct {
	Auth
	Orders
	Balance
	Withdrawals
}

// NewStorage возвращает указатель на Storage со встроенными интерфейсами
func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Auth:        NewAuthPostgres(db),
		Orders:      NewUserOrdersPostgres(db),
		Balance:     NewUserBalancePostgres(db),
		Withdrawals: NewWithdrawalsPostgres(db),
	}
}
