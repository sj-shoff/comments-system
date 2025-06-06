.PHONY: run-inmemory run-postgres docker-inmemory docker-postgres docker-down migrate gqlgen

# Config paths
INMEMORY_CONFIG := ./configs/inmemory.yaml
POSTGRES_CONFIG := ./configs/postgres.yaml

# Run with in-memory storage (local)
run-inmemory:
	@echo "Starting in-memory storage..."
	CONFIG_PATH=$(INMEMORY_CONFIG) go run ./cmd/comments-system/main.go

# Run with PostgreSQL (local)
run-postgres:
	@echo "Starting PostgreSQL storage..."
	CONFIG_PATH=$(POSTGRES_CONFIG) POSTGRES_PASSWORD=$$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2) go run ./cmd/comments-system/main.go

# Docker with in-memory storage
docker-inmemory:
	@echo "Starting Docker with in-memory storage..."
	docker-compose -f docker/docker-compose.yaml down -v
	CONFIG_FILE=inmemory.yaml docker-compose -f docker/docker-compose.yaml up --build --remove-orphans

# Docker with PostgreSQL
docker-postgres:
	@echo "Starting Docker with PostgreSQL..."
	docker-compose -f docker/docker-compose.yaml down -v
	CONFIG_FILE=postgres.yaml POSTGRES_PASSWORD=$$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2) \
	docker-compose -f docker/docker-compose.yaml up --build --remove-orphans

# Stop Docker
docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f docker/docker-compose.yaml down -v

# Apply migrations
migrate:
	@echo "Applying migrations..."
	CONFIG_PATH=./configs/postgres.yaml POSTGRES_PASSWORD=$$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2) \
	go run ./cmd/migrator/main.go

# Generate GraphQL code
gqlgen:
	@echo "Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate --config ./internal/graph/gqlgen.yml