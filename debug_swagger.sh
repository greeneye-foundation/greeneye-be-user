#!/bin/bash

echo "🔍 Module Diagnosis Report"
echo "========================="

echo -e "\n📦 Module Name:"
go list -m

echo -e "\n🔗 Dependency Graph:"
go mod graph | grep -E "swagger|gin"

echo -e "\n📝 Import Paths:"
go list -m all | grep greeneye

echo -e "\n🧩 Swagger Dependencies:"
go list -m all | grep -E "swag|swagger"