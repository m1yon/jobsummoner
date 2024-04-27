-- name: CreateCompany :exec
INSERT INTO companies (id, created_at, updated_at, name, url)
VALUES (?, ?, ?, ?, ?);

-- name: GetCompany :one
SELECT * from job_postings;