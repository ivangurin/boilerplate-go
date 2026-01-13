LOCAL_BIN:=$(CURDIR)/bin
VERSION ?= dev
DB_NAME=boilerplate
TEST_DB_NAME=boilerplate_test

GREEN := \033[32m
RED := \033[31m
RESET := \033[0m

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(RESET) %s\n", $$1, $$2}'

.PHONY: up
up: ## Start Docker Compose services
	docker-compose up -d

.PHONY: down
down: ## Stop Docker Compose services
	docker-compose down

.PHONY: prepare-run
prepare-run: ## Prepare database for running
	@echo "$(GREEN)▶$(RESET) Preparing database..."
	@docker exec boilerplate-go-postgres-1 createdb -Upostgres $(DB_NAME) 2> /dev/null || true;

.PHONY: run
run: prepare-run ## Run the application
	@echo "$(GREEN)▶$(RESET) Running application..."
	@go run ./cmd/boilerplate/main.go

.PHONY: build
build: ## Build the binary
	go build -ldflags "-X boilerplate/internal/pkg/version.gitVersion=${VERSION} -X 'boilerplate/internal/pkg/version.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S')' -s -w" -trimpath -o bin/boilerplate ./cmd/boilerplate/main.go

.PHONY: lint
lint: ## Run linters on changed files
	@echo "$(GREEN)▶$(RESET) Running linters on changed files..."
	@go tool golangci-lint run \
		--fix \
		--new-from-rev=master \
		--config=.golangci.yaml \
		--max-issues-per-linter=1000 \
		--max-same-issues=1000 \
		./...

.PHONY: lint-full
lint-full: ## Run linters on all files
	@echo "$(GREEN)▶$(RESET) Running linters on full codebase..."
	@go tool golangci-lint run \
		--fix \
		--config=.golangci.yaml \
		--max-issues-per-linter=1000 \
		--max-same-issues=1000 \
		./...

.PHONY: prepare-test
prepare-test: ## Prepare test database
	@echo "$(GREEN)▶$(RESET) Preparing test database..."
	@if [ -z "$$GITLAB_CI" ]; then \
		docker exec boilerplate-go-postgres-1 dropdb --if-exists -Upostgres $(TEST_DB_NAME); \
		docker exec boilerplate-go-postgres-1 createdb -Upostgres $(TEST_DB_NAME); \
	fi

.PHONY: test
test: prepare-test ## Run tests with coverage
	@echo "$(GREEN)▶$(RESET) Running tests with coverage..."
	@go test -v -count=1 ./... -coverprofile coverage.out.tmp 
	@grep -vE "mock.go|pb.go" coverage.out.tmp > coverage.out
	@rm -f coverage.out.tmp
	@go tool cover -func coverage.out | grep total | awk '{print "Coverage percent: " $$3}'
	@rm -f coverage.out

.PHONY: clean-mocks
clean-mocks: ## Clean old mocks
	@echo "$(RED)▶$(RESET) Cleaning old mocks..."
	@rm -rf $(find . -type d -name mocks)

.PHONY: gen-mocks
gen-mocks: clean-mocks ## Generate mocks
	@echo "$(GREEN)▶$(RESET) Generating mocks..."
	@go tool mockery

.PHONY: gen-swag
gen-swag: ## Generate swagger docs
	@echo "$(GREEN)▶$(RESET) Generating swagger docs..."
	@go tool swag fmt	
	@go tool swag init --parseDependency --parseInternal -g handlers.go -d internal/api/http/handlers -o internal/api/swagger

.PHONY: proto-deps
proto-deps: ## Export proto dependencies for IDE
	@echo "$(GREEN)▶$(RESET) Exporting proto dependencies for IDE..."
	@rm -rf .proto-deps
	@go tool buf export buf.build/envoyproxy/protoc-gen-validate --output .proto-deps
	@go tool buf export buf.build/googleapis/googleapis --output .proto-deps
	@go tool buf export buf.build/grpc-ecosystem/grpc-gateway --output .proto-deps
	@echo "$(GREEN)▶$(RESET) Proto dependencies exported to .proto-deps/"

.PHONY: install-proto-plugins
install-proto-plugins: ## Install proto plugins
	@echo "$(GREEN)▶$(RESET) Installing proto plugins..."
	@for plugin in google.golang.org/protobuf/cmd/protoc-gen-go \
		google.golang.org/grpc/cmd/protoc-gen-go-grpc \
		github.com/envoyproxy/protoc-gen-validate \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
		github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2; do \
		GOBIN=$(LOCAL_BIN) go install $$plugin@latest || exit 1; \
	done

.PHONY: gen-proto
gen-proto: install-proto-plugins ## Generate proto files
	@echo "$(GREEN)▶$(RESET) Updating buf dependencies..."
	@go tool buf dep update
	@echo "$(GREEN)▶$(RESET) Generating proto files..."
	@PATH=$(LOCAL_BIN):"$(PATH)" go tool buf generate
	@echo "$(GREEN)▶$(RESET) Merging swagger files..."
	@go tool swagger mixin --quiet --ignore-conflicts --compact -o internal/pkg/swagger/swagger.json pkg/pb/*.swagger.json
	@find pkg/pb -name "*.swagger.json" -type f -delete
	@echo "$(GREEN)▶$(RESET) Proto generation completed successfully"
	
.PHONY: lint-proto
lint-proto: ## Lint proto files
	@echo "$(GREEN)▶$(RESET) Linting proto files..."
	@go tool buf lint

.PHONY: format
format: ## Format code
	@echo "$(GREEN)▶$(RESET) Code formatting..."
	@go tool gofumpt -l -w -extra .

.PHONY: generate
generate: ## Run code generation (parallel)
	@echo "$(GREEN)▶$(RESET) Code generation (parallel)..."
	$(MAKE) -j 3 format gen-mocks gen-proto
	@echo "$(GREEN)▶$(RESET) Code generation completed successfully"

.PHONY: docker-build
docker-build: ## Build Docker image
	docker image rm boilerplate:latest || true
	docker build --build-arg VERSION=${VERSION} -t boilerplate:latest -f ./Dockerfile .

.PHONY: docker-run
docker-run: ## Run Docker container
	docker run -d -p 8080:8080 -e BOILERPLATE_DB_HOST=postgres -e BOILERPLATE_API_HOST=0.0.0.0 --name boilerplate --network boilerplate-network --rm boilerplate:latest

.PHONY: docker-stop
docker-stop: ## Stop Docker container
	docker stop boilerplate || true