@echo off
echo Building and starting Docker containers...
docker-compose up --build

if %errorlevel% neq 0 (
    echo Error starting containers. Check if Docker is running.
    pause
    exit /b %errorlevel%
)

echo Containers started successfully!
echo Frontend: http://localhost
echo Backend API: http://localhost/api/health
echo PostgreSQL: localhost:5432
pause