@echo off
echo ========================================
echo    Debugging Migrations for Tests
echo ========================================

echo Current directory:
cd

echo.
echo Checking migrations location...
dir migrations /b 2>nul
if %errorlevel% neq 0 (
    echo ❌ migrations not found in current directory
    echo Checking parent directories...
    dir ..\migrations /b 2>nul
    if %errorlevel% equ 0 (
        echo ✅ migrations found in parent directory
    ) else (
        echo ❌ migrations not found
    )
)

echo.
echo Testing migration path detection...
cd backend
go run -v ./tests/testutils/debug_migrations.go

echo.
pause