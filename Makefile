# Makefile — atalhos para o RCM Backoffice

mnum ?= 1

ifneq (,$(wildcard .env))
  include .env
  export
endif

.PHONY: help dev up down logs backend frontend test lint build

help:
	@echo "Targets principais:"
	@echo "  dev        - docker-compose up (build) + logs follow"
	@echo "  up         - docker-compose up -d (build)"
	@echo "  down       - docker-compose down"
	@echo "  logs       - docker-compose logs -f"
	@echo "  backend    - go run ./backend/cmd/server"
	@echo "  frontend   - (cd frontend && npm run dev)"
	@echo "  test       - (cd backend && go test ./...)"
	@echo "  lint       - (cd backend && go vet ./...)"
	@echo "  build      - docker-compose build"
	@echo "  build-be   - build da imagem backend"
	@echo "  build-fe   - build da imagem frontend"
	@echo "  migrate-create - cria nova migration (name=...)"
	@echo "  migrate-up - aplica todas as migrations"
	@echo "  migrate-down - desfaz migrations (mnum=1)"
	@echo "  migrate-up-force - forca versao da migration (version=...)"
dev: down up logs
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
	cd backend && go test ./...
lint:
	cd backend && go vet ./...
build: build-be	build-fe
migrate-create:
	docker run --rm \
		--network rgps.backoffice.network \
		-v $(PWD)/migrations:/migrations \
		--env-file .env \
		migrate/migrate:4 \
		--source file:///migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable" \
		create -ext sql -dir /migrations ${name}
		sudo chown "$(USER)":"$(USER)" -R $(PWD)/migrations
migrate-up:
	docker run --rm \
		--network rgps.backoffice.network \
		-v $(PWD)/migrations:/migrations \
		--env-file .env \
		migrate/migrate:4 \
		--source file:///migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable" \
		up
migrate-down:
	docker run --rm \
		--network rgps.backoffice.network \
		-v $(PWD)/migrations:/migrations \
		--env-file .env \
		migrate/migrate:4 \
		--source file:///migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable" \
		down $(mnum)
migrate-up-force:
	docker run --rm \
		--network rgps.backoffice.network \
		-v $(PWD)/migrations:/migrations \
		--env-file .env \
		migrate/migrate:4 \
		--source file:///migrations \
		-database "postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@db:5432/$(POSTGRES_DB)?sslmode=disable" \
		force $(version)
build-be:
	docker buildx build -t rgps.backoffice/backend:latest -f backend/Dockerfile backend
build-fe:
	docker buildx build -t rgps.backoffice/frontend:latest -f frontend/Dockerfile frontend
