# ASMO Site Backend

Backend сервис для сайта ASMO с использованием Go, PostgreSQL, Docker и Nginx.

## Структура проекта

- `backend/` - Go backend сервер
- `nginx/` - Nginx reverse proxy
- `frontend/` - Next.js frontend (отдельный репозиторий)

## Быстрый старт

### Предварительные требования
- Docker Desktop для Windows
- Go 1.21+
- Git

### Запуск проекта

1. Клонируйте репозиторий
2. Запустите `run.bat` для сборки и запуска контейнеров
3. Откройте http://localhost для доступа к приложению

### Остановка проекта
Запустите `stop.bat` для остановки контейнеров.

## API Endpoints

- `GET /api/health` - Health check сервера
- `GET /api/items` - Получить список items
- `POST /api/items` - Создать новый item

## Разработка

### Backend разработка
```bash
cd backend
go run cmd/server/main.go