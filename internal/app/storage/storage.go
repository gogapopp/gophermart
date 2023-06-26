package storage

import (
	"context"
	"database/sql"

	"github.com/gogapopp/gophermart/internal/app/models"
)

// Использования одной структуры storage позволит и 1 интерфейс сдееать на все
type Auth interface {
	CreateUser(user models.User) (int, error)
	GetUser(login, password string) (models.User, error)
}

type Orders interface {
	Create(userID int, order models.Order) (int, error)
	CheckUserOrder(userID int, order models.Order) error
	GetUserOrders(userID int) ([]models.Order, error)
}

type Balance interface {
	UpdateUserBalance(userID int, accrual float64) error
	GetUserBalance(userID int) (models.Balance, error)
}

type Withdrawals interface {
	UserWithdraw(userID int, withdraw models.Withdraw) error
	GetUserWithdrawals(userID int) ([]models.Withdraw, error)
}

type Storage struct {
	Auth
	Orders
	Balance
	Withdrawals
}

// NewStorage возвращает указатель на Storage со встроенными интерфейсами
func NewStorage(ctx context.Context, db *sql.DB) *Storage {
	return &Storage{
		Auth:        NewAuthPostgres(ctx, db),
		Orders:      NewUserOrdersPostgres(ctx, db),
		Balance:     NewUserBalancePostgres(ctx, db),
		Withdrawals: NewWithdrawalsPostgres(ctx, db),
	}
}
