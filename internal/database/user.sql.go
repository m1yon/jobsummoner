// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (name, email, hashed_password, created_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP)
`

type CreateUserParams struct {
	Name           string
	Email          string
	HashedPassword string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) error {
	_, err := q.db.ExecContext(ctx, createUser, arg.Name, arg.Email, arg.HashedPassword)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, name, email, hashed_password, created_at, updated_at
FROM users
WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Email,
		&i.HashedPassword,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserCredentials = `-- name: GetUserCredentials :one
SELECT id, hashed_password
FROM users
WHERE email = ?
`

type GetUserCredentialsRow struct {
	ID             int64
	HashedPassword string
}

func (q *Queries) GetUserCredentials(ctx context.Context, email string) (GetUserCredentialsRow, error) {
	row := q.db.QueryRowContext(ctx, getUserCredentials, email)
	var i GetUserCredentialsRow
	err := row.Scan(&i.ID, &i.HashedPassword)
	return i, err
}
