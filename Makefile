SHELL=/bin/sh

.PHONY: start-dev
start-dev:
	docker compose -f docker-compose.dev.yml up

.PHONY: build-migrator-image
build-migrator-image:
	docker build -t karten-backend-migrator:latest --target=migrator .

.PHONY: rm-dev-db
rm-dev-db:
	docker compose -f docker-compose.dev.yml rm -sv db
	docker volume rm karten-backend_pg_data

.PHONY: start-test-db
start-test-db:
	docker compose -f docker-compose.test.yml up -d karten-test-db

.PHONY: stop-test-db
stop-test-db:
	docker compose -f docker-compose.test.yml down -v --remove-orphans

.PHONY: test-db
test-db:
	STORE_DSN="postgres://tester:tester@127.0.0.1:5432/test?sslmode=disable&search_path=karten" \
		go test "github.com/lesnoi-kot/karten-backend/src/store"

.PHONY: docker-test-db
docker-test-db:
	docker compose -f docker-compose.test.yml up   \
		--abort-on-container-exit --force-recreate \
		--exit-code-from karten-intergation-tests || true
	docker compose -f docker-compose.test.yml down -v
