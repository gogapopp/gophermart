package storage

import "errors"

var ErrNoRows = errors.New("sql: no rows in result set")
var ErrRepeatValue = errors.New("ErrRepeatValue")
var ErrUserRepeatValue = errors.New("ErrUserRepeatValue")
