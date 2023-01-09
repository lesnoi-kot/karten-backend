# Karten-backend

## Dev

```sh
# Spin up a local database
docker compose -f docker-compose.dev.yml up

# Run the app
DEBUG=1 API_HOST="127.0.0.1:4000" STORE_DSN="postgres://karten:karten@127.0.0.1:5432/karten?sslmode=disable&search_path=karten" go run src/main.go

docker compose -f docker-compose.dev.yml down -v
```

## Run intergation tests

```sh
make docker-test-db

# Debugging integration tests
make start-test-db
make test-db
# ...
make stop-test-db
```

## TODO

- Functional tests with Hurl
