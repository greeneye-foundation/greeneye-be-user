.PHONY: dev build test clean

# Development with hot reload
dev:
	air

# Build the application
build:
	go build -o bin/app cmd/api/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf tmp/

# Install dependencies
deps:
	go mod download
	go mod tidy

# Update dependencies
update:
	go get -u all
	go mod tidy

# Generate API documentation (if you're using swag)
docs:
	swag init -g cmd/api/main.go -o api/docs