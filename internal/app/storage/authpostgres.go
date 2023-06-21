package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type AuthPostgres struct {
	ctx context.Context
	db  *sql.DB
}

func NewAuthPostgres(ctx context.Context, db *sql.DB) *AuthPostgres {
	return &AuthPostgres{ctx: ctx, db: db}
}

func (p *AuthPostgres) CreateUser(user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, password) values ($1, $2) RETURNING id", usersTable)
	row := p.db.QueryRowContext(p.ctx, query, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		fmt.Println(err, "6")
		return 0, err
	}

	createUserBalance := fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", usersBalance)
	_, err := p.db.ExecContext(p.ctx, createUserBalance, id)
	if err != nil {
		fmt.Println(err, "7")
		return 0, err
	}

	return id, nil
}

func (p *AuthPostgres) GetUser(login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password=$2", usersTable)
	row := p.db.QueryRowContext(p.ctx, query, login, password)
	if err := row.Scan(&user.ID); err != nil {
		fmt.Println(err, "8")
		return user, ErrNoRows
	}
	return user, nil
}
