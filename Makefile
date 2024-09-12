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
