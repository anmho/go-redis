


FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go vet ./...
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-redis /app/cmd/server


FROM alpine

WORKDIR /app

COPY --from=builder /app/go-redis /app/

EXPOSE 6379
CMD ["/app/go-redis"]





