-- name: AddPositionBlacklistedWordToScrape :exec
INSERT INTO scrape_position_blacklisted_words (created_at, updated_at, scrape_id, blacklisted_word)
VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, ?, ?);

