
.PHONY: help apigen migrate revision seed

help: ## Prints help for targets with comments
	@cat $(MAKEFILE_LIST) | grep -E '^[a-zA-Z_-]+:.*?## .*$$' | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

apigen:   ## Generate gRPC API code
	@echo "Generating gRPC API code..."
	@protoc --go_out=gen/go/grpc/v1 --go-grpc_out=gen/go/grpc/v1 api/grpc/v1/*.proto

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@migrate -path ./migrations -database ${DB_URL} up

revision: ## Create a new database migration
	@if [ -z "$(name)" ]; then echo "Error: name parameter is required. Usage: make revision name=your_migration_name"; exit 1; fi
	@echo "Creating new migration: $(name)"
	@migrate create -ext sql -dir ./migrations -seq $(name)

seed: ## Seed the database with initial data
	@echo "Seeding the database..."
	@go run migrations/seeder.go