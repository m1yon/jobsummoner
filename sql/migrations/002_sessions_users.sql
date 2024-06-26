-- +goose Up
CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BLOB NOT NULL,
	expiry REAL NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE users (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	hashed_password TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP
);

-- +goose Down
DROP TABLE sessions;
DROP TABLE users;

