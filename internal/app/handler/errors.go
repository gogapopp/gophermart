package handler

import "errors"

var ErrTooManyRequests = errors.New("TooManyRequests")
var ErrNoContent = errors.New("NoContent")
