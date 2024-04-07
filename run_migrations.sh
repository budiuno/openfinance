#!/bin/bash

# Database connection parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=root
DB_PASSWORD=root
DB_NAME=openfinanceDB
SSL_MODE=disable

# Directory containing migration SQL files
MIGRATIONS_DIR=schema/migrations

# Path to golang-migrate binary
MIGRATE_BIN=$(which migrate)

# Check if golang-migrate binary exists
if [ -z "$MIGRATE_BIN" ]; then
  echo "Error: golang-migrate binary not found in PATH"
  exit 1
fi

# Run migrations
echo "Running migrations..."

# Loop through SQL files in migrations directory
for FILE in $MIGRATIONS_DIR/*.up.sql; do
  # Extract migration version from filename
  VERSION=$(basename $FILE | cut -d'_' -f1)
  
  # Run up migration
  $MIGRATE_BIN -database "postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$SSL_MODE" -path $MIGRATIONS_DIR up $VERSION
done

echo "Migrations completed"
