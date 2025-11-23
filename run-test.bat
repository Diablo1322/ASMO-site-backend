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
echo 🧪 Building test database...
docker-compose -f docker-compose.test.yml build

echo.
echo 🚀 Starting test database...
docker-compose -f docker-compose.test.yml up -d

echo.
echo ⏳ Waiting for test database to be ready...
ping -n 10 127.0.0.1 >nul

echo.
echo 🔍 Debugging test setup...
call debug-test-db.bat

echo.
echo ╔══════════════════════════════════════╗
echo ║            🧪 RUNNING TESTS          ║
echo ╚══════════════════════════════════════╝
echo.

cd backend

echo.
echo "=== 🔬 UNIT TESTS ==="
go test -v -short ./tests/unit/...

if %errorlevel% neq 0 (
    echo.
    echo ❌ Unit tests failed!
    goto cleanup
)

echo.
echo "=== 🔍 INTEGRATION TESTS ==="
go test -v ./tests/integration/...

if %errorlevel% neq 0 (
    echo.
    echo ❌ Integration tests failed!
    goto cleanup
)

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
echo 🎯 Test execution completed!
echo.
echo 💡 Tip: Use 'make.bat test-unit' for quick unit tests
echo 💡 Tip: Use 'make.bat test-integration' for integration tests
echo.
pause