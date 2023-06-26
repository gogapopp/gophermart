package handler

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

var pgErr *pgconn.PgError
var ErrTooManyRequests = errors.New("too many requests")
var ErrGenerateToken = errors.New("error generate token")
var ErrUnmarshal = errors.New("error unmarshal request body")
var ErrReadBody = errors.New("error read request body")
var ErrDecodingReq = errors.New("error decoding request body")
var ErrUnvalidNumb = errors.New("unvalid order number")
var ErrGetBalance = errors.New("error get user balance")
var ErrEncogingResp = errors.New("error encoding response")
