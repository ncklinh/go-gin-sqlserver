#!/bin/bash
set -e

echo "🎬 Starting Postgres and Kafka in Docker for Dev..."

# Start all services (including Kafka)
docker compose -f docker-compose.dev.yml up -d

echo "⏳ Waiting for Postgres and Kafka to be ready..."
sleep 10

docker compose -f docker-compose.dev.yml ps
echo ""
echo "Postgres is running on port 5432"

# Create Kafka topic (film-events with 3 partitions)
echo "🛠️ Creating Kafka topic 'film-events' with 3 partitions..."
docker exec kafka kafka-topics.sh \
  --create \
  --if-not-exists \
  --topic film-events \
  --bootstrap-server localhost:9092 \
  --partitions 3 \
  --replication-factor 1
echo "✅ Kafka topic 'film-events' ready."

# Check if Air is installed
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

