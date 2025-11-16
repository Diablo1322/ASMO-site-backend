@echo off
echo ========================================
echo    ASMO Backend - Running Tests
echo ========================================

echo Building test database...
docker-compose -f docker-compose.test.yml build

echo Starting test database...
docker-compose -f docker-compose.test.yml up -d

echo Waiting for test database to be ready...
timeout /t 10 /nobreak > nul

echo Running tests...
cd backend

echo.
echo ===== UNIT TESTS =====
go test -v ./tests/unit/...

if %errorlevel% neq 0 (
    echo.
    echo ❌ Unit tests failed!
    goto :cleanup
)

echo.
echo ===== INTEGRATION TESTS =====
go test -v ./tests/integration/...

if %errorlevel% neq 0 (
    echo.
    echo ❌ Integration tests failed!
    goto :cleanup
)

echo.
echo ========================================
echo    ✅ ALL TESTS PASSED!
echo ========================================

:cleanup
echo.
echo Cleaning up test containers...
cd ..
docker-compose -f docker-compose.test.yml down

echo Test execution completed.
pause