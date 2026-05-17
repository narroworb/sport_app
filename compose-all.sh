#!/usr/bin/env bash
set -euo pipefail

# Wrapper to start all service docker-compose files together using multiple -f flags
BASE_DIR="$(cd "$(dirname "$0")" && pwd)"
FILES=(
  "${BASE_DIR}/services/auth_service/docker-compose.yaml"
  "${BASE_DIR}/services/core_api/docker-compose.yaml"
  "${BASE_DIR}/services/data_collector/docker-compose.yaml"
  "${BASE_DIR}/services/analytics_service/docker-compose.yaml"
  "${BASE_DIR}/services/frontend/docker-compose.yaml"
)
 
# Ensure network exists
if ! docker network ls --format '{{.Name}}' | grep -wq sport_network; then
  echo "Creating docker network: sport_network"
  docker network create sport_network >/dev/null
fi

## Ensure expected external volumes exist
VOLUMES=(
  pgdata
  ch_data
  ch_logs
  zk_data
  zk_logs
  kafka_data
  elasticsearch_data
  redis_data
  postgres_data
  clickhouse_data
)
for v in "${VOLUMES[@]}"; do
  if ! docker volume ls --format '{{.Name}}' | grep -wq "${v}"; then
    echo "Creating docker volume: ${v}"
    docker volume create "${v}" >/dev/null
  fi
done

ARGS=()
for f in "${FILES[@]}"; do
  ARGS+=( -f "$f" )
done

# If no args provided, default to build and up detached
if [ "$#" -eq 0 ]; then
  echo "Running: docker compose ${ARGS[*]} up -d --build"
  docker compose "${ARGS[@]}" up -d --build
else
  echo "Running: docker compose ${ARGS[*]} $*"
  docker compose "${ARGS[@]}" "$@"
fi
