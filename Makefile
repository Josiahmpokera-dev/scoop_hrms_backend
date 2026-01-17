.PHONY: generate-models generate-users build run help sqlc-generate sqlc-verify

# Generate GORM models for a module
# Usage: make generate-models MODULE=users TABLE=users
generate-models:
	@if [ -z "$(MODULE)" ] || [ -z "$(TABLE)" ]; then \
		echo "Usage: make generate-models MODULE=<module_name> TABLE=<table_name>"; \
		echo "Example: make generate-models MODULE=users TABLE=users"; \
		exit 1; \
	fi
	@echo "Generating GORM models for $(MODULE) module (table: $(TABLE))..."
	go run scripts/generators/generate_models.go $(MODULE) $(TABLE)

# Generate users models (example)
generate-users:
	@echo "Generating users models..."
	go run scripts/generators/generate_users.go

# Generate SQLC code (optional)
sqlc-generate:
	@echo "Generating SQLC code..."
	sqlc generate

# Verify SQLC configuration
sqlc-verify:
	@echo "Verifying SQLC configuration..."
	sqlc vet

# Build the application
build:
	@echo "Building application..."
	go build -o bin/api ./cmd/api

# Run the application
run:
	@echo "Running application..."
	go run ./cmd/api/main.go

# Help
help:
	@echo "Available commands:"
	@echo "  make generate-models MODULE=<name> TABLE=<name>  - Generate GORM models"
	@echo "  make generate-users                             - Generate users models"
	@echo "  make build                                       - Build the application"
	@echo "  make run                                         - Run the application"
	@echo ""
	@echo "Optional (SQLC):"
	@echo "  make sqlc-generate                               - Generate SQLC code"
	@echo "  make sqlc-verify                                 - Verify SQLC config"
