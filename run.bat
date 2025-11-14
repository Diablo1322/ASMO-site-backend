@echo off
echo ========================================
echo    ASMO Backend - Starting Services
echo ========================================

echo Building and starting Docker containers...
docker-compose up --build

if %errorlevel% neq 0 (
    echo.
    echo ERROR: Failed to start containers.
    echo Please check if Docker is running.
    echo.
    pause
    exit /b %errorlevel%
)

echo.
echo ========================================
echo    Services started successfully!
echo ========================================
echo Frontend:    http://localhost
echo Backend API: http://localhost/api/health
echo PostgreSQL:  localhost:5432
echo.
echo Press any key to stop services...
pause >nul

call stop.bat