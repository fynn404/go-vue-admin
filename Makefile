.PHONY: all build clean test run docker-build docker-run frontend-dev

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
BINARY_NAME=go-vue-admin
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build

build:
	cd cmd && $(GOBUILD) -o ../bin/$(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f bin/$(BINARY_NAME)
	rm -f bin/$(BINARY_UNIX)

test:
	$(GOTEST) -v ./...

run:
	cd cmd && $(GOBUILD) -o ../bin/$(BINARY_NAME) -v
	./bin/$(BINARY_NAME)

# Frontend commands
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME) .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Development helpers
dev: frontend-dev run

# Database migrations
migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

# Generate API documentation
gen-docs:
	swag init -g cmd/main.go

# Lint
lint:
	golangci-lint run

.DEFAULT_GOAL := build 