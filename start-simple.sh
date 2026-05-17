#!/usr/bin/env bash
set -euo pipefail

BASE_DIR="$(cd "$(dirname "$0")" && pwd)"

# Minimal: run each service-level docker-compose up detached, do NOT build or create volumes/networks.
for compose_file in "${BASE_DIR}"/services/*/docker-compose.yaml; do
  if [ -f "${compose_file}" ]; then
    echo "Starting (no-build): docker compose -f ${compose_file} up -d"
    docker compose -f "${compose_file}" up -d
  fi
done

echo "All requested services started (no build)."
