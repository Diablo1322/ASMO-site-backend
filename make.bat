@echo off
echo Running Go tests...
cd backend
go test -v ./tests/unit/...
if %errorlevel% neq 0 exit /b %errorlevel%

echo Running integration tests...
docker-compose -f ../docker-compose.test.yml up -d
timeout /t 5 /nobreak > nul
set TEST_DATABASE_URL=postgres://test:test@localhost:5433/testdb?sslmode=disable
go test -v ./tests/integration/...
set TEST_DATABASE_URL=
docker-compose -f ../docker-compose.test.yml down

echo All tests completed!