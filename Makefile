.PHONY: run, env, create_db

# Загружаем .env
ifneq (,$(wildcard .env))
    include .env
    export
endif

run:
	@go run cmd/server/main.go || true

env:
	@cp .env.example .env

create_db:
	@echo "Creating database $(POSTGRES_NAME) as user $(POSTGRES_USER)..."
	@sudo -u ${POSTGRES_USER} createdb ${POSTGRES_NAME} || echo "Database may already exist"
