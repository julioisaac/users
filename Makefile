.PHONY: help run docker/up docker/down docker/restar tools swagger upgrade/mod install integration/test unit/test lint

DB_DSN=host=localhost port=5432 user=user-app password=user-pass dbname=user-dev-db sslmode=disable connect_timeout=2 statement_timeout=2s

_goose_ = goose -dir ./migration postgres "$(DB_DSN)"

GO = go
GOLANGCILINT = golangci-lint

# Load envs from .env file if it exists
-include: .env
export


# To add description to a target, just put a comment with two # after the target definition
# Ex:
# target_name: target_dep1 target_dep2  ## i'm a description
# do anything

help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make [target]\033[36m\033[0m\n\nTargets:\n"} /^[a-zA-Z_/-]+:.*?##/ { printf "\033[36m%-18s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

run:  swagger ## Run and serve an api
	@$(GO) run main.go api

docker/up:  ## Build and run docker apps, pass profile as argument PROFILE="--profile development
	@docker compose --env-file .docker-compose.env $(PROFILE) up -d
	until docker exec user-postgres-1 pg_isready ; do sleep 1 ; done

docker/down:  ## Stop docker apps, pass profile as argument PROFILE="--profile development"
	@docker compose $(PROFILE) down

docker/restart: docker/down docker/up db/migrate-up ## Restart docker apps, pass profile as argument PROFILE="--profile development"

tools:  ## Install golangci-lint and swaggo
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@$(GO) install github.com/swaggo/swag/cmd/swag@latest
	@$(GO) install github.com/pressly/goose/v3/cmd/goose@latest

install: tools  ## Install go libs
	@$(GO) mod vendor && $(GO) mod tidy

upgrade/mod: ## upgrade dependencies
	@$(GO) get -u -t ./...

build: lint ## Build apps
	@$(GO) build -o /dev/null

lint:  ## Run lint tools
	@$(GOLANGCILINT) run

db/create-migration: ## Create a migration file (need to pass MIGRATION_NAME param)
	$(_goose_) create $(MIGRATION_NAME) sql

db/migrate-up: ## Apply migrations in db
	$(_goose_) up

db/migrate-down: ## Apply migrations rollback in db
	$(_goose_) down

db/migration-status: ## Check migrations status in db
	$(_goose_) status