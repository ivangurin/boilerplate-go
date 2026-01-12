LOCAL_BIN:=$(CURDIR)/bin
VERSION ?= dev
DB_NAME=boilerplate
TEST_DB_NAME=boilerplate_test

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: prepare-run
prepare-run:
	@echo "\033[32m▶\033[0m Preparing database..."
	@docker exec boilerplate-go-postgres-1 createdb -Upostgres $(DB_NAME) 2> /dev/null || true;

.PHONY: run
run: prepare-run
	@echo "\033[32m▶\033[0m Running application..."
	@go run ./cmd/boilerplate/main.go

.PHONY: build
build:
	go build -ldflags "-X boilerplate/internal/pkg/version.gitVersion=${VERSION} -X 'boilerplate/internal/pkg/version.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S')' -s -w" -trimpath -o bin/boilerplate ./cmd/boilerplate/main.go

.PHONY: lint
lint: 
	@echo "\033[32m▶\033[0m Running linters on changed files..."
	@go tool golangci-lint run \
	--new-from-rev=master \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

.PHONY: lint-full
lint-full: 
	@echo "\033[32m▶\033[0m Running linters on full codebase..."
	@go tool golangci-lint run \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

.PHONY: prepare-test
prepare-test:
	@echo "\033[32m▶\033[0m Preparing test database..."
	@if [ -z "$$GITLAB_CI" ]; then \
		docker exec boilerplate-go-postgres-1 dropdb --if-exists -Upostgres $(TEST_DB_NAME); \
		docker exec boilerplate-go-postgres-1 createdb -Upostgres $(TEST_DB_NAME); \
	fi

.PHONY: test
test: prepare-test
	@echo "\033[32m▶\033[0m Running tests with coverage..."
	@go test -v -count=1 ./... -coverprofile coverage.out.tmp 
	@grep -vE "mock.go|pb.go" coverage.out.tmp > coverage.out
	@rm -f coverage.out.tmp
	@go tool cover -func coverage.out | grep total | awk '{print "Coverage percent: " $$3}'
	@rm -f coverage.out

.PHONY: gen-mock
gen-mock: 
	@echo "\033[32m▶\033[0m Generating mocks..."
	@go tool mockery

.PHONY: gen-swag
gen-swag:
	@echo "\033[32m▶\033[0m Generating swagger docs..."
	@go tool swag fmt	
	@go tool swag init --parseDependency --parseInternal -g handlers.go -d internal/api/http/handlers -o internal/api/swagger

.PHONY: proto-deps
proto-deps:
	@echo "\033[32m▶\033[0m Exporting proto dependencies for IDE..."
	@rm -rf .proto-deps
	@go tool buf export buf.build/envoyproxy/protoc-gen-validate --output .proto-deps
	@go tool buf export buf.build/googleapis/googleapis --output .proto-deps
	@go tool buf export buf.build/grpc-ecosystem/grpc-gateway --output .proto-deps
	@echo "\033[32m▶\033[0m Proto dependencies exported to .proto-deps/"

.install-proto-plugins:
	@echo "\033[32m▶\033[0m Installing proto plugins..."
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

.PHONY: gen-proto
gen-proto: .install-proto-plugins
	@echo "\033[32m▶\033[0m Updating buf dependencies..."
	@go tool buf dep update
	@echo "\033[32m▶\033[0m Generating proto files..."
	@PATH=$(LOCAL_BIN):"$(PATH)" go tool buf generate
	@echo "\033[32m▶\033[0m Merging swagger files..."
	@go tool swagger mixin --quiet --ignore-conflicts --compact -o internal/pkg/swagger/swagger.json pkg/pb/*.swagger.json
	@find pkg/pb -name "*.swagger.json" -type f -delete
	@echo "\033[32m▶\033[0m Proto generation completed successfully"
	
.PHONY: lint-proto
lint-proto:
	@echo "\033[32m▶\033[0m Linting proto files..."
	@go tool buf lint

.PHONY: format
format:
	@echo "\033[32m▶\033[0m Code formatting..."
	@go tool gofumpt -l -w -extra .

.PHONY: generate
generate: format gen-mock gen-swag gen-proto
	@echo "\033[32m▶\033[0m Code generation completed successfully"

.PHONY: docker-build
docker-build:
	docker image rm boilerplate:latest || true
	docker build --build-arg VERSION=${VERSION} -t boilerplate:latest -f ./Dockerfile .

.PHONY: docker-run
docker-run:
	docker run -d -p 8080:8080 -e BOILERPLATE_DB_HOST=postgres -e BOILERPLATE_API_HOST=0.0.0.0 --name boilerplate --network boilerplate-network --rm boilerplate:latest

.PHONY: docker-stop
docker-stop:
	docker stop boilerplate || true