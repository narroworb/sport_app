@echo off
REM Minimal: sequentially run each service-level docker compose up -d without building or creating volumes/networks.
SETLOCAL
pushd "%~dp0" >nul 2>&1
for /D %%D in (services\*) do call :StartService "%%D"
popd >nul 2>&1

echo All requested services started (no build).
ENDLOCAL
exit /b 0

:StartService
if exist "%~1\docker-compose.yaml" (
  echo Starting (no-build): docker compose -f "%~1\docker-compose.yaml" up -d
  docker compose -f "%~1\docker-compose.yaml" up -d
)
exit /b 0
