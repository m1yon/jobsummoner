-- name: CreateScrape :exec
INSERT INTO scrapes (created_at, updated_at, last_scraped, name, location, work_type, user_id)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, ?, ?, ?, ?);

-- name: UpdateLastScraped :exec
UPDATE scrapes
SET last_scraped = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: GetLastScrapedDate :one
SELECT last_scraped
FROM scrapes
WHERE id = ?
ORDER BY last_scraped DESC
LIMIT 1;