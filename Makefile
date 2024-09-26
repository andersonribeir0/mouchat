# Simple Makefile for a Go project

# Build the application
all: build

build:
	@echo "Building..."
	@pnpx tailwindcss -i internal/web/views/css/app.css -o internal/web/public/styles.css --config ./internal/web/tailwind.config.js
	@templ generate internal/web/views
	@mkdir -p bin && go build -o bin/mouchat cmd/main.go

# Run the application
run:
	@go run cmd/main.go

# Create DB container
docker-run:
	@while read line; do export $line; done < .env
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Run migration up 
migrate-up:
	@go run cmd/migrate/main.go -migrate up

# Run migration down
migrate-down:
	@go run cmd/migrate/main.go -migrate down

# Create migration
create-migration:
	@goose -dir=cmd/migrate/migrations create $(NAME) sql

# Reset database
reset:
	@go run cmd/reset/main.go

# Live Reload
watch:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/cosmtrek/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

.PHONY: all build run test clean
