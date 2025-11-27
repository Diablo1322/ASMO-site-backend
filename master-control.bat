@echo off
chcp 65001 >nul
echo.
echo ╔══════════════════════════════════════╗
echo ║        ASMO BACKEND CONTROL          ║
echo ╚══════════════════════════════════════╝
echo.
echo 📍 Frontend: Separate repository
echo 📍 Backend:  Current project
echo.

:menu
echo 📋 Select Mode:
echo.
echo   1 🚀 DEVELOPMENT (Backend + DB)
echo   2 ✅ PRODUCTION (Backend + DB + Redis)
echo   3 🧪 Run Tests
echo   4 📊 Service Status
echo   5 🔧 Database Tools
echo   6 🛑 Stop All
echo   7 ❌ Exit
echo.

set /p choice=Choose (1-7):

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
    echo 🔧 Database Tools:
    echo   1. Backup database
    echo   2. View migrations
    echo   3. Check connections
    echo.
    set /p db_choice="Choose: "
    if "%db_choice%"=="1" (
        docker-compose exec backend ./migrate backup
    )
    goto menu
)
if "%choice%"=="6" (
    echo.
    echo 🛑 Stopping all services...
    docker-compose -f docker-compose.dev.yml down 2>nul
    docker-compose -f docker-compose.prod.yml down 2>nul
    echo ✅ All services stopped!
    timeout /t 2 >nul
    goto menu
)
if "%choice%"=="7" (
    echo 👋 Goodbye!
    timeout /t 1 >nul
    exit /b 0
)

echo ❌ Invalid choice. Please try again.
goto menu