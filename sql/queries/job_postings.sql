-- name: CreateJobPosting :exec
INSERT OR IGNORE INTO job_postings (created_at, updated_at, last_posted, position, url, company_id)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetJobPostings :many
SELECT job_postings.position, job_postings.url as job_posting_url, companies.name as company_name, last_posted from job_postings
JOIN companies on job_postings.company_id = companies.id
ORDER BY job_postings.last_posted DESC;