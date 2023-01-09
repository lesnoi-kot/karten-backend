SHELL=/bin/sh

start-test-db:
	docker compose -f docker-compose.test.yml up -d karten-test-db

stop-test-db:
	docker compose -f docker-compose.test.yml down -v --remove-orphans

test-db:
	STORE_DSN="postgres://tester:tester@127.0.0.1:5432/test?sslmode=disable&search_path=karten" \
		go test "github.com/lesnoi-kot/karten-backend/src/store"

docker-test-db:
	docker compose -f docker-compose.test.yml up   \
		--abort-on-container-exit --force-recreate \
		--exit-code-from karten-intergation-tests || true
	docker compose -f docker-compose.test.yml down -v
