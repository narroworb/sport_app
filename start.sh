#!/bin/bash

# Sports Analytics Platform Startup Script

echo "🚀 Starting Sports Analytics Platform..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "📦 Building images..."
docker-compose build

echo "🌐 Starting services..."
docker-compose up -d

echo "⏳ Waiting for services to be ready..."
sleep 10

echo ""
echo "✅ Services started successfully!"
echo ""
echo "📊 Frontend:          http://localhost"
echo "🔐 Auth Service:      http://localhost:8081"
echo "📈 Core API:          http://localhost:8080"
echo "📉 Analytics:         http://localhost:8082"
echo ""
echo "🗄️  Databases:"
echo "   PostgreSQL:        localhost:5432"
echo "   ClickHouse:        localhost:8123"
echo "   Redis:             localhost:6379"
echo "   Elasticsearch:     localhost:9200"
echo "   Kafka:             localhost:9092"
echo ""
echo "💡 To view logs: docker-compose logs -f"
echo "🛑 To stop:      docker-compose down"
echo ""
