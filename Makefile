.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down swagger

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $1, $2}'

build: ## Build the application
	go build -o bin/app cmd/app/main.go

run: ## Run the application
	go run cmd/app/main.go

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/

docker-up: ## Start Docker containers
	docker compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-rebuild: ## Rebuild and restart Docker containers
	docker-compose down
	docker-compose up -d --build

docker-logs: ## Show Docker logs
	docker-compose logs -f

mod-download: ## Download Go modules
	go mod download

mod-tidy: ## Tidy Go modules
	go mod tidy

seed: ## Seed database with initial data
	go run cmd/seed/main.go

install: ## Install dependencies
	go mod download
	go mod verify
	go install github.com/swaggo/swag/cmd/swag@latest

dev: ## Run in development mode
	GIN_MODE=debug go run cmd/app/main.go

swagger: ## Generate swagger documentation
	swag init -g cmd/app/main.go -o docs

swagger-serve: ## Generate and serve swagger docs
	swag init -g cmd/app/main.go -o docs
	@echo "Swagger docs generated. Run 'make dev' and visit http://localhost:8080/swagger/index.html"
