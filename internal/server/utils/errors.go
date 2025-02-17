package utils

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrPasswordCorrect = errors.New("invalid password")
)
