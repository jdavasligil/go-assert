.PHONY: all
all: fmt build test

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: build
build:
	@go build -o ./bin/dass ./cmd/delete_assertions.go

.PHONY: run
run:
	@./bin/das

.PHONY: test
test:
	@go test -v ./...

.PHONY: clean
clean:
	@go clean && rm -rf ./bin/*
