.PHONY: run-inmemory run-postgres docker-inmemory docker-postgres docker-down migrate gqlgen

INMEMORY_CONFIG := ./configs/inmemory.yaml
POSTGRES_CONFIG := ./configs/postgres.yaml

# (local)
run-inmemory:
	@echo "Starting in-memory storage..."
	CONFIG_PATH=$(INMEMORY_CONFIG) go run ./cmd/comments-system/main.go

# (local)
run-postgres:
	@echo "Starting PostgreSQL storage..."
	CONFIG_PATH=$(POSTGRES_CONFIG) POSTGRES_PASSWORD=$$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2) go run ./cmd/comments-system/main.go

docker-inmemory:
	@echo "Starting Docker with in-memory storage..."
	docker-compose -f docker/docker-compose.yaml down -v
	CONFIG_FILE=inmemory.yaml STORAGE_TYPE=inmemory \
	docker-compose -f docker/docker-compose.yaml up --build --remove-orphans --scale postgres=0

docker-postgres:
	@echo "Starting Docker with PostgreSQL..."
	docker-compose -f docker/docker-compose.yaml down -v
	CONFIG_FILE=postgres.yaml STORAGE_TYPE=postgres POSTGRES_PASSWORD=$$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2) \
	docker-compose -f docker/docker-compose.yaml up --build --remove-orphans

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f docker/docker-compose.yaml down -v

migrate:
	@echo "Applying migrations..."
	CONFIG_PATH=./configs/postgres.yaml POSTGRES_PASSWORD=$$(grep POSTGRES_PASSWORD .env | cut -d '=' -f2) \
	go run ./cmd/migrator/main.go

gqlgen:
	@echo "Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate --config ./internal/graph/gqlgen.yml

test:
	go test -v ./... | grep -v 'no test files'