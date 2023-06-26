package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type UserBalancePostgres struct {
	db *sql.DB
}

// NewUserBalancePostgres возвращает *UserBalancePostgres
func NewUserBalancePostgres(db *sql.DB) *UserBalancePostgres {
	return &UserBalancePostgres{db: db}
}

// UpdateUserBalance обновляет баланс юзера
func (p *UserBalancePostgres) UpdateUserBalance(ctx context.Context, userID int, accrual float64) error {
	var currentBalanceStr string
	var currentBalance, accrualBig big.Float
	userCurrentBalanceQuery := fmt.Sprintf("SELECT current_balance FROM %s WHERE user_id = $1", usersBalance)
	row := p.db.QueryRowContext(ctx, userCurrentBalanceQuery, userID)
	if err := row.Scan(&currentBalanceStr); err != nil {
		return err
	}
	accrualBig.SetFloat64(accrual)
	currentBalance.SetString(currentBalanceStr)
	currentBalance.Add(&currentBalance, &accrualBig)

	updateCurrentBalanceQuery := fmt.Sprintf("UPDATE %s SET current_balance = $1 WHERE user_id = $2", usersBalance)
	_, err := p.db.ExecContext(ctx, updateCurrentBalanceQuery, &currentBalance, userID)
	if err != nil {
		return err
	}

	return nil
}

// GetUserBalance получает баланс юзера
func (p *UserBalancePostgres) GetUserBalance(ctx context.Context, userID int) (models.Balance, error) {
	var userBalance models.Balance
	getUserBalanceQuery := fmt.Sprintf("SELECT current_balance, withdrawn FROM %s WHERE user_id = $1", usersBalance)
	row := p.db.QueryRowContext(ctx, getUserBalanceQuery, userID)
	err := row.Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		return userBalance, err
	}

	return userBalance, err
}
