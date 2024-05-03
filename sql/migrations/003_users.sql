-- +goose Up
CREATE TABLE users(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE user_job_postings(
  created_at TIMESTAMP NOT NULL,
  user_id INTEGER NOT NULL,
  position TEXT NOT NULL,
  company_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (position) REFERENCES job_postings(position),
  FOREIGN KEY (company_id) REFERENCES job_postings(company_id)
);

-- +goose Down
DROP TABLE users;
DROP TABLE user_job_postings;

