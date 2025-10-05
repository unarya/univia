#!/bin/bash
set -e

# Ensure the main binary is built before running
if [ ! -f "./bin/api" ]; then
  echo "Binary './api' not found. Building it now..."
  go build -buildvcs=false -o ./bin/api .
fi

# Use environment variables (with default fallbacks)
DB_USER=${DB_USER:-admin}
DB_PASS=${DB_PASS:-MysqlAdmin}
DB_HOST=${DB_HOST:-db}
DB_PORT=${DB_PORT:-3306}
DB_NAME=${DB_NAME:-UniviaDB}

# Compose migration URL dynamically
DB_URL="mysql://${DB_USER}:${DB_PASS}@tcp(${DB_HOST}:${DB_PORT})/${DB_NAME}"

echo "✅ Running migrations with database URL: ${DB_URL}"
migrate -path ./migrations -database "${DB_URL}" up

# Start the application with Air
echo "🚀 Starting the application with Air..."
exec air
