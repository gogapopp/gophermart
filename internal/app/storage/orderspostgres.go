package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type UserOrdersPostgres struct {
	db *sql.DB
}

// NewUserOrdersPostgres возвращает *UserOrdersPostgres
func NewUserOrdersPostgres(db *sql.DB) *UserOrdersPostgres {
	return &UserOrdersPostgres{db: db}
}

// Create создаёт ордер в orders и добавляет юзер айди и айди ордера таблицу users_orders
func (p *UserOrdersPostgres) Create(ctx context.Context, userID int, order models.Order) (int, error) {
	var orderID int
	createOrderQuery := fmt.Sprintf("INSERT INTO %s (number, status, accrual, uploaded_at) VALUES ($1, $2, $3, $4) RETURNING id", ordersTable)
	row := p.db.QueryRowContext(ctx, createOrderQuery, order.Number, order.Status, order.Accrual, order.UploadedAt)
	if err := row.Scan(&orderID); err != nil {
		return 0, ErrRepeatValue
	}

	createUsersOrdersQuery := fmt.Sprintf("INSERT INTO %s (user_id, order_id) VALUES ($1, $2)", usersOrdersTable)
	_, err := p.db.ExecContext(ctx, createUsersOrdersQuery, userID, orderID)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

// CheckUserOrder проверяет есть ли у юзера ордер с определённым номером
func (p *UserOrdersPostgres) CheckUserOrder(ctx context.Context, userID int, order models.Order) error {
	checkOrderQuery := fmt.Sprintf("SELECT user_id FROM %s ur INNER JOIN %s o ON ur.order_id = o.id WHERE o.number = $1", usersOrdersTable, ordersTable)
	row := p.db.QueryRowContext(ctx, checkOrderQuery, order.Number)
	var existingUserID int
	err := row.Scan(&existingUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		} else {
			return err
		}
	} else {
		if existingUserID == userID {
			return ErrUserRepeatValue
		} else {
			return ErrRepeatValue
		}
	}
}

// GetUserOrders находит все ордера которые соответствую userID
func (p *UserOrdersPostgres) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	var UserOrders []models.Order
	userOrdersQuery := fmt.Sprintf("SELECT o.number, o.status, o.accrual, o.uploaded_at FROM %s o INNER JOIN %s ur on o.id = ur.order_id WHERE ur.user_id = $1 ORDER BY o.uploaded_at ASC", ordersTable, usersOrdersTable)
	rows, err := p.db.QueryContext(ctx, userOrdersQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		UserOrders = append(UserOrders, order)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return UserOrders, nil
}
