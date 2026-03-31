#!/usr/bin/env bash
set -e

read -p "⚠️  This will DROP the database. Continue? (y/n): " confirm
if [ "$confirm" != "y" ]; then
  echo "❌ Aborted"
  exit 0
fi

if [ ! -f ".env" ]; then
  echo "❌ .env file not found"
  exit 1
fi

export $(grep -v '^#' .env | xargs)

echo "🔥 Dropping database $DATABASE_NAME..."

docker compose exec db psql \
  -U "$DATABASE_USER" \
  -c "DROP DATABASE IF EXISTS $DATABASE_NAME;"

docker compose exec db psql \
  -U "$DATABASE_USER" \
  -c "CREATE DATABASE $DATABASE_NAME;"

echo "♻️ Database recreated"

echo "🚀 Running migrations..."
go run ./cmd/migrate/main.go

echo "🌱 Running seed..."
go run ./cmd/seed/main.go

echo "✅ Reset DB completed"
