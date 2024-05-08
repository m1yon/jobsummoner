-- +goose Up
ALTER TABLE user_job_postings
ADD COLUMN status INTEGER NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE user_job_postings DROP COLUMN status;