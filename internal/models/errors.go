package models

import "errors"

var (
	ErrDuplicateEmail = errors.New("user: duplicate email")
)
