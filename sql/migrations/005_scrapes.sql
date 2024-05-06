-- +goose Up
CREATE TABLE scrapes(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  last_scraped TIMESTAMP,
  name TEXT NOT NULL,
  location TEXT NOT NULL,
  work_type INT CHECK (work_type BETWEEN 1 AND 3) NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE scrape_keywords(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  scrape_id INTEGER NOT NULL,
  keyword TEXT NOT NULL,
  FOREIGN KEY (scrape_id) REFERENCES scrapes(id)
);

-- +goose Down
DROP TABLE scrapes;
DROP TABLE scrape_keywords;
