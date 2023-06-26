package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type WithdrawalsPostgres struct {
	db *sql.DB
}

// NewWithdrawalsPostgres возвращает *WithdrawalsPostgres
func NewWithdrawalsPostgres(db *sql.DB) *WithdrawalsPostgres {
	return &WithdrawalsPostgres{db: db}
}

// UserWithdraw записывает withdraw юзера в БД
func (p *WithdrawalsPostgres) UserWithdraw(ctx context.Context, userID int, withdraw models.Withdraw) error {
	userWithdrawQuery := fmt.Sprintf("INSERT INTO %s (user_id, order_id, sum, processed_at) VALUES ($1, $2, $3, $4)", usersWithdrawals)
	_, err := p.db.ExecContext(ctx, userWithdrawQuery, userID, withdraw.Order, withdraw.Sum, withdraw.ProcessedAt)
	if err != nil {
		return err
	}

	updateWithdrawBalanceQuery := fmt.Sprintf("UPDATE %s SET withdrawn = $1 WHERE user_id = $2", usersBalance)
	_, err = p.db.ExecContext(ctx, updateWithdrawBalanceQuery, &withdraw.Sum, userID)
	if err != nil {
		return err
	}

	return err
}

// GetUserWithdrawals получает все withdraw юзера из БД
func (p *WithdrawalsPostgres) GetUserWithdrawals(ctx context.Context, userID int) ([]models.Withdraw, error) {
	var userWithdrawals []models.Withdraw
	userWithdrawalsQuery := fmt.Sprintf("SELECT order_id, sum, processed_at FROM %s WHERE user_id=$1 ORDER BY processed_at ASC", usersWithdrawals)
	rows, err := p.db.QueryContext(ctx, userWithdrawalsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var withdraw models.Withdraw
		if err := rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt); err != nil {
			return nil, err
		}
		userWithdrawals = append(userWithdrawals, withdraw)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userWithdrawals, nil
}
