package models

import "errors"

var (
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrNoRecord           = errors.New("models: no matching record found")
)

const (
	DBErrUniqueConstraint = "UNIQUE constraint failed"
	DBErrNoResults        = "no rows in result set"
)
