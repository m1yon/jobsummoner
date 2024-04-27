.PHONY:

include .env
export

dev:
	go run ./cmd/jobsummoner

migrate-up:
	goose -dir "./sql/schema" sqlite $(DB_CONNECTION) up

migrate-down:
	goose -dir "./sql/schema" sqlite $(DB_CONNECTION) down

generate:
	sqlc generate