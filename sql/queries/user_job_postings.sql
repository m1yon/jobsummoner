-- name: CreateUserJobPosting :exec
INSERT INTO user_job_postings (created_at, user_id, position, company_id)
VALUES (CURRENT_TIMESTAMP, ?, ?, ?);

-- name: GetUserJobPostings :many
SELECT job_postings.position, job_postings.url as job_posting_url, companies.name as company_name, last_posted, companies.avatar as company_avatar 
FROM user_job_postings
JOIN companies on job_postings.company_id = companies.id
JOIN job_postings on user_job_postings.position = job_postings.position AND user_job_postings.company_id = job_postings.company_id
WHERE user_job_postings.user_id = ?
ORDER BY job_postings.last_posted DESC;