#!/bin/bash

# Film Rental API Development Cleanup Script
set -e

echo "🧹 Cleaning up Film Rental API Development environment..."

# Stop and remove containers, volumes, and networks
echo "📦 Removing development containers, volumes, and networks..."
docker-compose -f docker-compose.dev.yml down -v --remove-orphans

# Remove Kafka container (if still exists)
echo "📋 Removing leftover containers..."
docker rm -f kafka film_rental_db_dev film_rental_app_dev 2>/dev/null || echo "No leftover containers to remove"

docker exec -it kafka kafka-topics.sh \
  --delete \
  --topic film-events \
  --bootstrap-server localhost:9092
  
# Remove custom app image (optional)
echo "🗑️  Removing development Docker image..."
docker rmi go-gin-rawsql-postgre-app:latest 2>/dev/null || echo "No development app image to remove"
docker rmi $(docker images -q go-gin-rawsql-postgre-app) 2>/dev/null || echo "No other development app images found"

# Remove development volumes
echo "💾 Removing development volumes..."
docker volume rm film_rental_postgres_data_dev 2>/dev/null || echo "No Postgres volume to remove"
docker volume rm $(docker volume ls -q --filter name=kafka) 2>/dev/null || echo "No Kafka volumes to remove"

# Remove development networks
echo "🌐 Removing development networks..."
docker network rm film_rental_network_dev 2>/dev/null || echo "No dev network to remove"

echo ""
echo "✅ Development cleanup completed successfully!"
echo ""
echo "🎬 To start fresh in development mode:"
echo "   ./scripts/start-dev.sh"
echo ""
echo "💡 All Film Rental API dev Docker resources (Postgres + Kafka) are now clean."
