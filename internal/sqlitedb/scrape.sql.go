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

const getLastScrape = `-- name: GetLastScrape :one
SELECT id, source_id, created_at FROM scrapes
WHERE source_id = ?
ORDER BY created_at DESC
LIMIT 1
`

func (q *Queries) GetLastScrape(ctx context.Context, sourceID string) (Scrape, error) {
	row := q.db.QueryRowContext(ctx, getLastScrape, sourceID)
	var i Scrape
	err := row.Scan(&i.ID, &i.SourceID, &i.CreatedAt)
	return i, err
}
