SHELL := /usr/bin/bash
.ONESHELL:
.SHELLFLAGS := -euo pipefail -c

ifneq ("$(wildcard .env)","")
	include .env
	export $(shell sed -n 's/^\([^#][A-Za-z0-9_]*\)=.*/\1/p' .env)
endif

GO ?= go

.PHONY: help run-ram run-redis run-sqlite build-ram build-redis build-sqlite

help:
	@echo "Targets:"
	@echo "  make run-ram"
	@echo "  make run-redis"
	@echo "  make run-sqlite"
	@echo "  make build-ram"
	@echo "  make build-redis"
	@echo "  make build-sqlite"

run-ram:
	$(GO) run ./cmd/ram-app/main.go

run-redis:
	$(GO) run ./cmd/redis-app/main.go

run-sqlite:
	$(GO) run ./cmd/sqlite-app/main.go

build-ram:
	go build -o ./bin/ram-app ./cmd/ram-app/main.go

build-redis:
	go build -o ./bin/redis-app ./cmd/redis-app/main.go

build-sqlite:
	go build -o ./bin/sqlite-app ./cmd/sqlite-app/main.go
