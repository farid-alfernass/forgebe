.PHONY: build test fmt lint vet tidy

build:
	go build -o bin/forgebe ./cmd/forgebe

test:
	go test ./...

fmt:
	gofmt -w ./cmd ./internal

lint:
	golangci-lint run ./...

vet:
	go vet ./...

tidy:
	go mod tidy
