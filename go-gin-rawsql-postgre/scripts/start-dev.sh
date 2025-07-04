#!/bin/bash
set -e

echo "🎬 Starting Postgres in Docker for Dev..."

# Start Postgres only
docker compose -f docker-compose.dev.yml up -d

# Wait for Postgres
echo "⏳ Waiting for Postgres to be ready..."
sleep 5

docker compose -f docker-compose.dev.yml ps

echo ""
echo "Postgres is running on port 5432"

# Check if air is installed
if ! command -v air &> /dev/null
then
    echo "⚠️  Air is not installed. Installing it now..."
    go install github.com/air-verse/air@latest

    echo ""
    echo "✅ Air installed. Please add your GOPATH/bin to PATH if needed:"
    echo "   export PATH=\$PATH:\$(go env GOPATH)/bin"
    echo "   (then re-run this script)"
    exit 1
fi

echo ""
echo "✅ Air is installed. Starting hot-reload server..."

# Actually run Air
# air

