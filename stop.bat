@echo off
echo Stopping Docker containers...
docker-compose down

echo Cleaning up...
docker system prune -f

echo Project stopped successfully!
pause