#!/usr/bin/env bash
set -e

# Initialize multiple databases for Duangdee microservices
echo "Initializing databases: auth_db, tarot_db, reading_db, payment_db, notification_db..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE auth_db;
    CREATE DATABASE tarot_db;
    CREATE DATABASE reading_db;
    CREATE DATABASE payment_db;
    CREATE DATABASE notification_db;
EOSQL

echo "Multiple databases initialized successfully!"
