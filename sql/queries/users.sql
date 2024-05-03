-- name: CreateUser :exec
INSERT INTO users (created_at, updated_at)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;