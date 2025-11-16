@echo off
echo ========================================
echo    Debugging Test Database
echo ========================================

echo Starting test database...
docker-compose -f docker-compose.test.yml up -d

echo Waiting for database to start...
timeout /t 5 /nobreak > nul

echo Checking test database connection...
docker exec asmo-site-backend-test-database-1 psql -U test -d testdb -c "SELECT version();"

if %errorlevel% neq 0 (
    echo.
    echo ❌ Cannot connect to test database!
    echo.
    echo Checking what databases exist...
    docker exec asmo-site-backend-test-database-1 psql -U test -d postgres -c "\l"
) else (
    echo.
    echo ✅ Test database is accessible!
    echo.
    echo Checking tables...
    docker exec asmo-site-backend-test-database-1 psql -U test -d testdb -c "\dt"
)

echo.
echo Test database logs:
docker-compose -f docker-compose.test.yml logs test-database

echo.
pause