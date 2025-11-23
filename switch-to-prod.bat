@echo off
chcp 65001 >nul
echo.
echo ╔══════════════════════════════════════╗
echo ║    ASMO Backend - PRODUCTION Mode    ║
echo ╚══════════════════════════════════════╝
echo.

echo 🔒 Checking production requirements...

if not exist "ssl\asmo-backend.crt" (
    echo 🔐 Generating SSL certificates...
    call create-ssl-certs.bat
)

echo 🚀 Setting up production environment...
copy .env.production .env >nul 2>&1

echo ⚠️  Please edit .env file for production values!
ping -n 2 127.0.0.1 >nul

echo 🛑 Stopping any running services...
docker-compose -f docker-compose.dev.yml down 2>nul
docker-compose -f docker-compose.prod.yml down 2>nul

echo 🚀 Starting production stack...
docker-compose -f docker-compose.prod.yml up --build -d

echo.
echo ⏳ Waiting for services to start...
ping -n 8 127.0.0.1 >nul

echo.
echo ✅ PRODUCTION Mode Activated!
echo.
echo 📍 Endpoints:
echo    HTTPS API: https://localhost/api/health
echo    HTTP Redirect: http://localhost
echo.
echo 🔒 Features:
echo    ✅ HTTPS enforced
echo    ✅ Security headers
echo    ✅ Production logging
echo    ✅ CORS for production domain
echo.
echo 💡 Remember: Update CORS in server/main.go for your production domain!
echo.
pause