# Boilerplate for GO app

Шаблонный проект на Go для быстрого старта разработки RESTful API приложений с использованием PostgreSQL.

## Описание

Boilerplate предоставляет готовую архитектуру для построения веб-сервисов с аутентификацией, работой с базой данных и полным набором инструментов для разработки и тестирования.

## Структура проекта

```
boilerplate-go/
├── cmd/
│   └── boilerplate/         # Точка входа приложения
│       └── main.go          # Инициализация, парсинг конфигурации, запуск
├── internal/
│   ├── api/                 # HTTP слой
│   │   ├── handlers/        # REST API handlers (auth, users)
│   │   ├── middleware/      # Middleware (auth, logging, recovery)
│   │   └── swagger/         # Swagger документация
│   ├── app/
│   │   └── boilerplate/     # Инициализация приложения
│   ├── model/               # Модели данных и конфигурация
│   │   └── config.go        # Структура конфигурации
│   ├── repository/          # Слой работы с БД
│   │   ├── repo.go          # Интерфейс репозитория
│   │   └── users.go         # Репозиторий пользователей
│   ├── services/            # Бизнес-логика
│   │   ├── auth/            # Сервис аутентификации
│   │   └── users/           # Сервис пользователей
│   ├── service_provider/    # Провайдер сервисов
│   └── pkg/                 # Внутренние пакеты
│       ├── clients/         # Клиенты (DB)
│       ├── closer/          # Graceful shutdown
│       ├── easyscan/        # Упрощенная работа с SQL
│       ├── errors/          # Типизированные ошибки
│       ├── gin/             # Gin utilities
│       ├── http_server/     # HTTP сервер
│       ├── jwt/             # JWT токены
│       ├── logger/          # Логирование (zap)
│       ├── metadata/        # Работа с метаданными
│       ├── pwd/             # Работа с паролями
│       ├── suite/           # Тестовые утилиты
│       ├── utils/           # Общие утилиты
│       └── version/         # Информация о версии
├── migrations/              # SQL миграции (goose)
│   ├── 001_create_table_users.sql
│   └── migrations.go        # Запуск миграций
├── bin/                     # Собранные бинарники
├── docker-compose.yaml      # PostgreSQL для разработки
├── Dockerfile               # Сборка образа приложения
├── Makefile                 # Команды для сборки/запуска/тестов
└── go.mod                   # Зависимости

```

## Основные компоненты

### 1. HTTP API (internal/api)
- **Handlers**: Обработчики HTTP запросов для аутентификации и управления пользователями
- **Middleware**: Аутентификация JWT, логирование, обработка паник
- **Swagger**: Автогенерируемая документация API

### 2. Сервисы (internal/services)
- **Auth Service**: 
  - Логин (выдача access/refresh токенов)
  - Обновление токенов (refresh)
  - Валидация токенов
  - Получение информации о текущем пользователе
- **Users Service**: 
  - CRUD операции над пользователями
  - Поиск пользователей с фильтрацией

### 3. Репозитории (internal/repository)
- Слой доступа к данным
- Использование Squirrel для построения SQL запросов
- Поддержка транзакций
- Типобезопасное сканирование результатов (easyscan)

### 4. База данных
- **PostgreSQL** 16
- Миграции через **Goose**
- Connection pooling (pgx/v5)
- Поддержка `pg_stat_statements` для мониторинга

### 5. Аутентификация и безопасность
- JWT токены (access и refresh)
- Хеширование паролей (bcrypt)
- Middleware для проверки авторизации
- Разделение публичных и приватных эндпоинтов

### 6. Инфраструктура
- **Logger**: Структурированное логирование (zap)
- **Closer**: Graceful shutdown с корректным закрытием ресурсов
- **HTTP Server**: Gin framework с настроенными middleware
- **Version**: Встраивание версии и даты сборки в бинарник

### 7. Тестирование
- Тестовая база данных (автоматическое создание/удаление)
- Test suite с фабриками и провайдерами
- Моки через mockery
- Линтер golangci-lint с расширенной конфигурацией

## Конфигурация

Приложение настраивается через:
- **Флаги командной строки**
- **Переменные окружения** (префикс `BOILERPLATE_`)
- **YAML конфиг-файл** (опционально)

Основные параметры:
```yaml
log-level: debug
db:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: boilerplate
  ssl-mode: false
api:
  host: 0.0.0.0
  port: 8080
  access-private-key: "your-secret-key"
  access-token-ttl: 3600      # 1 час
  refresh-token-ttl: 604800   # 7 дней
```

## Запуск

### Разработка
```bash
# Запуск PostgreSQL
make up

# Запуск приложения (с автоматическим созданием БД и миграциями)
make run
```

### Сборка
```bash
# Сборка бинарника с версией
make build VERSION=1.0.0

# Запуск бинарника
./bin/boilerplate --log-level=info --db-host=localhost ...
```

### Тесты
```bash
# Запуск тестов (создает тестовую БД автоматически)
make test

# Линтинг
make lint
```

### Docker
```bash
# Сборка образа
docker build -t boilerplate:latest .
```

## API Endpoints

### Public
- `POST /api/v1/auth/login` - Аутентификация
- `POST /api/v1/auth/refresh` - Обновление токена
- `POST /api/v1/auth/validate` - Валидация токена

### Protected (требуют JWT)
- `GET /api/v1/auth/me` - Информация о текущем пользователе
- `POST /api/v1/users` - Создание пользователя
- `GET /api/v1/users/:id` - Получение пользователя
- `PUT /api/v1/users/:id` - Обновление пользователя
- `DELETE /api/v1/users/:id` - Удаление пользователя
- `GET /api/v1/users` - Поиск пользователей

## Технологический стек

- **Язык**: Go 1.25
- **Web Framework**: Gin
- **БД**: PostgreSQL 16 + pgx/v5
- **Миграции**: Goose
- **SQL Builder**: Squirrel
- **Логирование**: Zap
- **JWT**: golang-jwt/jwt
- **Валидация**: go-playground/validator
- **CLI**: Cobra + Viper
- **Swagger**: swaggo/swag
- **Тестирование**: testify, mockery
- **Линтинг**: golangci-lint

## Особенности

✅ Clean Architecture (handlers → services → repository → db)  
✅ Dependency Injection через service provider  
✅ Graceful shutdown всех ресурсов  
✅ Структурированное логирование  
✅ Автоматические миграции при старте  
✅ JWT аутентификация с refresh токенами  
✅ Swagger документация  
✅ Docker-ready  
✅ Makefile для автоматизации  
✅ Тестовое окружение  
✅ Версионирование бинарников  

