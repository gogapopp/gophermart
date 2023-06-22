package service

import (
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
func (s *UserOrders) Create(userID int, order models.Order) (int, error) {
	return s.storage.Create(userID, order)
}

// CheckUserOrder передаёт данные на слой storage
func (s *UserOrders) CheckUserOrder(userID int, order models.Order) error {
	return s.storage.CheckUserOrder(userID, order)
}

// GetUserOrders передаёт данные на слой storage
func (s *UserOrders) GetUserOrders(userID int) ([]models.Order, error) {
	return s.storage.GetUserOrders(userID)
}
