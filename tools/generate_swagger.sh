#!/bin/bash

# Generate Swagger docs with verbose output
swag init -g "./cmd/api/main.go" -o ./docs --parseDependency --parseInternal --parseDepth 1