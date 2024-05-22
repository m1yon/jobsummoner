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
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/jobsummoner"

dev-show: 
	go run ./cmd/jobsummoner -rod=show

reset-db:
	rm -f ./db/database.db && make migrate-up 

docker-dev:
	./scripts/build-docker.sh && docker compose up --build

query-db:
	sqlite3 db/database.db

test:
	go test ./...