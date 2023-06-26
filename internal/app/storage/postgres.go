package storage

import (
	"context"
	"database/sql"

	"github.com/gogapopp/gophermart/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	usersTable       = "users"
	ordersTable      = "orders"
	usersOrdersTable = "users_orders"
	usersBalance     = "user_balance"
	usersWithdrawals = "withdrawals"
)

// NewDB создаёт таблицы и индексы
func NewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.DatabaseURI)
	if err != nil {
		return nil, err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id serial not null unique,
			login varchar(256),
			password varchar(256)
		)
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS orders (
			id serial not null unique,
			number varchar(256),
			status varchar(256),
			accrual decimal,
			uploaded_at timestamptz
		)
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS user_balance (
			id serial not null unique,
			user_id int,
			current_balance decimal default 0,
			withdrawn decimal default 0
		)
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS withdrawals (
			id serial not null unique,
			user_id int,
			order_id varchar(256),
			sum decimal,
			processed_at timestamptz
		)
	`)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users_orders (
			id serial not null unique,
			user_id int references users (id) on delete cascade,
			order_id int references orders (id) on delete cascade
		)
	`)
	// Rollback обычно вызывают тоже через defer
	// Вообще можно посмотреть как оптимизировать транзакцию https://go.dev/doc/database/execute-transactions
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS login_id ON users(login)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.ExecContext(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS order_id ON orders(number)")
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return db, tx.Commit()
}
