#!/bin/bash
set -e

# This script runs when the PostgreSQL container starts for the first time
# It sets up the initial database configuration

echo "Initializing PostgreSQL database..."

# Wait for PostgreSQL to be ready
until pg_isready -U $POSTGRES_USER -d $POSTGRES_DB; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 2
done

echo "PostgreSQL is ready!"

# Grant necessary permissions to the user
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Grant all privileges on database
    GRANT ALL PRIVILEGES ON DATABASE film_rental TO filmuser;
    
    -- Grant usage on schema
    GRANT USAGE ON SCHEMA public TO filmuser;
    
    -- Grant all privileges on all tables (will be created by the Go app)
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO filmuser;
    
    -- Grant all privileges on all sequences (will be created by the Go app)
    GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO filmuser;
    
    -- Set default privileges for future tables
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO filmuser;
    
    -- Set default privileges for future sequences
    ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO filmuser;
EOSQL

echo "Database initialization completed!" 