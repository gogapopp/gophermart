package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type UserOrdersPostgres struct {
	ctx context.Context
	db  *sql.DB
}

func NewUserOrdersPostgres(ctx context.Context, db *sql.DB) *UserOrdersPostgres {
	return &UserOrdersPostgres{ctx: ctx, db: db}
}

func (p *UserOrdersPostgres) Create(userID int, order models.Order) (int, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return 0, err
	}

	var orderID int
	createOrderQuery := fmt.Sprintf("INSERT INTO %s (number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4) RETURNING id", ordersTable)
	row := tx.QueryRowContext(p.ctx, createOrderQuery, order.Number, order.Status, order.Accrual, order.UploadedAt)
	if err := row.Scan(orderID); err != nil {
		if err := row.Scan(&orderID); err != nil {
			return 0, err
		}
		tx.Rollback()
		return 0, err
	}

	createUsersOrdersQuery := fmt.Sprintf("INSERT INTO %s (user_id, order_id) VALUES ($1, $2)", usersOrdersTable)
	_, err = tx.ExecContext(p.ctx, createUsersOrdersQuery, userID, orderID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return orderID, tx.Commit()
}

func (p *UserOrdersPostgres) GetUserOrders(userID int) ([]models.Order, error) {
	var UserOrders []models.Order
	userOrdersQuery := fmt.Sprintf("SELECT * FROM %s o INNER JOIN %s ur on o.id = ur.order_id WHERE ur.user_id = $1", ordersTable, usersOrdersTable)
	rows, err := p.db.QueryContext(p.ctx, userOrdersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		UserOrders = append(UserOrders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return UserOrders, nil
}
