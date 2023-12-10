all: start

start:
	go run ./cmd/main.go input.txt

test:
	go test ./... | { grep -v 'no test files'; true; }