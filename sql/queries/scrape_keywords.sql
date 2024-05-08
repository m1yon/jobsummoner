-- name: GetAllScrapesWithKeywords :many
SELECT *, GROUP_CONCAT(scrape_keywords.keyword, " OR ") AS keywords
FROM scrapes
JOIN scrape_keywords on scrapes.id = scrape_keywords.scrape_id
GROUP BY scrapes.id;

-- name: AddKeywordToScrape :exec
INSERT INTO scrape_keywords (created_at, updated_at, scrape_id, keyword)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?);

