# ASMO Backend Service

Production-ready backend API для управления проектами веб-приложений, мобильных приложений, ботов и сотрудниками.

## 🏗️ Архитектура
ASMO-site-backend/
├── backend/ # Go backend сервер
│ ├── cmd/ # Точки входа
│ ├── internal/ # Внутренние пакеты
│ ├── migrations/ # Миграции БД
│ └── pkg/ # Внешние пакеты
├── nginx/ # Reverse proxy + SSL
├── tests/ # Unit и интеграционные тесты
└── docker-compose.*.yml # Окружения

text

## 🚀 Быстрый старт

### Предварительные требования
- Docker & Docker Compose
- Go 1.25.4+ (для разработки)

### Development режим
```bash
# Автоматическая настройка
switch-to-dev.bat

# Или вручную
docker-compose -f docker-compose.dev.yml up --build
Development endpoints:

🚀 API: http://localhost/api

🗄️ PGAdmin: http://localhost:5050 (admin@asmo.com/admin)

📊 База данных: localhost:5432

Production режим
bash
# С SSL сертификатами
switch-to-prod.bat

# Или вручную
docker-compose -f docker-compose.prod.yml up --build -d
📡 API Endpoints
Health Check
GET /api/health - Статус сервиса и БД

Web Applications
GET /api/WebApplications - Список веб-проектов

GET /api/WebApplications/:id - Проект по ID

POST /api/WebApplications - Создать проект

Mobile Applications
GET /api/MobileApplications - Список мобильных проектов

GET /api/MobileApplications/:id - Проект по ID

POST /api/MobileApplications - Создать проект

Bots
GET /api/Bots - Список бот-проектов

GET /api/Bots/:id - Проект по ID

POST /api/Bots - Создать проект

Staff
GET /api/Staff - Список сотрудников

GET /api/Staff/:id - Сотрудник по ID

POST /api/Staff - Добавить сотрудника

🗃️ Модели данных
WebProjects / MobileProjects / BotsProjects
json
{
  "id": 1,
  "name": "Название проекта (15-100 символов)",
  "description": "Описание (20-1500 символов)",
  "img": "https://example.com/image.jpg",
  "price": 1500.50,
  "time_develop": 30,
  "created_at": "2024-01-01T00:00:00Z",
  "update_at": "2024-01-01T00:00:00Z"
}
Staff
json
{
  "id": 1,
  "name": "ФИО сотрудника (15-100 символов)",
  "description": "Описание (20-1500 символов)",
  "img": "https://example.com/photo.jpg",
  "role": "Должность (1-50 символов)",
  "created_at": "2024-01-01T00:00:00Z",
  "update_at": "2024-01-01T00:00:00Z"
}
🧪 Тестирование
bash
# Все тесты
run-test.bat

# Только unit тесты
cd backend && go test ./tests/unit/...

# Только интеграционные тесты
cd backend && go test ./tests/integration/...
🔧 Утилиты
Миграции БД
bash
# Создать миграции
create-migrations.bat

# Применить миграции
docker-compose exec backend ./migrate up

# Откатить миграции
docker-compose exec backend ./migrate down
SSL сертификаты
bash
# Сгенерировать self-signed certificates
create-ssl-certs.bat
Дебаггинг
bash
# Проверить структуру проекта
check-structure.bat

# Дебаг миграций
debug-migrations.bat

# Дебаг тестовой БД
debug-test-db.bat
⚙️ Конфигурация
Environment Variables
Development (.env.dev):

env
DATABASE_URL=postgres://user:password@postgres:5432/asmo_db?sslmode=disable
PORT=3000
LOG_LEVEL=DEBUG
ENVIRONMENT=development
Production (.env.production):

env
DB_HOST=postgres
DB_PORT=5432
DB_USER=asmo_prod_user
DB_PASSWORD=secure_password
DB_NAME=asmo_production
DB_SSL_MODE=require
PORT=3000
LOG_LEVEL=INFO
ENVIRONMENT=production
ALLOWED_ORIGINS=https://your-domain.com
🔒 Безопасность
✅ HTTPS (Production)

✅ CORS настройки

✅ Rate limiting

✅ Security headers

✅ Валидация входных данных

✅ SQL injection protection

📊 Логирование
Структурированные JSON логи с уровнями:

DEBUG - Детальная отладка

INFO - Основные события

WARN - Предупреждения

ERROR - Ошибки

🐳 Docker команды
bash
# Просмотр логов
docker-compose logs -f backend

# Проверка здоровья
docker-compose ps

# Остановка сервисов
docker-compose down

# Пересборка
docker-compose build --no-cache
🚀 Деплой
GitHub
bash
deploy-to-github.bat
Production деплой
Настройте .env.production

Обновите ALLOWED_ORIGINS для фронтенда

Замените SSL сертификаты на реальные

Запустите: switch-to-prod.bat

📞 Поддержка
При проблемах проверьте:

Docker запущен и порты свободны

.env файлы настроены правильно

Миграции применены: docker-compose exec backend ./migrate up

Логи: docker-compose logs -f backend

🏆 Особенности
✅ Production-ready архитектура

✅ Автоматические миграции БД

✅ Полная тестовая покрытие

✅ HTTPS & Security headers

✅ Rate limiting & CORS

✅ Структурированное логирование

✅ Health checks

✅ Docker-optimized