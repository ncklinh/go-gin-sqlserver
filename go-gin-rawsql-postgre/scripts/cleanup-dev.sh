#!/bin/bash

# Film Rental API Development Cleanup Script
set -e

echo "🧹 Cleaning up Film Rental API Development environment..."

# Stop and remove containers, networks, and volumes
echo "📦 Removing development containers, networks, and volumes..."
docker-compose -f docker-compose.dev.yml down -v

# Remove only our project's development Docker image
echo "🗑️  Removing development Docker image..."
docker rmi go-gin-rawsql-postgre-app:latest 2>/dev/null || echo "No development app image to remove"
docker rmi $(docker images -q go-gin-rawsql-postgre-app) 2>/dev/null || echo "No development app images found"

# Remove only our project's development volumes
echo "💾 Cleaning up development volumes..."
docker volume rm go-gin-rawsql-postgre_postgres_data_dev 2>/dev/null || echo "No development volumes to remove"

# Remove only our project's development networks
echo "🌐 Cleaning up development networks..."
docker network rm go-gin-rawsql-postgre_film_rental_network_dev 2>/dev/null || echo "No development networks to remove"

# Remove only our project's development containers (if any are still running)
echo "📋 Removing development containers..."
docker rm -f film_rental_app_dev film_rental_db_dev 2>/dev/null || echo "No development containers to remove"

echo ""
echo "✅ Development cleanup completed successfully!"
echo ""
echo "🎬 To start fresh in development mode:"
echo "   ./scripts/start-dev.sh"
echo ""
echo "💡 Only Film Rental API development Docker resources have been removed." 