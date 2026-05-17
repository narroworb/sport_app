@echo off
REM Wrapper to run multiple docker-compose files together on Windows (CMD)
SETLOCAL ENABLEDELAYEDEXPANSION
set "BASE=%~dp0"
set "FILES=-f %BASE%services\auth_service\docker-compose.yaml -f %BASE%services\core_api\docker-compose.yaml -f %BASE%services\data_collector\docker-compose.yaml -f %BASE%services\analytics_service\docker-compose.yaml -f %BASE%services\frontend\docker-compose.yaml"

REM Ensure docker network exists (create if missing)
docker network ls --format "{{.Name}}" | findstr /R /C:"^sport_network$" >nul 2>&1
if errorlevel 1 (
  echo Creating docker network: sport_network
  docker network create sport_network >nul
)
if "%*"=="" (
  echo Running: docker compose %FILES% up -d --build
  docker compose %FILES% up -d --build
) else (
  echo Running: docker compose %FILES% %*
  docker compose %FILES% %*
)
ENDLOCAL

REM Ensure expected external volumes exist
for %%V in (pgdata ch_data ch_logs zk_data zk_logs kafka_data elasticsearch_data redis_data postgres_data clickhouse_data) do (
  docker volume ls --format "{{.Name}}" | findstr /R /C:"^%%V$" >nul 2>&1
  if errorlevel 1 (
    echo Creating docker volume: %%V
    docker volume create %%V >nul
  )
)

ENDLOCAL
