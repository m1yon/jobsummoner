.PHONY:

include .env
export

migrate-up:
	goose -dir "./sql/migrations" sqlite ./db/database.db up

migrate-down:
	goose -dir "./sql/migrations" sqlite ./db/database.db down

generate:
	sqlc generate

dev: 
	go run ./cmd/jobsummoner

dev-show: 
	go run ./cmd/jobsummoner -rod=show

reset-db:
	rm -f ./database.db && make migrate-up 

docker-dev:
	./scripts/build-docker.sh && docker compose up --build