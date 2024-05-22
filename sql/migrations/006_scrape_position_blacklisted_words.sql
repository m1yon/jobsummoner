-- +goose Up
CREATE TABLE scrape_position_blacklisted_words(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  scrape_id INTEGER NOT NULL,
  blacklisted_word TEXT NOT NULL,
  FOREIGN KEY (scrape_id) REFERENCES scrapes(id)
);

-- +goose Down
DROP TABLE scrape_position_blacklisted_words;
