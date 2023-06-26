package handler

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

// Вот тут было бы хорошо описать все дублирующиеся ошибки
var pgErr *pgconn.PgError
var ErrTooManyRequests = errors.New("TooManyRequests")
