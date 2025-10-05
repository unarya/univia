#!/bin/bash
set -e
echo "Generating swagger docs..."

cd "$(dirname "${BASH_SOURCE[0]}")/../../.."

mkdir -p api/swagger

swag init \
  -g cmd/api/main.go \
  --parseInternal \
  --parseDependency \
  --output api/swagger

echo "✅ Swagger docs generated at: $(pwd)/api/swagger"
