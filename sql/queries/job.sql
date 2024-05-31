-- name: AddJob :exec
INSERT INTO jobs (id, created_at, position, url, company_id, location, source_id)
VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?,  ?);