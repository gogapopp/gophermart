package service

import (
	"context"

	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type UserBalance struct {
	storage storage.Balance
}

// NewUserBalanceService создаёт структуру UserBalance
func NewUserBalanceService(storage storage.Balance) *UserBalance {
	return &UserBalance{storage: storage}
}

// UpdateUserBalanc передаёт данные на слой storage
func (s *UserBalance) UpdateUserBalance(ctx context.Context, userID int, accrual float64) error {
	return s.storage.UpdateUserBalance(ctx, userID, accrual)
}

// GetUserBalance передаёт данные на слой storage
func (s *UserBalance) GetUserBalance(ctx context.Context, userID int) (models.Balance, error) {
	return s.storage.GetUserBalance(ctx, userID)
}
