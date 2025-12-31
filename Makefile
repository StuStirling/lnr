.PHONY: build test lint clean install

# Build the binary
build:
	CGO_ENABLED=0 go build -o lnr .

# Run tests
test:
	CGO_ENABLED=0 go test ./...

# Run tests with coverage
test-coverage:
	CGO_ENABLED=0 go test -cover ./...

# Run linter
lint:
	golangci-lint run

# Clean build artifacts
clean:
	rm -f lnr

# Install to GOPATH/bin
install:
	go install .

# Run the CLI
run:
	go run . $(ARGS)

# Tidy dependencies
tidy:
	go mod tidy

# Format code
fmt:
	go fmt ./...

# All checks (test + lint)
check: test lint
