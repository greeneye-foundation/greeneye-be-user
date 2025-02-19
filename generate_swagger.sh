#!/bin/bash
set -e

# Validate Go environment
go version
go env GO111MODULE

# Clean previous swagger docs
rm -rf docs/swagger/*

# Generate Swagger Docs with comprehensive options
swag init \
    -g ./cmd/api/main.go \
    --parseDependency \
    --parseInternal \
    --generatedTime \
    --output ./docs/swagger \
    --outputTypes go,json,yaml
