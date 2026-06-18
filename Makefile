BINARY=starter-api
CMD=./cmd/main.go

.PHONY: all build run dev tidy seed help

all: build

## Download dependencies and tidy go.mod
tidy:
	go mod tidy

## Build binary
build: tidy
	go build -o $(BINARY) $(CMD)

## Run directly (with hot-reload friendly output)
run:
	go run $(CMD)

## Run with Air hot-reload (requires: go install github.com/air-verse/air@latest)
dev:
	air

## Run tests
test:
	go test ./...

## Clean build artifacts
clean:
	rm -f $(BINARY)

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  tidy    Download and tidy dependencies"
	@echo "  build   Compile the binary"
	@echo "  run     Run the application"
	@echo "  dev     Run with Air hot-reload"
	@echo "  test    Run all tests"
	@echo "  clean   Remove build artifacts"
