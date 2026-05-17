#!/bin/bash

# Sports Analytics Platform - Stop and Cleanup

echo "🛑 Stopping Sports Analytics Platform..."

docker-compose down

echo "✅ All services stopped."
echo "💾 Data volumes preserved (for development)"
echo ""
echo "To also remove data volumes, run:"
echo "  docker-compose down -v"
