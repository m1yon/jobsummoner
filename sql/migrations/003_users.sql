-- +goose Up
CREATE TABLE users(
  id INTEGER PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE TABLE user_job_postings(
  created_at TIMESTAMP NOT NULL,
  user_id INTEGER NOT NULL,
  job_posting_id TEXT NOT NULL,
  FOREIGN KEY (user_id) REFERENCES users(id),
  FOREIGN KEY (job_posting_id) REFERENCES job_postings(id),
  PRIMARY KEY (user_id, job_posting_id)
);

-- +goose Down
DROP TABLE users;
DROP TABLE user_job_postings;

