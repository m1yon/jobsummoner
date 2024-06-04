-- name: CreateScrape :exec
INSERT INTO scrapes (source_id, created_at)
VALUES (?, CURRENT_TIMESTAMP);