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

-- name: GetAllScrapes :many
SELECT *, scrapes.id, rtrim(replace(group_concat(DISTINCT LOWER(scrape_keywords.keyword)||'@!'), '@!,', ' OR '),'@!') AS keywords, rtrim(replace(group_concat(DISTINCT scrape_position_blacklisted_words.blacklisted_word||'@!'), '@!,', ','),'@!') AS blacklisted_words
FROM scrapes
JOIN scrape_keywords on scrapes.id = scrape_keywords.scrape_id
JOIN scrape_position_blacklisted_words on scrapes.id = scrape_position_blacklisted_words.scrape_id
GROUP BY scrapes.id;