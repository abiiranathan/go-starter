# Load .env file
ifneq (,$(wildcard ./.env))
	include .env
	export $(shell sed 's/=.*//' .env)
endif

DATABASE_URL := $(shell echo $$DATABASE_URL)

MIGRATIONS = sqlc/migrations
TARGET = bin/server

# Build the app
build: generate
	CGO_ENABLED=0 go build -ldflags "-s -w" -o $(TARGET) cmd/server/*.go

# Install dependencies
install:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


# Create a new migration
create:
	@echo "Enter migration name: "; \
	read name; \
	if [ -z "$$name" ]; then \
		echo "Migration name is required"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir $(MIGRATIONS) -seq "$$name"

# Migrate up all migrations
up:
	migrate -source file://$(MIGRATIONS) -database $(DATABASE_URL) up

# force a specific version
force:
	migrate -source file://$(MIGRATIONS) -database $(DATABASE_URL) force $(VERSION)

# Migrate down all migrations
down:
	migrate -source file://$(MIGRATIONS) -database $(DATABASE_URL) down 1

# Generate sqlc and tailwindcss
generate: 
	sqlc generate
	bunx tailwindcss -i assets/static/css/input.css -o assets/static/css/styles.css

# Watch tailwindcss changes
watch: 
	bunx tailwindcss -i assets/static/css/input.css -o assets/static/css/styles.css --watch
	
# dev uses air to watch for changes and rebuild the app
dev:
	air -c .air.toml

.PHONY: up down generate build create install dev 
