package handler

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var pgErr *pgconn.PgError
var ErrTooManyRequests = errors.New("TooManyRequests")
