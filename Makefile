include secret/.env

## connects to the database through the console
.PHONY: db/psql
db/psql:
	@echo "Connecting to the database...\nTo quit type \\q\n"
	@psql $(DBSTRING)

## runs the application
.PHONY: run/api
run/api:
	@go run cmd/api/main.go

## runs rabbitmq
.PHONY: amqp/run
amqp/run:
	docker run -d --name rabbit -p 15672:15672 -p 5672:5672 rabbitmq:3-management

## build/docs: generate API documentation using Swagger
.PHONY: build/docs
build/docs:
	@echo 'Building docs'
	swag init -g ./cmd/api/main.go

MIGRATIONS_DIR = ./migrations

.PHONY: migration/create
migration/create:
	@read -p "Enter migration name: " migration_name; \
	migrate create -seq -ext .sql -dir $(MIGRATIONS_DIR) $$migration_name

migration/up:
	@read -p "Enter the number of migrations to apply (ignore if you want to migrate up as possible): " count; \
	migrate -path $(MIGRATIONS_DIR) -database $(DBSTRING) up $$count

# Цель для отката миграций с подтверждением действия
migration/down:
	@echo "WARNING: You are about to roll back migrations!"
	@read -p "Are you sure you want to continue? (y/n): " confirm; \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		read -p "Enter the number of migrations to rollback (leave blank to rollback all migrations): " count; \
		migrate -path $(MIGRATIONS_DIR) -database $(DBSTRING) down $$count; \
		echo "Migrations rolled back successfully."; \
	else \
		echo "Migration rollback aborted."; \
	fi
