-- name: CreateCompany :exec
INSERT OR IGNORE INTO companies (id, created_at, updated_at, name, url, avatar)
VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?);