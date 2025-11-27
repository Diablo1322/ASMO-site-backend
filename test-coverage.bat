@echo off
chcp 65001 >nul
echo.
echo ╔══════════════════════════════════════╗
echo ║         TEST COVERAGE REPORT         ║
echo ╚══════════════════════════════════════╝
echo.

if not exist "test-results" (
    echo ❌ No test results found!
    echo Run tests first: run-test.bat
    pause
    exit /b 1
)

echo 📊 Generating detailed coverage reports...

cd backend

echo.
echo "=== OVERALL COVERAGE ==="
go test -coverprofile=../test-results/total-coverage.out ./...
go tool cover -func=../test-results/total-coverage.out

echo.
echo "=== PACKAGE COVERAGE ==="
for /f "tokens=1" %%p in ('go list ./...') do (
    echo 📦 %%p
    go test -coverprofile=../test-results/%%~np-coverage.out %%p
    go tool cover -func=../test-results/%%~np-coverage.out | findstr "total:"
)

echo.
echo "=== HTML REPORTS ==="
go tool cover -html=../test-results/total-coverage.out -o ../test-results/total-coverage.html
echo ✅ Total coverage: test-results/total-coverage.html

echo.
echo 📈 Coverage reports generated in test-results/ folder:
dir test-results\*.html /b

echo.
echo 🚀 Opening coverage report...
start test-results\total-coverage.html

pause