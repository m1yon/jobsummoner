// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: job.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createJob = `-- name: CreateJob :exec
INSERT INTO jobs (id, created_at, last_posted, position, url, company_id, location, source_id, last_posted)
VALUES (?, CURRENT_TIMESTAMP, ?, ?, ?, ?, ?, ?, ?)
`

type CreateJobParams struct {
	ID           string
	LastPosted   time.Time
	Position     string
	Url          string
	CompanyID    string
	Location     sql.NullString
	SourceID     string
	LastPosted_2 time.Time
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) error {
	_, err := q.db.ExecContext(ctx, createJob,
		arg.ID,
		arg.LastPosted,
		arg.Position,
		arg.Url,
		arg.CompanyID,
		arg.Location,
		arg.SourceID,
		arg.LastPosted_2,
	)
	return err
}

const getJob = `-- name: GetJob :one
SELECT jobs.position, jobs.location, jobs.url AS job_url, companies.url AS company_url, companies.name AS company_name, companies.id AS company_id, jobs.source_id, companies.avatar AS company_avatar, companies.url AS company_url, jobs.last_posted
FROM jobs
JOIN companies ON jobs.company_id = companies.id
WHERE jobs.id = ?
`

type GetJobRow struct {
	Position      string
	Location      sql.NullString
	JobUrl        string
	CompanyUrl    string
	CompanyName   string
	CompanyID     string
	SourceID      string
	CompanyAvatar sql.NullString
	CompanyUrl_2  string
	LastPosted    time.Time
}

func (q *Queries) GetJob(ctx context.Context, id string) (GetJobRow, error) {
	row := q.db.QueryRowContext(ctx, getJob, id)
	var i GetJobRow
	err := row.Scan(
		&i.Position,
		&i.Location,
		&i.JobUrl,
		&i.CompanyUrl,
		&i.CompanyName,
		&i.CompanyID,
		&i.SourceID,
		&i.CompanyAvatar,
		&i.CompanyUrl_2,
		&i.LastPosted,
	)
	return i, err
}

const getJobs = `-- name: GetJobs :many
SELECT jobs.position, jobs.location, jobs.url AS job_url, companies.url AS company_url, companies.name AS company_name, companies.id AS company_id, jobs.source_id, companies.avatar AS company_avatar, companies.url AS company_url, jobs.last_posted
FROM jobs
JOIN companies ON jobs.company_id = companies.id
`

type GetJobsRow struct {
	Position      string
	Location      sql.NullString
	JobUrl        string
	CompanyUrl    string
	CompanyName   string
	CompanyID     string
	SourceID      string
	CompanyAvatar sql.NullString
	CompanyUrl_2  string
	LastPosted    time.Time
}

func (q *Queries) GetJobs(ctx context.Context) ([]GetJobsRow, error) {
	rows, err := q.db.QueryContext(ctx, getJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetJobsRow
	for rows.Next() {
		var i GetJobsRow
		if err := rows.Scan(
			&i.Position,
			&i.Location,
			&i.JobUrl,
			&i.CompanyUrl,
			&i.CompanyName,
			&i.CompanyID,
			&i.SourceID,
			&i.CompanyAvatar,
			&i.CompanyUrl_2,
			&i.LastPosted,
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
