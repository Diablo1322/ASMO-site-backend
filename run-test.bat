@echo off
chcp 65001 >nul
echo.
echo ╔══════════════════════════════════════╗
echo ║        ASMO Backend - TESTS          ║
echo ╚══════════════════════════════════════╝
echo.

echo 🔄 Switching to test mode...

echo 🛑 Stopping any running services...
docker-compose -f docker-compose.dev.yml down 2>nul
docker-compose -f docker-compose.prod.yml down 2>nul

echo.
echo 🧪 Building test environment...
docker-compose -f docker-compose.test.yml build

echo.
echo 🚀 Starting test services (DB + Redis)...
docker-compose -f docker-compose.test.yml up -d

echo.
echo ⏳ Waiting for test services to be ready...
timeout /t 10 /nobreak >nul

echo.
echo 🔍 Checking test database connection...
docker-compose -f docker-compose.test.yml exec -T test-database psql -U test -d testdb -c "SELECT version();"

echo.
echo 🔍 Checking test Redis connection...
docker-compose -f docker-compose.test.yml exec -T test-redis redis-cli ping

echo.
echo ╔══════════════════════════════════════╗
echo ║            🧪 RUNNING TESTS          ║
echo ╚══════════════════════════════════════╝
echo.

cd backend

echo.
echo "=== 🔬 UNIT TESTS ==="
go test -v -short ./tests/unit/... -cover -coverprofile=../test-results/unit-coverage.out

if %errorlevel% neq 0 (
    echo.
    echo ❌ Unit tests failed!
    goto cleanup
)

echo.
echo "=== 🔍 INTEGRATION TESTS ==="
go test -v ./tests/integration/... -cover -coverprofile=../test-results/integration-coverage.out

if %errorlevel% neq 0 (
    echo.
    echo ❌ Integration tests failed!
    goto cleanup
)

echo.
echo "=== 📊 GENERATING COVERAGE REPORT ==="
go tool cover -html=../test-results/unit-coverage.out -o ../test-results/unit-coverage.html
go tool cover -html=../test-results/integration-coverage.out -o ../test-results/integration-coverage.html

echo 📈 Coverage reports generated:
echo    - test-results/unit-coverage.html
echo    - test-results/integration-coverage.html

echo.
echo ╔══════════════════════════════════════╗
echo ║           ✅ ALL TESTS PASSED!       ║
echo ╚══════════════════════════════════════╝

:cleanup
echo.
echo 🧹 Cleaning up test containers...
cd ..
docker-compose -f docker-compose.test.yml down

echo.
echo 📊 Test execution completed!
echo 📁 Results saved in: test-results/
echo.
pause