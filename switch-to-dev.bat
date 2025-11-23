@echo off
chcp 65001 >nul
echo.
echo ╔══════════════════════════════════════╗
echo ║    ASMO Backend - DEVELOPMENT Mode   ║
echo ╚══════════════════════════════════════╝
echo.

echo 🔧 Setting up development environment...
copy .env.dev .env >nul 2>&1

echo 🛑 Stopping any running services...
docker-compose -f docker-compose.prod.yml down 2>nul
docker-compose -f docker-compose.dev.yml down 2>nul

echo 🗑️  Cleaning up old containers and images...
docker system prune -f

echo 🚀 Starting development stack...
docker-compose -f docker-compose.dev.yml up --build

echo.
echo ⏳ Waiting for services to start...
ping -n 10 127.0.0.1 >nul

echo.
echo ✅ DEVELOPMENT Mode Activated!
echo.
echo 📍 Endpoints:
echo    Backend API: http://localhost:3000
echo    Frontend:    http://localhost:3001 (Next.js)
echo    PGAdmin:     http://localhost:5050
echo    Nginx:       http://localhost
echo.
echo 🔧 Features:
echo    ✅ Hot reload enabled
echo    ✅ Debug logging
echo    ✅ CORS for Next.js on :3001
echo    ✅ Database management UI
echo.
echo 🛑 To stop: Ctrl+C or run stop-dev.bat
echo.
pause