package service

import (
	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type UserBalance struct {
	storage storage.Balance
}

func NewUserBalanceService(storage storage.Balance) *UserBalance {
	return &UserBalance{storage: storage}
}

func (s *UserBalance) UpdateUserBalance(userID int, accrual float64) error {
	return s.storage.UpdateUserBalance(userID, accrual)
}

func (s *UserBalance) GetUserBalance(userID int) (models.Balance, error) {
	return s.storage.GetUserBalance(userID)
}
