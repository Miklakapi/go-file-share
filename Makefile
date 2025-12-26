SHELL := /usr/bin/bash
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed -n 's/^\([^#][A-Za-z0-9_]*\)=.*/\1/p' .env)
endif

GO ?= go

.PHONY: help run-ram run-redis run-sqlite run-test

help:
	@echo "Targets:"
	@echo "  make run-ram      - run cmd/ram-app"
	@echo "  make run-redis    - run cmd/redis-app"
	@echo "  make run-sqlite   - run cmd/sqlite-app"
	@echo "  make run-test     - run cmd/test-app"

run-ram:
	$(GO) run ./cmd/ram-app/main.go

run-redis:
	$(GO) run ./cmd/redis-app/main.go

run-sqlite:
	$(GO) run ./cmd/sqlite-app/main.go

run-test:
	$(GO) run ./cmd/test-app/main.go