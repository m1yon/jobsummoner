-- +goose Up
CREATE TABLE job_postings(
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  last_posted TIMESTAMP NOT NULL,
  role TEXT NOT NULL,
  url TEXT NOT NULL,
  company_id TEXT NOT NULL,
  FOREIGN KEY (company_id) REFERENCES companies(id),
  PRIMARY KEY (role, company_id)
);

CREATE TABLE companies(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL
);


-- +goose Down
DROP TABLE job_postings;
DROP TABLE companies;

