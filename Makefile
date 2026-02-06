# HRMS Backend Makefile
# Provides convenient commands for development and Docker operations

.PHONY: help build run dev stop clean logs test

# Default target
help:
	@echo "HRMS Backend - Available Commands"
	@echo "=================================="
	@echo ""
	@echo "Development:"
	@echo "  make dev-deps      - Start development dependencies (PostgreSQL, RabbitMQ)"
	@echo "  make dev-deps-stop - Stop development dependencies"
	@echo "  make run           - Run the application locally (requires Go)"
	@echo "  make build-local   - Build binary locally"
	@echo "  make test          - Run tests"
	@echo ""
	@echo "Docker (Full Stack):"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start all services (app + dependencies)"
	@echo "  make docker-down   - Stop all services"
	@echo "  make docker-logs   - View application logs"
	@echo "  make docker-clean  - Remove containers, volumes, and images"
	@echo "  make docker-rebuild- Rebuild and restart services"
	@echo ""
	@echo "Database:"
	@echo "  make db-shell      - Open PostgreSQL shell"
	@echo "  make db-dump       - Dump database to file"
	@echo "  make db-restore    - Restore database from file"
	@echo ""
	@echo "Utilities:"
	@echo "  make logs          - Tail all logs"
	@echo "  make status        - Show status of all containers"
	@echo "  make shell         - Open shell in app container"

# ==============================================================================
# Development Commands
# ==============================================================================

# Start only development dependencies (run your app locally)
dev-deps:
	docker compose -f docker-compose.dev.yml up -d
	@echo ""
	@echo "Development dependencies started!"
	@echo "PostgreSQL: localhost:5432"
	@echo "RabbitMQ:   localhost:5672 (Management: localhost:15672)"
	@echo ""
	@echo "Run 'make run' or 'go run ./cmd/api' to start the application"

dev-deps-stop:
	docker compose -f docker-compose.dev.yml down

dev-deps-clean:
	docker compose -f docker-compose.dev.yml down -v

# Run application locally (requires Go installed)
run:
	go run ./cmd/api

# Build binary locally
build-local:
	CGO_ENABLED=0 go build -ldflags="-w -s" -o ./bin/hrms-backend ./cmd/api

# Run tests
test:
	go test -v ./...

# ==============================================================================
# Docker Commands (Full Stack)
# ==============================================================================

# Build Docker image
docker-build:
	docker compose build

# Start all services
docker-up:
	docker compose up -d
	@echo ""
	@echo "All services started!"
	@echo "Application: http://localhost:8080"
	@echo "RabbitMQ Management: http://localhost:15672"

# Stop all services
docker-down:
	docker compose down

# View logs
docker-logs:
	docker compose logs -f app

# View all logs
logs:
	docker compose logs -f

# Clean up everything (containers, volumes, images)
docker-clean:
	docker compose down -v --rmi local
	@echo "Cleaned up containers, volumes, and local images"

# Rebuild and restart
docker-rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d

# Show status
status:
	docker compose ps

# Open shell in app container
shell:
	docker compose exec app sh

# ==============================================================================
# Database Commands
# ==============================================================================

# Open PostgreSQL shell
db-shell:
	docker compose exec postgres psql -U postgres -d hrms_db

# Dump database
db-dump:
	@mkdir -p ./backups
	docker compose exec postgres pg_dump -U postgres hrms_db > ./backups/hrms_db_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "Database dumped to ./backups/"

# Restore database (usage: make db-restore FILE=./backups/hrms_db_xxx.sql)
db-restore:
	@if [ -z "$(FILE)" ]; then echo "Usage: make db-restore FILE=./backups/hrms_db_xxx.sql"; exit 1; fi
	docker compose exec -T postgres psql -U postgres -d hrms_db < $(FILE)
	@echo "Database restored from $(FILE)"

# ==============================================================================
# Production Commands
# ==============================================================================

# Build production image with tag
docker-build-prod:
	docker build -t hrms-backend:latest .
	docker tag hrms-backend:latest hrms-backend:$(shell git rev-parse --short HEAD)

# Run in production mode
docker-prod:
	ENV=production docker compose up -d
