-- name: CreateScrape :exec
INSERT INTO scrapes (created_at, updated_at, last_scraped, name, location, work_type, user_id)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, NULL, ?, ?, ?, ?);

-- name: UpdateLastScraped :exec
UPDATE scrapes
SET last_scraped = CURRENT_TIMESTAMP, SET updated_at = CURRENT_TIMESTAMP
WHERE id = ?;