@echo off
cd backend

if "%1"=="" (
    echo.
    echo ========================================
    echo    ASMO Backend - Make Utility
    echo ========================================
    echo.
    echo Available commands:
    echo   make.bat test        - Run all tests
    echo   make.bat test-unit   - Run unit tests only
    echo   make.bat test-integration - Run integration tests
    echo   make.bat migrate-up  - Run migrations
    echo   make.bat migrate-down - Rollback migrations
    echo   make.bat build       - Build application
    echo   make.bat run         - Run application locally
    echo.
    exit /b 0
)

if "%1"=="test" (
    echo Running all tests...
    docker-compose -f ../docker-compose.test.yml up -d
    timeout /t 5 /nobreak > nul
    go test -v ./tests/...
    docker-compose -f ../docker-compose.test.yml down
) else if "%1"=="test-unit" (
    echo Running unit tests...
    go test -v ./tests/unit/...
) else if "%1"=="test-integration" (
    echo Running integration tests...
    docker-compose -f ../docker-compose.test.yml up -d
    timeout /t 5 /nobreak > nul
    go test -v ./tests/integration/...
    docker-compose -f ../docker-compose.test.yml down
) else if "%1"=="migrate-up" (
    echo Running migrations...
    go run cmd/migrate/main.go up
) else if "%1"=="migrate-down" (
    echo Rolling back migrations...
    go run cmd/migrate/main.go down
) else if "%1"=="build" (
    echo Building application...
    go build -o main ./cmd/server/
    go build -o migrate ./cmd/migrate/
    echo Build completed!
) else if "%1"=="run" (
    echo Running application locally...
    go run cmd/server/main.go
) else (
    echo Unknown command: %1
    echo Run make.bat without arguments to see available commands.
    exit /b 1
)