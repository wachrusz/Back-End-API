include secret/.env

## connects to the database through the console
.PHONY: db/psql
db/psql:
	@echo "Connecting to the database...\nTo quit type \\q\n"
	@psql postgres://cadvadmin:${PASSWORD_DB}@${HOST}:${PORT}/cadvdb?sslmode=disable
