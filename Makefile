.PHONY: all
all: fmt build test

.PHONY: fmt
fmt:
	@go fmt ./...

.PHONY: build
build:
	@go build -o ./bin/dass ./cmd/dass/dass.go

.PHONY: run
run:
	@./bin/dass

.PHONY: test
test:
	@go test -v ./...

.PHONY: clean
clean:
	@go clean && rm -rf ./bin/*
