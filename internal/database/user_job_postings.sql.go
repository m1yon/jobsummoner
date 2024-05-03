// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user_job_postings.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createUserJobPosting = `-- name: CreateUserJobPosting :exec
INSERT INTO user_job_postings (created_at, user_id, job_posting_id)
VALUES (CURRENT_TIMESTAMP, ?, ?)
`

type CreateUserJobPostingParams struct {
	UserID       int64
	JobPostingID string
}

func (q *Queries) CreateUserJobPosting(ctx context.Context, arg CreateUserJobPostingParams) error {
	_, err := q.db.ExecContext(ctx, createUserJobPosting, arg.UserID, arg.JobPostingID)
	return err
}

const getUserJobPostings = `-- name: GetUserJobPostings :many
SELECT job_postings.position, job_postings.url as job_posting_url, companies.name as company_name, last_posted, companies.avatar as company_avatar, companies.id as company_id, job_postings.id as job_posting_id
FROM user_job_postings
JOIN companies on job_postings.company_id = companies.id
JOIN job_postings on user_job_postings.job_posting_id = job_postings.id
WHERE user_job_postings.user_id = ? AND user_job_postings.hidden = false
ORDER BY job_postings.last_posted DESC
`

type GetUserJobPostingsRow struct {
	Position      string
	JobPostingUrl string
	CompanyName   string
	LastPosted    time.Time
	CompanyAvatar sql.NullString
	CompanyID     string
	JobPostingID  string
}

func (q *Queries) GetUserJobPostings(ctx context.Context, userID int64) ([]GetUserJobPostingsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserJobPostings, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserJobPostingsRow
	for rows.Next() {
		var i GetUserJobPostingsRow
		if err := rows.Scan(
			&i.Position,
			&i.JobPostingUrl,
			&i.CompanyName,
			&i.LastPosted,
			&i.CompanyAvatar,
			&i.CompanyID,
			&i.JobPostingID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const hideUserJobPosting = `-- name: HideUserJobPosting :exec
UPDATE user_job_postings
SET hidden = 1
WHERE user_id = ? AND job_posting_id = ?
`

type HideUserJobPostingParams struct {
	UserID       int64
	JobPostingID string
}

func (q *Queries) HideUserJobPosting(ctx context.Context, arg HideUserJobPostingParams) error {
	_, err := q.db.ExecContext(ctx, hideUserJobPosting, arg.UserID, arg.JobPostingID)
	return err
}
