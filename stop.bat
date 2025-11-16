@echo off
echo ========================================
echo    ASMO Backend - Stopping Services
echo ========================================

echo Stopping Docker containers...
docker-compose down

echo Cleaning up unused Docker resources...
docker system prune -f

echo.
echo ========================================
echo    Services stopped successfully!
echo ========================================
echo.
pause