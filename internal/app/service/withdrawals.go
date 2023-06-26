package service

import (
	"context"

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
func (p *WithdrawalsService) UserWithdraw(ctx context.Context, userID int, withdraw models.Withdraw) error {
	return p.storage.UserWithdraw(ctx, userID, withdraw)
}

// UserWithdraw передаёт данные на слой storage
func (p *WithdrawalsService) GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error) {
	return p.storage.GetUserWithdrawals(ctx, userID)
}
