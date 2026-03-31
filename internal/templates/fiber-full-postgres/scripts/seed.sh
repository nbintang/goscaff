#!/usr/bin/env bash
set -e

echo "🌱 Running database seed..."

if [ ! -f ".env" ]; then
  echo "❌ .env file not found"
  exit 1
fi

export $(grep -v '^#' .env | xargs)

go run ./cmd/seed/main.go

echo "✅ Seeding completed"
