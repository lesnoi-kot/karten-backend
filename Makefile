SHELL=/bin/sh

###
prepare:
	go mod tidy
	cp .env.example .env

###
run:
	docker compose up

###
run-local:
	USE_DOTENV=true go run ./src/cmd/karten/main.go

run-services:
	docker compose up db migrator static-server

###
test:
	go test ./...

test-integration:
	INTEGRATION_TESTS=1 go test -run ^TestIntegration ./...

docker-test-db:
	docker compose -f docker-compose.test.yml up --abort-on-container-exit

###
rm-dev-db:
	docker compose rm -sv db
	docker volume rm karten-backend_pg_data

rm-test-db:
	docker compose -f docker-compose.test.yml rm -sv karten-test-db

###
new-migration:
	go run src/cmd/migrator/main.go db create_sql unnamed_migration

migrate:
	go run src/cmd/migrator/main.go db migrate

rollback:
	go run src/cmd/migrator/main.go db rollback

###
build-migrator-image:
	docker build -t karten-backend-migrator:local --target=migrator .
