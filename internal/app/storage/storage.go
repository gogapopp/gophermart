package storage

import (
	"context"
	"database/sql"

	"github.com/gogapopp/gophermart/internal/app/models"
)

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
}

type Storage struct {
	Auth
	Orders
	Balance
}

func NewStorage(ctx context.Context, db *sql.DB) *Storage {
	return &Storage{
		Auth:   NewAuthPostgres(ctx, db),
		Orders: NewUserOrdersPostgres(ctx, db),
	}
}
