#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"

# Ensure network (use inspect for reliability)
if ! docker network inspect sport_network >/dev/null 2>&1; then
  echo "Creating docker network: sport_network"
  docker network create sport_network >/dev/null
fi

# Ensure volumes exist (do NOT create automatically). Use inspect for reliability.
VOLUMES=(pgdata ch_data ch_logs zk_data zk_logs kafka_data elasticsearch_data redis_data postgres_data clickhouse_data)
MISSING=()
for v in "${VOLUMES[@]}"; do
  if ! docker volume inspect "${v}" >/dev/null 2>&1; then
    MISSING+=("${v}")
  fi
done

if [ ${#MISSING[@]} -ne 0 ]; then
  echo "Error: the following required docker volumes are missing: ${MISSING[*]}"
  echo "Create them manually (or restore from backup) before running this script. Example:"
  echo "  docker volume create <name>"
  exit 1
fi

# Find service-level docker-compose files and bring them up sequentially
for compose_file in "${BASE_DIR}"/services/*/docker-compose.yaml; do
  if [ -f "${compose_file}" ]; then
    echo "Starting services from: ${compose_file}"
    docker compose -f "${compose_file}" up -d --build
  fi
done

echo "All services started."
