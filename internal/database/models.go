// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"time"
)

type Company struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
}

type JobPosting struct {
	CreatedAt  time.Time
	UpdatedAt  time.Time
	LastPosted time.Time
	Position   string
	Url        string
	CompanyID  string
}
