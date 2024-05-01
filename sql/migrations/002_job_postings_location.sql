-- +goose Up
ALTER TABLE job_postings
ADD COLUMN location_type INTEGER CHECK ( location_type BETWEEN 1 AND 3 ) NOT NULL DEFAULT 2;

ALTER TABLE job_postings
ADD COLUMN location TEXT;

-- +goose Down
ALTER TABLE job_postings DROP COLUMN location_type;
ALTER TABLE job_postings DROP COLUMN location;
