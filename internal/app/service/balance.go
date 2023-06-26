package service

import (
	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type UserBalance struct {
	storage storage.Balance
}

// Не очень понятна этот функционал со "слоями", почему просто интерфейсы не подошли?
// NewUserBalanceService создаёт структуру UserBalance
func NewUserBalanceService(storage storage.Balance) *UserBalance {
	return &UserBalance{storage: storage}
}

// UpdateUserBalanc передаёт данные на слой storage
func (s *UserBalance) UpdateUserBalance(userID int, accrual float64) error {
	return s.storage.UpdateUserBalance(userID, accrual)
}

// GetUserBalance передаёт данные на слой storage
func (s *UserBalance) GetUserBalance(userID int) (models.Balance, error) {
	return s.storage.GetUserBalance(userID)
}
