package user

import "errors"

var (
	ErrDuplicateEmail = errors.New("user: duplicate email")
)
