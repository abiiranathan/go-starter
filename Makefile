# Load .env file
ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

DATABASE_URL := $(shell echo $$DATABASE_URL)

ifeq ($(DATABASE_URL),)
$(error DATABASE_URL is not set)
endif

MIGRATIONS = cmd/app/sqlc/migrations
TARGET = bin/app

install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

create:
	@echo "Enter migration name: "; \
	read name; \
	if [ -z "$$name" ]; then \
		echo "Migration name is required"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir $(MIGRATIONS) -seq "$$name"

up:
	migrate -source file://$(MIGRATIONS) -database $(DATABASE_URL) up

force:
	migrate -source file://$(MIGRATIONS) -database $(DATABASE_URL) force $(VERSION)

down:
	migrate -source file://$(MIGRATIONS) -database $(DATABASE_URL) down 1

generate: 
	sqlc generate

build: generate
	go build -ldflags "-s -w" -o $(TARGET) cmd/app/*.go

watch: 
	npx tailwindcss -i cmd/app/static/css/input.css -o cmd/app/static/css/styles.css --watch
	
dev:
	air -c .air.toml

.PHONY: up down generate build create install dev 
