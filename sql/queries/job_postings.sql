-- name: CreateJobPosting :exec
INSERT INTO job_postings (created_at, updated_at, last_posted, position, url, company_id)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?, ?, ?);

-- name: GetJobPostings :many
SELECT job_postings.position, job_postings.url as job_posting_url, companies.name as company_name, last_posted, companies.avatar as company_avatar from job_postings
JOIN companies on job_postings.company_id = companies.id
ORDER BY job_postings.last_posted DESC;

-- name: UpdateJobPostingLastPosted :exec
UPDATE job_postings
SET last_posted = ?, updated_at = CURRENT_TIMESTAMP
WHERE position = ? AND company_id = ?;