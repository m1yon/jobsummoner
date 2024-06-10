test:
	gotestsum --watch

test-ci:
	go test ./...

migrate-up:
	goose -dir "./sql/migrations" sqlite ./db/database.db up

migrate-down:
	goose -dir "./sql/migrations" sqlite ./db/database.db down

reset-db:
	rm -f ./db/database.db && go run ./cmd/migrator/main.go ./db/database.db

dev: 
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/server"

migrate:
	go run ./cmd/migrator/main.go

deploy-server:
	docker/server/build.sh && fly deploy --config docker/server/fly.toml --dockerfile docker/server/Dockerfile

deploy-scraper:
	docker/scraper/build.sh && fly deploy --config docker/scraper/fly.toml --dockerfile docker/scraper/Dockerfile

deploy: deploy-server deploy-scraper