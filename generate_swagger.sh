#!/bin/bash

# Ensure you're in the correct project root directory
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$PROJECT_ROOT"


# Add verbose output for debugging
echo "Generating Swagger docs..."
echo "Project Root: $PROJECT_ROOT"
echo "Main Go File: $MAIN_PATH"

# Generate Swagger docs with verbose output
swag init -g "./cmd/api/main.go" -o ./docs --parseDependency --parseInternal --parseDepth 1