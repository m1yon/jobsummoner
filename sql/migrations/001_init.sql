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
  FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE TABLE companies(
  id TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  avatar TEXT
);

CREATE TABLE users(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP
);

CREATE TABLE user_jobs(
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  user_id INTEGER NOT NULL,
  job_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (job_id) REFERENCES jobs(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, job_id)
);

CREATE TABLE scrapes(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  last_scraped TIMESTAMP,
  name TEXT NOT NULL,
  location TEXT NOT NULL,
  work_type TEXT NOT NULL,
  job_type TEXT NOT NULL,
  salary_range TEXT NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE scrape_keywords(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  keyword TEXT NOT NULL,
  scrape_id INTEGER NOT NULL,
  FOREIGN KEY (scrape_id) REFERENCES scrapes(id) ON DELETE CASCADE
);

CREATE TABLE scrape_work_types(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  work_type TEXT NOT NULL,
  scrape_id INTEGER NOT NULL,
  FOREIGN KEY (scrape_id) REFERENCES scrapes(id) ON DELETE CASCADE
);

CREATE TABLE scrape_job_types(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP,
  job_type TEXT NOT NULL,
  scrape_id INTEGER NOT NULL,
  FOREIGN KEY (scrape_id) REFERENCES scrapes(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE jobs;
DROP TABLE companies;
DROP TABLE users;
DROP TABLE user_jobs;
DROP TABLE scrapes;
DROP TABLE scrape_keywords;
DROP TABLE scrape_work_types;
DROP TABLE scrape_job_types;
