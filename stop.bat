@echo off
REM Sports Analytics Platform - Stop and Cleanup

echo.
echo Stopping Sports Analytics Platform...
echo.

docker-compose down

echo.
echo All services stopped.
echo Data volumes preserved for development.
echo.
echo To also remove data volumes, run:
echo   docker-compose down -v
echo.
pause
