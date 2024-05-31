-- name: AddJob :exec
INSERT INTO jobs (id, created_at, last_posted, position, url, company_id, location, source_id)
VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?, ?,  ?);

-- name: GetJob :one
SELECT * FROM jobs
WHERE id = ?;

