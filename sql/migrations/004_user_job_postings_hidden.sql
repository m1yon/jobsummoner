-- +goose Up
ALTER TABLE user_job_postings
ADD COLUMN hidden BOOLEAN NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE user_job_postings DROP COLUMN hidden;