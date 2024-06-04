-- name: CreateScrape :exec
INSERT INTO scrapes (source_id, created_at)
VALUES (?, CURRENT_TIMESTAMP);

-- name: GetLastScrape :one
SELECT * FROM scrapes
WHERE source_id = ?
ORDER BY created_at DESC
LIMIT 1;