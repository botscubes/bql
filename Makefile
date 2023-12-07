all: start

start:
	go run ./cmd/main.go

test:
	go test ./...