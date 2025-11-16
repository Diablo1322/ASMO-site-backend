@echo off
echo ========================================
echo    ASMO Backend - Production Mode
echo ========================================

echo Checking for SSL certificates...
if not exist "ssl\asmo-backend.crt" (
    echo SSL certificates not found. Creating self-signed certificates...
    call create-ssl-certs.bat
)

echo Copying production environment...
copy .env.prod .env >nul

echo Starting production services...
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up --build -d

if %errorlevel% neq 0 (
    echo.
    echo ERROR: Failed to start production services
    echo.
    pause
    exit /b %errorlevel%
)

echo.
echo ========================================
echo    Production services started!
echo ========================================
echo HTTP redirect: http://localhost
echo HTTPS API: https://localhost/api/health
echo PostgreSQL: localhost:5432
echo.
echo Mode: Production (HTTPS)
echo Containers running in background.
echo.
echo To view logs: docker-compose -f docker-compose.prod.yml logs -f
echo To stop: docker-compose -f docker-compose.prod.yml down
echo.
pause