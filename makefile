.PHONY: test
test:
	gotestsum --watch

.PHONY: test-ci
test-ci:
	go test ./...

.PHONY: migrate-up
migrate-up:
	goose -dir "./sql/migrations" sqlite ./db/database.db up

.PHONY: migrate-down
migrate-down:
	goose -dir "./sql/migrations" sqlite ./db/database.db down

.PHONY: reset-db
reset-db:
	rm -f ./db/database.db && go run ./cmd/migrator/main.go ./db/database.db

.PHONY: dev
dev: 
	templ generate --watch --proxy="http://localhost:3000" --cmd="go run ./cmd/server"

.PHONY: migrate
migrate:
	go run ./cmd/migrator/main.go

.PHONY: build-server
build-server:
	docker/server/build.sh

.PHONY: deploy-server
deploy-server:
	fly deploy --config docker/server/fly.toml --dockerfile docker/server/Dockerfile

.PHONY: build-scraper
build-scraper:
	docker/scraper/build.sh

.PHONY: deploy-scraper
deploy-scraper:
	fly deploy --config docker/scraper/fly.toml --dockerfile docker/scraper/Dockerfile

.PHONY: build-all
build-all: build-server build-scraper

.PHONY: deploy-all
deploy-all: deploy-server deploy-scraper

.PHONY: start-services
start-services:
	docker-compose up --build

.PHONY: dev
dev: build-all start-services