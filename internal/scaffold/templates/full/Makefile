.PHONY: dev migrate seed build clean

ENV_FILE := .env.local
APP_NAME := app
BUILD_DIR := bin


COUNT ?=
SEED_ARGS :=
ifneq ($(strip $(COUNT)),)
  SEED_ARGS += --count=$(COUNT)
endif

ifeq ($(OS),Windows_NT)
SHELL := powershell.exe
.SHELLFLAGS := -NoProfile -ExecutionPolicy Bypass -Command
AIR := air
LOAD_ENV := ./scripts/load-env.ps1 -EnvFile $(ENV_FILE);
else
SHELL := /bin/bash
.SHELLFLAGS := -lc
AIR := air
LOAD_ENV := set -a; [ -f $(ENV_FILE) ] && . $(ENV_FILE); set +a;
endif

docker:
	docker compose up -d

dev:
	@echo "🚀 Running with Air..."
	@$(LOAD_ENV) $(AIR)

migrate:
	@$(LOAD_ENV) go run ./cmd/migrate

migrate-reset:
	@$(LOAD_ENV) go run ./cmd/migrate --reset

seed:

	@$(LOAD_ENV) go run ./cmd/seed $(SEED_ARGS)

build:
	@echo "🔨 Building binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

clean:
	@rm -rf $(BUILD_DIR)
