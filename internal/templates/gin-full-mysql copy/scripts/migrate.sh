#!/usr/bin/env bash
set -e

echo "🚀 Running database migrations..."

if [ ! -f ".env" ]; then
  echo "❌ .env file not found"
  exit 1
fi

export $(grep -v '^#' .env | xargs)

go run ./cmd/migrate/main.go

echo "✅ Migration completed"
