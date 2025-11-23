@echo off
chcp 65001 >nul
echo.
echo ╔══════════════════════════════════════╗
echo ║        ASMO BACKEND CONTROL          ║
echo ╚══════════════════════════════════════╝
echo.

:menu
echo 📋 Select Mode:
echo.
echo   1 🚀 DEVELOPMENT (Next.js on :3001)
echo   2 ✅ PRODUCTION (HTTPS)
echo   3 🧪 Run Tests
echo   4 📊 Service Status
echo   5 🛑 Stop All
echo   6 ❌ Exit
echo.

set /p choice=Choose (1-6):

if "%choice%"=="1" (
    call switch-to-dev.bat
    goto menu
)
if "%choice%"=="2" (
    call switch-to-prod.bat
    goto menu
)
if "%choice%"=="3" (
    call run-test.bat
    goto menu
)
if "%choice%"=="4" (
    echo.
    echo 📊 Running Services:
    docker-compose -f docker-compose.dev.yml ps
    docker-compose -f docker-compose.prod.yml ps
    echo.
    pause
    goto menu
)
if "%choice%"=="5" (
    echo.
    echo 🛑 Stopping all services...
    docker-compose -f docker-compose.dev.yml down 2>nul
    docker-compose -f docker-compose.prod.yml down 2>nul
    echo ✅ All services stopped!
    ping -n 2 127.0.0.1 >nul
    goto menu
)
if "%choice%"=="6" (
    echo 👋 Goodbye!
    ping -n 1 127.0.0.1 >nul
    exit /b 0
)

echo ❌ Invalid choice. Please try again.
goto menu