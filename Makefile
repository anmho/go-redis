



all: server cli

.PHONY: server
server:
	go build -o bin/go-redis ./cmd/server

.PHONY: cli
cli:
	go build -o bin/go-redis-cli ./cmd/cli

.PHONY: test
test:
	go test ./...