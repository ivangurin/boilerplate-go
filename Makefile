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
	docker exec boilerplate-go-postgres-1 createdb -Upostgres $(DB_NAME) || true;

.PHONY: run
run: prepare-run
	go run ./cmd/boilerplate/main.go

.PHONY: build
build:
	go build -ldflags "-X boilerplate/internal/pkg/version.gitVersion=${VERSION} -X 'boilerplate/internal/pkg/version.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S')' -s -w" -trimpath -o bin/boilerplate ./cmd/boilerplate/main.go

.PHONY: lint
lint: 
	go tool golangci-lint run \
	--new-from-rev=master \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

.PHONY: lint-full
lint-full: 
	go tool golangci-lint run \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

.PHONY: prepare-test
prepare-test:
	@if [ -z "$$GITLAB_CI" ]; then \
		echo "Preparing test database locally"; \
		docker exec boilerplate-go-postgres-1 dropdb --if-exists -Upostgres $(TEST_DB_NAME); \
		docker exec boilerplate-go-postgres-1 createdb -Upostgres $(TEST_DB_NAME); \
	fi

.PHONY: test
test: prepare-test
	go test -v -count=1 ./... -coverprofile coverage.out.tmp 
	@grep -vE "mock.go|pb.go" coverage.out.tmp > coverage.out
	@rm -f coverage.out.tmp
	@go tool cover -func coverage.out | grep total | awk '{print "Coverage percent: " $$3}'
	@rm -f coverage.out

.PHONY: gen-mock
gen-mock: 
	go tool mockery

.PHONY: gen-swag
gen-swag:
	go tool swag fmt	
	go tool swag init --parseDependency --parseInternal -g handlers.go -d internal/api/handlers -o internal/api/swagger

.PHONY: proto-deps
proto-deps:
	@echo "Exporting proto dependencies for IDE..."
	@rm -rf .proto-deps
	@go tool buf export buf.build/envoyproxy/protoc-gen-validate --output .proto-deps
	@go tool buf export buf.build/googleapis/googleapis --output .proto-deps
	@go tool buf export buf.build/grpc-ecosystem/grpc-gateway --output .proto-deps
	@echo "Proto dependencies exported to .proto-deps/"

.install-proto-plugins:
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	GOBIN=$(LOCAL_BIN) go install github.com/envoyproxy/protoc-gen-validate@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

.PHONY: gen-proto
gen-proto: .install-proto-plugins
	@echo "Updating buf dependencies..."
	@go tool buf dep update
	@echo "Generating proto files..."
	@PATH=$(LOCAL_BIN):"$(PATH)" go tool buf generate
	@echo "Merging swagger files..."
	@go tool swagger mixin --quiet --ignore-conflicts --compact -o internal/pkg/swagger/swagger.json pkg/pb/*.swagger.json
	@echo "Proto generation completed successfully"
	
.PHONY: lint-proto
lint-proto:
	go tool buf lint

.PHONY: generate
generate: gen-mock gen-swag gen-proto

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