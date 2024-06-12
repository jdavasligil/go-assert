.PHONY: all
all: build run

.PHONY: build
build:
	@go build -o ./bin/delete_assertions ./cmd/delete_assertions.go

.PHONY: run
run:
	@./bin/delete_assertions

.PHONY: test
test:
	@go test -v ./...

.PHONY: clean
clean:
	@go clean && rm -rf ./bin/*
