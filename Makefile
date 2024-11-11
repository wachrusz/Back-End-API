include secret/.env

## connects to the database through the console
.PHONY: db/psql
db/psql:
	@echo "Connecting to the database...\nTo quit type \\q\n"
	@psql postgres://cadvadmin:${DB_PASSWORD}@${HOST}:5432/cadvdb?sslmode=disable

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
