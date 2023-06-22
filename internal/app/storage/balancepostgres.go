package storage

import (
	"context"
	"database/sql"
	"fmt"
	"math/big"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type UserBalancePostgres struct {
	ctx context.Context
	db  *sql.DB
}

func NewUserBalancePostgres(ctx context.Context, db *sql.DB) *UserBalancePostgres {
	return &UserBalancePostgres{ctx: ctx, db: db}
}

func (p *UserBalancePostgres) UpdateUserBalance(userID int, accrual float64) error {
	var currentBalanceStr string
	var currentBalance, accrualBig big.Float
	userCurrentBalanceQuery := fmt.Sprintf("SELECT current_balance FROM %s WHERE user_id = $1", usersBalance)
	row := p.db.QueryRowContext(p.ctx, userCurrentBalanceQuery, userID)
	if err := row.Scan(&currentBalance); err != nil {
		fmt.Println(err)
		return err
	}
	accrualBig.SetFloat64(accrual)
	currentBalance.SetString(currentBalanceStr)
	currentBalance.Add(&currentBalance, &accrualBig)

	updateCurrentBalanceQuery := fmt.Sprintf("UPDATE %s SET current_balance = $1 WHERE user_id = $2", usersBalance)
	_, err := p.db.ExecContext(p.ctx, updateCurrentBalanceQuery, &currentBalance, userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (p *UserBalancePostgres) GetUserBalance(userID int) (models.Balance, error) {
	var userBalance models.Balance
	getUserBalanceQuery := fmt.Sprintf("SELECT current_balance, withdrawn FROM %s WHERE user_id = $1", usersBalance)
	row := p.db.QueryRowContext(p.ctx, getUserBalanceQuery, userID)
	err := row.Scan(&userBalance.Current, &userBalance.Withdrawn)
	if err != nil {
		fmt.Println(err)
		return userBalance, err
	}

	return userBalance, err
}
