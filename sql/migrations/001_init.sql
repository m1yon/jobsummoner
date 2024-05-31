-- +goose Up
CREATE TABLE jobs(
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  last_posted TIMESTAMP NOT NULL,
  position TEXT NOT NULL,
  url TEXT NOT NULL,
  company_id TEXT NOT NULL,
  location TEXT,
  source_id TEXT NOT NULL,
  FOREIGN KEY (source_id) REFERENCES scraping_sources(id) ON DELETE CASCADE,
  FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE TABLE companies(
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  avatar TEXT,
  source_id TEXT NOT NULL,
  FOREIGN KEY (source_id) REFERENCES scraping_sources(id) ON DELETE CASCADE
);

CREATE TABLE scraping_sources (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL,
  config JSON
);

-- +goose Down
DROP TABLE jobs;
DROP TABLE companies;
DROP TABLE scraping_sources;

