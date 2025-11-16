@echo off
echo ========================================
echo    Checking Project Structure
echo ========================================

echo Checking for migrations folder...
if exist "backend\migrations" (
    echo ✅ migrations folder exists
    dir "backend\migrations\*.sql" /b
) else (
    echo ❌ migrations folder NOT found!
    echo Creating migrations folder...
    mkdir backend\migrations
)

echo.
echo Checking for migration files...
if exist "backend\migrations\*.sql" (
    echo ✅ SQL migration files exist
) else (
    echo ❌ No SQL migration files found!
    echo Please create migration files.
)

echo.
echo Project structure check completed.
pause