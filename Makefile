GO ?= go
APP ?= goscaff
ARGS ?=
ROOT ?= ./internal/templates
OLD_MODULE ?=
RENAME_DIR ?= ./internal/templates

.PHONY: help cmd cli tools renames renames-dry templater build build-cli build-tools install install-cli install-tools install-renames install-templater

help:
	@echo "Available targets:"
	@echo "  make cmd ARGS=\"new myapp\""
	@echo "  make cli ARGS=\"new myapp\""
	@echo "  make renames RENAME_DIR=./internal/templates/fiber-full-postgres"
	@echo "  make renames-dry RENAME_DIR=./internal/templates/fiber-full-postgres"
	@echo "  make templater ROOT=./internal/templates OLD_MODULE=github.com/nbintang/old"
	@echo "  make build"
	@echo "  make build-tools"
	@echo "  make install"
	@echo "  make install-tools"

cmd:
	$(GO) run . $(ARGS)

cli: cmd

tools:
	@echo "Tool targets:"
	@echo "  make renames RENAME_DIR=./internal/templates/fiber-full-postgres"
	@echo "  make renames-dry RENAME_DIR=./internal/templates/fiber-full-postgres"
	@echo "  make templater ROOT=./internal/templates OLD_MODULE=github.com/nbintang/old"
	@echo "  make install-renames"
	@echo "  make install-templater"

renames:
	$(GO) run ./cmd/tools/renames $(RENAME_DIR)

renames-dry:
	$(GO) run ./cmd/tools/renames $(RENAME_DIR) --dry-run

templater:
ifndef OLD_MODULE
	$(error OLD_MODULE is required. Example: make templater ROOT=./internal/templates OLD_MODULE=github.com/nbintang/old)
endif
	$(GO) run ./cmd/tools/templater $(ROOT) $(OLD_MODULE)

build: build-cli build-tools

build-cli:
	$(GO) build -o ./bin/$(APP) .

build-tools:
	$(GO) build -o ./bin/renames ./cmd/tools/renames
	$(GO) build -o ./bin/templater ./cmd/tools/templater

install: install-cli

install-cli:
	$(GO) install .

install-tools: install-renames install-templater

install-renames:
	$(GO) install ./cmd/tools/renames

install-templater:
	$(GO) install ./cmd/tools/templater
