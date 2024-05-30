test:
	gotestsum --watch

test-ci:
	go test ./

migrate-up:
	goose -dir "./sql/migrations" sqlite ./db/database.db up

migrate-down:
	goose -dir "./sql/migrations" sqlite ./db/database.db down

reset-db:
	rm -f ./db/database.db && make migrate-up 