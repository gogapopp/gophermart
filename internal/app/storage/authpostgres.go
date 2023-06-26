package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gogapopp/gophermart/internal/app/models"
)

type AuthPostgres struct {
	db *sql.DB
}

// NewAuthPostgres возвращает *AuthPostgres
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateUser создаёт юзера и задаёт ему нулевой баланс
func (p *AuthPostgres) CreateUser(ctx context.Context, user models.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (login, password) values ($1, $2) RETURNING id", usersTable)
	row := p.db.QueryRowContext(ctx, query, user.Login, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}

	createUserBalance := fmt.Sprintf("INSERT INTO %s (user_id) values ($1)", usersBalance)
	_, err := p.db.ExecContext(ctx, createUserBalance, id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetUser ищет по паролю и логину юзера и возвращает его
func (p *AuthPostgres) GetUser(ctx context.Context, login, password string) (models.User, error) {
	var user models.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE login=$1 AND password=$2", usersTable)
	row := p.db.QueryRowContext(ctx, query, login, password)
	if err := row.Scan(&user.ID); err != nil {
		return user, ErrNoRows
	}
	return user, nil
}
