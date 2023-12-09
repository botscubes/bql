all: start

start:
	go run ./cmd/main.go

test:
	go test ./... | { grep -v 'no test files'; true; }