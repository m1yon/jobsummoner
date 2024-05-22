-- name: AddKeywordToScrape :exec
INSERT INTO scrape_keywords (created_at, updated_at, scrape_id, keyword)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?);

