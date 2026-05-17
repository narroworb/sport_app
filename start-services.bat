@echo off
REM Sequentially start service-level docker-compose files from services/*
SETLOCAL ENABLEDELAYEDEXPANSION
set "ROOT=%~dp0"

REM Ensure network (use inspect for reliability)
docker network inspect sport_network >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
  echo Creating docker network: sport_network
  docker network create sport_network >nul 2>&1
)

REM Ensure volumes exist (do NOT create automatically). Use inspect for reliability.
set "MISSING="
for %%V in (pgdata ch_data ch_logs zk_data zk_logs kafka_data elasticsearch_data redis_data postgres_data clickhouse_data) do (
  docker volume inspect %%V >nul 2>&1
  if NOT !ERRORLEVEL! EQU 0 (
    set "MISSING=!MISSING! %%V"
  )
)

if defined MISSING (
  echo Error: the following required docker volumes are missing:%MISSING%
  echo Create them manually before running this script, e.g. "docker volume create ^<name^>"
  exit /b 1
)

REM Iterate service dirs
for /D %%D in ("%ROOT%services\*") do (
  if exist "%%D\docker-compose.yaml" (
    echo Starting services from: %%D\docker-compose.yaml
    docker compose -f "%%D\docker-compose.yaml" up -d --build
  )
)

echo All services started.
ENDLOCAL
