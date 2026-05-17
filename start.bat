@echo off
REM Sports Analytics Platform Startup Script for Windows

echo.
echo Starting Sports Analytics Platform...
echo.

REM Check if Docker is running
docker info >nul 2>&1
if errorlevel 1 (
    echo Docker is not running. Please start Docker Desktop and try again.
    exit /b 1
)

echo Building images...
docker-compose build

echo.
echo Starting services...
docker-compose up -d

echo.
echo Waiting for services to be ready...
timeout /t 10 /nobreak

echo.
echo ===== Services Started Successfully! =====
echo.
echo Frontend:          http://localhost
echo Auth Service:      http://localhost:8081
echo Core API:          http://localhost:8080
echo Analytics:         http://localhost:8082
echo.
echo Databases:
echo   PostgreSQL:      localhost:5432
echo   ClickHouse:      localhost:8123
echo   Redis:           localhost:6379
echo   Elasticsearch:   localhost:9200
echo   Kafka:           localhost:9092
echo.
echo To view logs:      docker-compose logs -f
echo To stop:           docker-compose down
echo.
pause
