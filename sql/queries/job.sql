-- name: CreateJob :exec
INSERT INTO jobs (id, created_at, last_posted, position, url, company_id, location, source_id, last_posted)
VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?, ?, ?);

-- name: GetJob :one
SELECT jobs.position, jobs.location, jobs.url AS job_url, companies.url AS company_url, companies.name AS company_name, companies.id AS company_id, jobs.source_id, companies.avatar AS company_avatar, companies.url AS company_url, jobs.last_posted
FROM jobs
JOIN companies ON jobs.company_id = companies.id
WHERE jobs.id = ?;

-- name: GetJobs :many
SELECT jobs.position, jobs.location, jobs.url AS job_url, companies.url AS company_url, companies.name AS company_name, companies.id AS company_id, jobs.source_id, companies.avatar AS company_avatar, companies.url AS company_url, jobs.last_posted
FROM jobs
JOIN companies ON jobs.company_id = companies.id;
