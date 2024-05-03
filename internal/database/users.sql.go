// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: users.sql

package database

import (
	"context"
)

const createUser = `-- name: CreateUser :exec
INSERT INTO users (created_at, updated_at)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
`

func (q *Queries) CreateUser(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, createUser)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, updated_at FROM users WHERE id = ?
`

func (q *Queries) GetUser(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(&i.ID, &i.CreatedAt, &i.UpdatedAt)
	return i, err
}
