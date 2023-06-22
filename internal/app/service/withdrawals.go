package service

import (
	"github.com/gogapopp/gophermart/internal/app/models"
	"github.com/gogapopp/gophermart/internal/app/storage"
)

type WithdrawalsService struct {
	storage storage.Withdrawals
}

func NewWithdrawalsService(storage storage.Withdrawals) *WithdrawalsService {
	return &WithdrawalsService{storage: storage}
}

// UserWithdraw передаёт данные на слой storage
func (p *WithdrawalsService) UserWithdraw(userID int, withdraw models.Withdraw) error {
	return p.storage.UserWithdraw(userID, withdraw)
}

// UserWithdraw передаёт данные на слой storage
func (p *WithdrawalsService) GetUserWithdrawals(userID int) ([]models.Withdraw, error) {
	return p.storage.GetUserWithdrawals(userID)
}
