@echo off
echo ========================================
echo    ASMO Backend - Development Mode
echo ========================================

echo Copying development environment...
copy .env.dev .env >nul

echo Starting development services...
docker-compose -f docker-compose.dev.yml down
docker-compose -f docker-compose.dev.yml up --build

if %errorlevel% neq 0 (
    echo.
    echo ERROR: Failed to start development services
    echo.
    pause
    exit /b %errorlevel%
)

echo.
echo ========================================
echo    Development services started!
echo ========================================
echo API: http://localhost/api/health
echo Direct backend: http://localhost:3000/api/health
echo PostgreSQL: localhost:5432
echo.
echo Mode: Development (HTTP)
echo.
pause