-- name: CreateUser :exec
INSERT INTO users (name, email, hashed_password, created_at)
VALUES (?, ?, ?, CURRENT_TIMESTAMP);

-- name: GetUserCredentials :one
SELECT id, hashed_password
FROM users
WHERE email = ?;