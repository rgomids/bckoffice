# Makefile — atalhos para o RCM Backoffice

.PHONY: help dev up down logs backend frontend test lint build

help:
	@echo "Targets principais:"
	@echo "  dev        - docker-compose up (build) + logs follow"
	@echo "  up         - docker-compose up -d (build)"
	@echo "  down       - docker-compose down"
	@echo "  logs       - docker-compose logs -f"
	@echo "  backend    - go run ./backend/cmd/server"
	@echo "  frontend   - (cd frontend && npm run dev)"
	@echo "  test       - go test ./..."
	@echo "  lint       - go vet ./..."
	@echo "  build      - docker-compose build"

dev: up logs

up:
	docker-compose -f infra/docker-compose.yml up -d --build

down:
	docker-compose -f infra/docker-compose.yml down

logs:
	docker-compose -f infra/docker-compose.yml logs -f

backend:
	go run ./backend/cmd/server

frontend:
	cd frontend && npm run dev

test:
	go test ./...

lint:
	go vet ./...

build:
	docker-compose -f infra/docker-compose.yml build
