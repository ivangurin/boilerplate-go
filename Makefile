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

.PHONY: gen-proto
gen-proto:
	go tool buf generate
	go tool swagger mixin --ignore-conflicts -o internal/pkg/swagger/swagger.json pkg/pb/*.swagger.json
	@find pkg/pb -name "*.swagger.json" -type f -delete
	@echo "Swagger files merged successfully"

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