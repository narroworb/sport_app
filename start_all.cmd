@echo off
echo ========================================
echo Starting services in defined order
echo ========================================

REM Определяем порядок запуска
set SERVICES=data_collector auth_service core_api analytics_service frontend

REM Задержка между сервисами (в секундах)
set DELAY=25

for %%s in (%SERVICES%) do (
    if exist "services\%%s" (
        echo.
        echo [%%s] Starting service...
        cd services\%%s
        
        REM Запускаем сервис
        docker-compose up -d
        
        if %errorlevel% equ 0 (
            echo [%%s] Started successfully
            echo Waiting %DELAY% seconds for service to initialize...
            timeout /t %DELAY% /nobreak >nul
        ) else (
            echo [%%s] FAILED to start!
            pause
            exit /b 1
        )
        
        cd ..\..
    ) else (
        echo [%%s] WARNING - Folder "services\%%s" not found, skipping...
    )
)

echo.
echo ========================================
echo All services started successfully!
echo ========================================
pause