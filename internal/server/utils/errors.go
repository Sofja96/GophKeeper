package utils

import "errors"

var (
	ErrUserExists       = errors.New("user already exists")
	ErrUserDataNotFound = errors.New("no data found")
)
