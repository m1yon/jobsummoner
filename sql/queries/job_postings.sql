-- name: CreateJobPosting :exec
INSERT OR IGNORE INTO job_postings (created_at, updated_at, last_posted, position, url, company_id)
VALUES (?, ?, ?, ?, ?, ?);