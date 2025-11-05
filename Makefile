DB_DSN = "postgres://postgres:postgres@localhost:5432/tech-e-market?sslmode=disable"
MIGRATIONS_DIR = db

# Apply all pending migrations
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) up

# Rollback the last migration
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) down-to 0

# Reset all migrations
migrate-reset:
	goose -dir $(MIGRATIONS_DIR) postgres $(DB_DSN) reset

# Create a new migration file
create-migration:
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

# Start server
start:
	go run cmd/main.go