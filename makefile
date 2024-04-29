.PHONY:

include .env
export

migrate-up:
	goose -dir "./sql/schema" sqlite $(DB_CONNECTION) up

migrate-down:
	goose -dir "./sql/schema" sqlite $(DB_CONNECTION) down

generate:
	sqlc generate

dev: 
	go run ./cmd/jobsummoner

dev-show: 
	go run ./cmd/jobsummoner -rod=show