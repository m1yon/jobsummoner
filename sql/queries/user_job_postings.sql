-- name: CreateUserJobPosting :exec
INSERT INTO user_job_postings (created_at, user_id, job_posting_id)
VALUES (CURRENT_TIMESTAMP, ?, ?);

-- name: GetUserJobPostingsByStatus :many
SELECT job_postings.position, job_postings.url as job_posting_url, companies.name as company_name, last_posted, companies.avatar as company_avatar, companies.id as company_id, job_postings.id as job_posting_id, user_job_postings.status
FROM user_job_postings
JOIN companies on job_postings.company_id = companies.id
JOIN job_postings on user_job_postings.job_posting_id = job_postings.id
WHERE user_job_postings.user_id = ? AND user_job_postings.status = ?
ORDER BY job_postings.last_posted DESC;

-- name: UpdateUserJobPostingStatus :exec
UPDATE user_job_postings
SET status = ?
WHERE user_id = ? AND job_posting_id = ?;