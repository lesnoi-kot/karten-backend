ARG GO_VERSION="1.21"

FROM golang:${GO_VERSION}-alpine as deps
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

FROM deps as builder
COPY src ./src/
RUN GOOS=linux go build -o karten -ldflags "-s -w" ./src/cmd/karten/main.go

FROM deps as migratorBuilder
RUN mkdir -p src/cmd/migrator
COPY src/cmd/migrator src/cmd/migrator
RUN GOOS=linux go build -o migrator -ldflags "-s -w" ./src/cmd/migrator/main.go

FROM alpine:3.18 as app
WORKDIR /root
EXPOSE 4000
COPY --from=builder /app/karten karten
ENTRYPOINT ["/root/karten"]

FROM alpine:3.18 as migrator
WORKDIR /root
COPY --from=migratorBuilder /app/migrator migrator
ENTRYPOINT ["/root/migrator"]
CMD ["db"]
