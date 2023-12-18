# Karten-backend

## Dev

```sh
git clone https://github.com/lesnoi-kot/karten-backend.git
make prepare

# Run the app in docker environment with live reload
make run

# Or run the app locally
make run-local
```

## Run intergation tests

```sh
make test-integration-docker

# Debugging integration tests
make start-test-db
make test-db
# ...
make stop-test-db
```

## TODO

- Functional tests with Hurl
- Contexted logging in a model operation boundary (eg: implicitly pass id of a project to logger in order to mantain and filter logs)
- minio s3
- testcontainers
- mockery
