



all: server cli

.PHONY: vet
vet:
	go vet ./...

.PHONY: server

server: vet
	go build -o bin/go-redis ./cmd/server

.PHONY: cli
cli: vet
	go build -o bin/go-redis-cli ./cmd/cli

.PHONY: test
test:
	go test -cover ./...

image:
	docker build -t go-redis .
