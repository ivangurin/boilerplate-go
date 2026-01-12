# Boilerplate for GO app

Шаблонный проект на Go для быстрого старта разработки микросервисных приложений с поддержкой gRPC и HTTP Gateway или только HTTP API.

## Структура проекта

```text
boilerplate-backend/
├── cmd/
│   └── boilerplate/                      # Точка входа в приложение
├── internal/
│   ├── api/
│   │   ├── grpc/                      # gRPC API
│   │   │   ├── handlers/              # gRPC handlers
│   │   │   │   ├── auth/              # Auth gRPC handlers
│   │   │   │   ├── users/             # Users gRPC handlers
│   │   │   │   └── handlers.go        # gRPC handlers registry
│   │   │   └── middleware/            # gRPC middleware
│   │   │       ├── auth.go            # JWT аутентификация
│   │   │       ├── logger.go          # Логирование запросов
│   │   │       ├── middleware.go      # Middleware interface
│   │   │       ├── panic.go           # Recover от паник
│   │   │       ├── tracer.go          # Request ID и IP трейсинг
│   │   │       └── validate.go        # Валидация запросов
│   │   └── http/                      # HTTP REST API
│   │       ├── handlers/              # HTTP handlers
│   │       │   ├── auth/              # Приватные auth endpoints
│   │       │   ├── auth_public/       # Публичные auth endpoints
│   │       │   ├── users/             # Приватные users endpoints
│   │       │   ├── users_public/      # Публичные users endpoints
│   │       │   └── handlers.go        # HTTP handlers registry
│   │       ├── middleware/            # HTTP middleware
│   │       │   ├── auth.go            # JWT аутентификация
│   │       │   ├── logger.go          # Логирование запросов
│   │       │   ├── middleware.go      # Middleware interface
│   │       │   └── tracer.go          # Request ID и IP трейсинг
│   │       └── swagger/               # Swagger документация
│   ├── app/
│   │   └── boilerplate/                  # Инициализация и запуск приложения
│   ├── migrations/                    # Миграции базы данных
│   ├── model/                         # Модели данных и конфигурация
│   ├── pkg/                           # Переиспользуемые пакеты
│   │   ├── clients/db/                # Клиент базы данных (rel)
│   │   ├── closer/                    # Graceful shutdown
│   │   ├── errors/                    # Кастомные ошибки (400, 401, 403, 404)
│   │   ├── gin/                       # Gin утилиты
│   │   ├── grpc_server/               # gRPC сервер
│   │   ├── http_server/               # HTTP сервер
│   │   ├── jwt/                       # JWT токены
│   │   ├── logger/                    # Логирование (zap)
│   │   ├── metadata/                  # Context metadata (request_id, ip, user_id)
│   │   ├── pwd/                       # Password hashing (bcrypt)
│   │   ├── swagger/                   # Swagger 
│   │   ├── suite/                     # Тестовые утилиты
│   │   │   ├── factory/               # Фабрики тестовых данных
│   │   │   └── provider/              # Провайдеры для тестов
│   │   ├── utils/                     # Общие утилиты
│   │   └── version/                   # Версионирование
│   ├── repository/                    # Слой работы с БД
│   │   └── model/                     # DB модели
│   ├── service_provider/              # Провайдер для сервисов
│   └── services/                      # Бизнес-логика
│       ├── auth/                      # Сервис аутентификации
│       ├── log_events/                # Сервис логирования событий
│       └── users/                     # Сервис управления пользователями
│           └── mocks/                 # Mock для тестов
├── pkg/
│   └── pb/                            # Сгенерированные protobuf файлы
│       ├── auth.*                     # Auth (pb.go, pb.gw.go, pb.validate.go, grpc.pb.go)
│       ├── users.*                    # Users (pb.go, pb.gw.go, pb.validate.go, grpc.pb.go)
├── proto/                             # Protobuf определения
│   ├── auth.proto                     # Auth API схема
│   ├── users.proto                    # Users API схема
├── bin/                               # Скомпилированные бинарники
├── .github/                           # GitHub Actions
├── .gitlab-ci.yml                     # GitLab CI/CD
├── .golangci.yaml                     # Golangci-lint конфигурация
├── .mockery.yaml                      # Mockery конфигурация
├── buf.yaml                           # Buf конфигурация модуля
├── buf.gen.yaml                       # Buf генерация настройки
├── docker-compose.yaml                # Docker окружение
├── Dockerfile                         # Dockerfile для сборки образа
├── go.mod                             # Go модули и зависимости
├── Makefile                           # Команды для сборки и тестирования
└── README.md                          # Документация проекта
```

## Основные компоненты

### cmd/boilerplate

Точка входа приложения. Обрабатывает флаги командной строки, переменные окружения и конфигурационные файлы с помощью Cobra и Viper.

### proto/

Protobuf определения для gRPC и REST API:

- **auth.proto** - API аутентификации (login, logout, refresh, me)
- **users.proto** - API управления пользователями (CRUD)
- Автоматическая генерация:
  - Go структур (protoc-gen-go)
  - gRPC серверов/клиентов (protoc-gen-go-grpc)
  - REST Gateway (protoc-gen-grpc-gateway)
  - OpenAPI спецификации (protoc-gen-openapiv2)
  - Валидации (protoc-gen-validate)

### internal/api

REST и gRPC API:

- **grpc/** - gRPC сервер и middleware
- **middleware/** - gRPC interceptors (auth, logger, tracer с извлечением IP)
- **handlers/** - обработчики HTTP запросов
- **middleware/** - промежуточные обработчики HTTP (аутентификация, логирование, трейсинг)
- **swagger/** - автогенерируемая документация API

### internal/services

Бизнес-логика приложения:

- **auth** - аутентификация (login, logout, refresh tokens)
- **users** - управление пользователями (CRUD операции)
- **log_events** - логирование событий системы

### internal/repository

Слой доступа к данным, работа с PostgreSQL.

### internal/pkg

Переиспользуемые утилиты и библиотеки проекта:

- **metadata** - управление контекстом (request_id, ip, user_id, user_name)
  - `SetRequestID/GetRequestID` - уникальный идентификатор запроса
  - `SetIP/GetIP` - IP адрес клиента (извлекается из peer в gRPC)
  - `SetUserID/GetUserID` - ID аутентифицированного пользователя
  - `SetUserName/GetUserName` - имя пользователя
- **logger** - структурированное логирование
- **jwt** - работа с JWT токенами
- **errors** - типизированные ошибки приложения
- **utils** - общие утилиты (UUID, группировки, лимиты и т.д.)

## Docker сервисы

Проект использует следующие сервисы в Docker окружении:

### PostgreSQL

- **Версия**: 16
- **Порт**: 5432
- **Credentials**: postgres/postgres
- **Расширения**: `pg_stat_statements` для сбора статистики запросов
- **Volume**: `pg-data` для персистентного хранения данных

База данных для хранения всех данных приложения (пользователи, события логирования и т.д.).

### Сеть

Все сервисы работают в единой Docker сети `boilerplate-network`, что позволяет им взаимодействовать друг с другом по именам сервисов.

## Конфигурация

Приложение поддерживает несколько источников конфигурации с приоритетом:

1. Флаги командной строки
2. Переменные окружения (с префиксом `BOILERPLATE_`)
3. Конфигурационный файл (`config.yaml`)
4. Значения по умолчанию

Пример переменных окружения:

```bash
# Уровень логирования (debug, info, warn, error)
BOILERPLATE_LOG_LEVEL=info

# Настройки API
BOILERPLATE_API_HOST=127.0.0.1
BOILERPLATE_API_HTTP_PORT=8080
BOILERPLATE_API_GRPC_PORT=8082
BOILERPLATE_API_ACCESS_PRIVATE_KEY=your-secret-key
BOILERPLATE_API_ACCESS_TOKEN_TTL=3600        # в секундах (1 час)
BOILERPLATE_API_REFRESH_TOKEN_TTL=604800     # в секундах (7 дней)

# Настройки базы данных
BOILERPLATE_DB_HOST=localhost
BOILERPLATE_DB_PORT=5432
BOILERPLATE_DB_USER=postgres
BOILERPLATE_DB_PASSWORD=postgres
BOILERPLATE_DB_NAME=boilerplate
BOILERPLATE_DB_SSL_MODE=false
```

## Сборка и запуск

```bash
# Запуск Docker окружения
make up

# Запуск приложения (создаст БД если нужно)
make run

# Сборка бинарника
make build

# Запуск тестов
make test

# Генерация кода (mocks, swagger, protobuf)
make generate
```

## Работа с Protobuf

Проект использует [Buf](https://buf.build/) для управления protobuf файлами.

### Установленные плагины

- **protoc-gen-go** - генерация Go структур
- **protoc-gen-go-grpc** - генерация gRPC серверов и клиентов
- **protoc-gen-grpc-gateway** - генерация REST Gateway
- **protoc-gen-openapiv2** - генерация OpenAPI/Swagger спецификации
- **protoc-gen-validate** - генерация валидации полей

### Структура proto файлов

```
proto/
├── auth.proto              # Auth API
├── users.proto             # Users API
├── buf.yaml                # Buf конфигурация модуля
└── buf.gen.yaml            # Настройки генерации
```

### Добавление нового API

1. Создайте `.proto` файл в `proto/`
2. Определите сервисы и сообщения
3. Добавьте HTTP аннотации для REST Gateway
4. Запустите `make gen-proto`
