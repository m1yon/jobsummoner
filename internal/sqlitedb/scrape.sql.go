// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: scrape.sql

package sqlitedb

import (
	"context"
)

const createScrape = `-- name: CreateScrape :exec
INSERT INTO scrapes (source_id, created_at)
VALUES (?, CURRENT_TIMESTAMP)
`

func (q *Queries) CreateScrape(ctx context.Context, sourceID string) error {
	_, err := q.db.ExecContext(ctx, createScrape, sourceID)
	return err
}