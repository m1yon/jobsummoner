-- name: GetCompany :one
SELECT * FROM companies
WHERE id = ?;

-- name: CreateCompany :exec
INSERT INTO companies (id, created_at, name, url, avatar, source_id)
VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?);