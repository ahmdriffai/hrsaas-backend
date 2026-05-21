BINARY   := hrsaas
PKG      := ./...
COVERAGE := coverage.out

.PHONY: test test-v test-race test-cover cover-html test-pkg lint build run clean

## Run all tests
test:
	go test $(PKG)

## Run all tests with verbose output
test-v:
	go test -v $(PKG)

## Run all tests with race detector
test-race:
	go test -race $(PKG)

## Run tests and generate coverage report
test-cover:
	go test -coverprofile=$(COVERAGE) $(PKG)
	go tool cover -func=$(COVERAGE)

## Open coverage report in browser
cover-html: test-cover
	go tool cover -html=$(COVERAGE)

## Run tests for a specific package  (usage: make test-pkg PKG=./internal/usecase/...)
test-pkg:
	go test -v $(PKG)

## Run the linter (requires golangci-lint)
lint:
	golangci-lint run ./...

## Build the binary
build:
	go build -o $(BINARY) ./cmd/...

## Run the app
run:
	go run ./cmd/...

## Remove build artifacts and coverage output
clean:
	rm -f $(BINARY) $(COVERAGE)
