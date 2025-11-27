@echo off
echo ========================================
echo    Debugging Test Database & Redis
echo ========================================

echo Starting test services...
docker-compose -f docker-compose.test.yml up -d

echo Waiting for services to start...
timeout /t 5 /nobreak >nul

echo.
echo 🔍 Checking test database connection...
docker-compose -f docker-compose.test.yml exec -T test-database psql -U test -d testdb -c "SELECT version();"

if %errorlevel% neq 0 (
    echo.
    echo ❌ Cannot connect to test database!
    echo.
    echo Checking what databases exist...
    docker-compose -f docker-compose.test.yml exec -T test-database psql -U test -d postgres -c "\l"
) else (
    echo.
    echo ✅ Test database is accessible!
    echo.
    echo Checking tables...
    docker-compose -f docker-compose.test.yml exec -T test-database psql -U test -d testdb -c "\dt"
)

echo.
echo 🔍 Checking test Redis connection...
docker-compose -f docker-compose.test.yml exec -T test-redis redis-cli ping

if %errorlevel% neq 0 (
    echo ❌ Cannot connect to test Redis!
) else (
    echo ✅ Test Redis is accessible!
    echo.
    echo Testing Redis operations...
    docker-compose -f docker-compose.test.yml exec -T test-redis redis-cli set test_key "hello_world"
    docker-compose -f docker-compose.test.yml exec -T test-redis redis-cli get test_key
)

echo.
echo 📊 Test services logs:
docker-compose -f docker-compose.test.yml logs --tail=20

echo.
pause