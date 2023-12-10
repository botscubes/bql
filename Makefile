all: start

start:
	go run ./cmd/main.go input.txt

build:
	go build ./cmd/main.go

test:
	go test ./... | { grep -v 'no test files'; true; }