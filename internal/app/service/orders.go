package service

import (
	"context"

	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type UserOrders struct {
	storage storage.Orders
}

// NewUserOrdersService создаём структуру UserOrders
func NewUserOrdersService(storage storage.Orders) *UserOrders {
	return &UserOrders{storage: storage}
}

// Create передаёт данные на слой storage
func (s *UserOrders) Create(ctx context.Context, userID int, order models.Order) (int, error) {
	return s.storage.Create(ctx, userID, order)
}

// CheckUserOrder передаёт данные на слой storage
func (s *UserOrders) CheckUserOrder(ctx context.Context, userID int, order models.Order) error {
	return s.storage.CheckUserOrder(ctx, userID, order)
}

// GetUserOrders передаёт данные на слой storage
func (s *UserOrders) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	return s.storage.GetUserOrders(ctx, userID)
}
