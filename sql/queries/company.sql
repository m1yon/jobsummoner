-- name: GetCompany :one
SELECT * FROM companies
WHERE id = ?;

-- name: AddCompany :exec
INSERT INTO companies (id, created_at, name, url, avatar, source_id)
VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?);