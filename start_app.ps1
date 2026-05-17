# start_app_simple.ps1
$ServicesPath = ".\services"

Write-Host "Поиск docker-compose.yaml в папке: $ServicesPath" -ForegroundColor Cyan

# Находим все docker-compose.yaml файлы
$composeFiles = Get-ChildItem -Path $ServicesPath -Recurse -Include "docker-compose.yaml", "docker-compose.yml" -File

if ($composeFiles.Count -eq 0) {
    Write-Host "Файлы docker-compose.yaml не найдены!" -ForegroundColor Red
    exit 1
}

Write-Host "Найдено $($composeFiles.Count) сервисов" -ForegroundColor Green

foreach ($composeFile in $composeFiles) {
    $serviceDir = $composeFile.DirectoryName
    $serviceName = Split-Path $serviceDir -Leaf
    
    Write-Host "`n>>> Запуск сервиса: $serviceName" -ForegroundColor Yellow
    
    # Переходим в папку сервиса
    Set-Location $serviceDir
    
    # Запускаем контейнеры
    docker-compose up -d
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "    ✓ Сервис $serviceName запущен" -ForegroundColor Green
    }
    else {
        Write-Host "    ✗ Ошибка при запуске $serviceName" -ForegroundColor Red
    }
}

# Возвращаемся в исходную папку
Set-Location $PSScriptRoot

Write-Host "`nГотово!" -ForegroundColor Cyan