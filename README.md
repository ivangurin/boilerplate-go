# Go Microservice Boilerplate

A comprehensive boilerplate for building production-ready microservices in Go with support for gRPC and HTTP REST APIs, message queuing, file storage, email notifications, and PDF generation.

## Features

### Core Architecture
- **Dual API Support**: gRPC and HTTP/REST with automatic gateway
- **Protocol Buffers**: Type-safe API definitions with code generation
- **JWT Authentication**: Secure token-based authentication with access and refresh tokens
- **Modular Design**: Clean architecture with clear separation of concerns
- **Dependency Injection**: Service provider pattern for managing dependencies

### Infrastructure & Integration
- **Database**: PostgreSQL 16 with migrations and connection pooling
- **Message Broker**: NATS JetStream for asynchronous messaging
- **Object Storage**: MinIO (S3-compatible) for file storage
- **Email Service**: SMTP integration for email notifications
- **PDF Generation**: Headless Chrome for HTML to PDF conversion
- **Structured Logging**: Zap-based logging with context propagation

### Developer Experience
- **Docker Compose**: Complete local development environment
- **Code Generation**: Automated mock generation with Mockery
- **API Documentation**: Auto-generated Swagger/OpenAPI specs
- **Testing**: Comprehensive test suite with coverage reports
- **Linting**: golangci-lint for code quality
- **Makefile**: Convenient commands for all common tasks

## Project Structure

```
boilerplate-go/
├── cmd/
│   └── boilerplate/           # Application entry point
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── grpc/              # gRPC server implementation
│   │   │   ├── handlers/      # gRPC service handlers
│   │   │   └── middleware/    # gRPC interceptors
│   │   ├── http/              # HTTP REST server
│   │   │   ├── handlers/      # HTTP handlers
│   │   │   └── middleware/    # HTTP middleware
│   │   └── swagger/           # Generated API documentation
│   ├── app/                   # Application initialization
│   ├── model/                 # Domain models and interfaces
│   ├── pkg/
│   │   ├── clients/
│   │   │   ├── chrome/        # Headless Chrome client for PDF
│   │   │   ├── db/            # PostgreSQL client
│   │   │   ├── mail/          # Email client
│   │   │   ├── nats/          # NATS messaging client
│   │   │   └── s3/            # S3/MinIO storage client
│   │   ├── servers/
│   │   │   ├── grpc/          # gRPC server
│   │   │   ├── http/          # HTTP server (Gin)
│   │   │   └── nats/          # NATS consumer server
│   │   ├── closer/            # Graceful shutdown manager
│   │   ├── convert/           # Type conversion utilities
│   │   ├── errors/            # Custom error types
│   │   ├── gateway/           # gRPC-Gateway configuration
│   │   ├── jwt/               # JWT token management
│   │   ├── logger/            # Structured logging
│   │   ├── metadata/          # Context metadata handling
│   │   ├── pwd/               # Password hashing
│   │   ├── swagger/           # Swagger UI integration
│   │   ├── utils/             # Common utilities
│   │   └── version/           # Version information
│   ├── repository/            # Data access layer
│   ├── service_provider/      # Dependency injection
│   └── services/
│       ├── auth/              # Authentication service
│       └── users/             # User management service
├── migrations/                # Database migration files
├── pkg/pb/                    # Generated Protocol Buffer code
├── proto/                     # Protocol Buffer definitions
│   ├── auth.proto             # Authentication API
│   └── users.proto            # User management API
├── bin/                       # Compiled binaries and tools
├── docker-compose.yaml        # Docker services configuration
├── Dockerfile                 # Application container
├── Makefile                   # Build automation
├── go.mod                     # Go module dependencies
└── README.md
```

## Architecture Components

### API Layer (`internal/api`)
Implements both gRPC and HTTP REST interfaces:
- **gRPC**: Native high-performance RPC with Protocol Buffers
- **HTTP**: RESTful JSON API via gRPC-Gateway
- **Swagger**: Interactive API documentation at `/swagger/index.html`
- **Middleware**: Authentication, logging, error handling, CORS

### Services (`internal/services`)
Business logic layer:
- **auth**: User authentication (login, logout, refresh, validate)
- **users**: User CRUD operations with search and filtering

### Repository (`internal/repository`)
Data access layer with Squirrel query builder for PostgreSQL.

### Clients (`internal/pkg/clients`)
Integration clients for external services:

#### Database Client (`db`)
- Connection pooling with pgx/v5
- Transaction management
- Health checks

#### NATS Client (`nats`)
- JetStream support for persistent messaging
- Publisher and consumer patterns
- Stream and consumer management
- Request-reply messaging
- Context metadata propagation

#### S3/MinIO Client (`s3`)
- Bucket management
- File upload/download
- File deletion
- Compatible with AWS S3 and MinIO

#### Email Client (`mail`)
- SMTP integration with TLS
- Attachment support
- HTML and plain text emails

#### Chrome Client (`chrome`)
- Headless browser automation
- HTML to PDF conversion
- Configurable page settings (landscape/portrait)
- Remote Chrome connection

### Servers (`internal/pkg/servers`)
Server implementations:

#### gRPC Server (`grpc`)
- TLS support (optional)
- Interceptor chain for middleware
- Reflection API for development
- Health checks

#### HTTP Server (`http`)
- Gin web framework
- gRPC-Gateway integration
- Swagger UI
- CORS middleware
- Graceful shutdown

#### NATS Server (`nats`)
- Message consumer management
- Concurrent message processing
- Error handling and retries

### Utilities (`internal/pkg`)

#### JWT (`jwt`)
- Access and refresh token generation
- Token validation and parsing
- Configurable expiration

#### Logger (`logger`)
- Structured logging with Zap
- Context-aware logging
- Request ID tracking
- Log level configuration

#### Metadata (`metadata`)
- Request ID propagation
- User context (ID, name)
- IP address tracking
- Organization ID support

#### Errors (`errors`)
- Typed errors: BadRequest, Unauthorized, Forbidden, NotFound
- gRPC error code mapping
- HTTP status code mapping

#### Password (`pwd`)
- Bcrypt password hashing
- Secure password comparison

#### Utils (`utils`)
- UUID generation and validation
- Map utilities (keys, values)
- Set operations
- String manipulation
- Limit/offset pagination
- GroupBy operations
- Pointer helpers

## Docker Services

### PostgreSQL
- **Image**: postgres:16
- **Port**: 5432
- **User/Password**: postgres/postgres
- **Extensions**: pg_stat_statements
- **Volume**: pg-data

### MinIO (S3-Compatible Storage)
- **Image**: minio/minio:latest
- **API Port**: 9000
- **Console Port**: 9001
- **User/Password**: admin/password
- **Volume**: minio-data

### Headless Chrome
- **Image**: browserless/chrome:latest
- **Port**: 3000
- **Max Concurrent Sessions**: 10
- **Connection Timeout**: 60s
- **Shared Memory**: 2GB

## Configuration

Multi-source configuration with priority:
1. Command-line flags
2. Environment variables (prefix: `BOILERPLATE_`)
3. Configuration file (`config.yaml`)
4. Default values

### Environment Variables

```bash
# Logging
BOILERPLATE_LOG_LEVEL=info  # debug, info, warn, error

# API Server
BOILERPLATE_API_HOST=0.0.0.0
BOILERPLATE_API_HTTP_PORT=8080
BOILERPLATE_API_GRPC_PORT=8082

# JWT Authentication
BOILERPLATE_API_ACCESS_PRIVATE_KEY=your-secret-key
BOILERPLATE_API_ACCESS_TOKEN_TTL=3600      # 1 hour
BOILERPLATE_API_REFRESH_TOKEN_TTL=604800   # 7 days

# PostgreSQL Database
BOILERPLATE_DB_HOST=localhost
BOILERPLATE_DB_PORT=5432
BOILERPLATE_DB_USER=postgres
BOILERPLATE_DB_PASSWORD=postgres
BOILERPLATE_DB_NAME=boilerplate
BOILERPLATE_DB_SSL_MODE=disable

# NATS Messaging
BOILERPLATE_NATS_HOST=localhost
BOILERPLATE_NATS_PORT=4222

# MinIO/S3 Storage
BOILERPLATE_S3_ENDPOINT=localhost:9000
BOILERPLATE_S3_ACCESS_KEY=admin
BOILERPLATE_S3_SECRET_KEY=password
BOILERPLATE_S3_BUCKET=boilerplate
BOILERPLATE_S3_USE_SSL=false

# Email Service
BOILERPLATE_MAIL_HOST=smtp.gmail.com
BOILERPLATE_MAIL_PORT=587
BOILERPLATE_MAIL_USER=your-email@gmail.com
BOILERPLATE_MAIL_PASSWORD=your-password
BOILERPLATE_MAIL_FROM=noreply@example.com

# Chrome PDF Service
BOILERPLATE_CHROME_HOST=localhost
BOILERPLATE_CHROME_PORT=3000
BOILERPLATE_CHROME_TIMEOUT=30  # seconds
```

## Getting Started

### Prerequisites
- Go 1.25+
- Docker and Docker Compose
- Make

### Quick Start

```bash
# Clone the repository
git clone <repository-url>
cd boilerplate-go

# Start infrastructure services (PostgreSQL, MinIO, Chrome)
make up

# Run the application
make run
```

The application will:
1. Connect to PostgreSQL and run migrations
2. Start gRPC server on port 8082
3. Start HTTP server on port 8080
4. Serve Swagger UI at http://localhost:8080/swagger/index.html

### Build Commands

```bash
# Build the binary
make build

# Build with version
VERSION=v1.0.0 make build

# Run tests with coverage
make test

# Run linters on changed files
make lint

# Run linters on all files
make lint-full

# Stop Docker services
make down
```

### Code Generation

```bash
# Generate all (mocks, swagger, protobuf)
make generate

# Generate only mocks
make gen-mocks

# Generate only Swagger docs
make gen-swag

# Generate only Protocol Buffer code
make gen-proto

# Export proto dependencies for IDE
make proto-deps
```

## API Documentation

### Swagger UI
Interactive API documentation available at:
- **URL**: http://localhost:8080/swagger/index.html
- **Features**: Try-it-out functionality, request/response examples, authentication

### Available APIs

#### Authentication API (`/api/auth`)
- `POST /api/auth/login` - User login (returns access & refresh tokens)
- `POST /api/auth/logout` - User logout
- `POST /api/auth/refresh` - Refresh access token
- `GET /api/auth/me` - Get current user info
- `POST /api/auth/validate` - Validate token

#### Users API (`/api/users`)
- `POST /api/users` - Create user
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete user
- `POST /api/users/search` - Search users with filters

## Working with Protocol Buffers

### Tools
The project uses [Buf](https://buf.build/) for Protocol Buffer management:
- `protoc-gen-go` - Generate Go structs
- `protoc-gen-go-grpc` - Generate gRPC services
- `protoc-gen-grpc-gateway` - Generate HTTP/REST gateway
- `protoc-gen-openapiv2` - Generate OpenAPI/Swagger specs
- `protoc-gen-validate` - Generate validation code

### Adding a New Service

1. Create a `.proto` file in `proto/`:
```protobuf
syntax = "proto3";
package myservice;

import "google/api/annotations.proto";
import "validate/validate.proto";

service MyServiceAPI {
  rpc DoSomething (DoSomethingRequest) returns (DoSomethingResponse) {
    option (google.api.http) = {
      post: "/api/myservice/action"
      body: "*"
    };
  }
}
```

2. Generate code:
```bash
make gen-proto
```

3. Implement the service in `internal/services/myservice/`

4. Register handlers in `internal/api/grpc/handlers/` and `internal/api/http/handlers/`

## Testing

### Run Tests
```bash
# Run all tests with coverage
make test

# Run tests for specific package
go test -v ./internal/services/auth/...

# Run tests with race detector
go test -race ./...
```

### Test Database
Tests use a separate database (`boilerplate_test`) which is automatically created and cleaned before each test run.

### Mocks
Generated using Mockery. Mocks are located in `mocks/` directories next to the interfaces they implement.

## Project Conventions

### Error Handling
Use typed errors from `internal/pkg/errors`:
```go
return errors.NewNotFoundError("user not found")
return errors.NewBadRequestError("invalid email")
return errors.NewUnauthorizedError("invalid credentials")
return errors.NewForbiddenError("access denied")
```

### Logging
Use structured logging with context:
```go
logger.Info(ctx, "user created", zap.String("user_id", userID))
logger.Error(ctx, "database error", zap.Error(err))
```

### Context Metadata
Propagate request metadata using the metadata package:
```go
ctx = metadata.SetRequestID(ctx, requestID)
ctx = metadata.SetUserID(ctx, userID)
requestID := metadata.GetRequestID(ctx)
```

## Contributing

1. Follow Go best practices and idioms
2. Write tests for new functionality
3. Run linters before committing: `make lint`
4. Update API documentation in `.proto` files
5. Generate code after changes: `make generate`

## License

This project is licensed under the MIT License.
