VERSION ?= dev
DB_NAME=boilerplate
TEST_DB_NAME=boilerplate_test

up:
	docker-compose up -d

down:
	docker-compose down

prepare-run:
	docker exec boilerplate-go-postgres-1 createdb -Upostgres $(DB_NAME) || true;

run: prepare-run
	go run ./cmd/boilerplate/main.go

.PHONY: build
build:
	go build -ldflags "-X boilerplate/internal/pkg/version.gitVersion=${VERSION} -X 'boilerplate/internal/pkg/version.buildDate=$(shell date -u '+%Y-%m-%d %H:%M:%S')' -s -w" -trimpath -o bin/boilerplate ./cmd/boilerplate/main.go

lint: 
	go tool golangci-lint run \
	--new-from-rev=master \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

lint-full: 
	go tool golangci-lint run \
	--config=.golangci.yaml \
	--max-issues-per-linter=1000 \
	--max-same-issues=1000 \
	./...

prepare-test:
	@if [ -z "$$GITLAB_CI" ]; then \
		echo "Preparing test database locally"; \
		docker exec boilerplate-go-postgres-1 dropdb --if-exists -Upostgres $(TEST_DB_NAME); \
		docker exec boilerplate-go-postgres-1 createdb -Upostgres $(TEST_DB_NAME); \
	fi

test: prepare-test
	go test -v -count=1 ./... -coverprofile coverage.out.tmp 
	@grep -vE "mock.go|pb.go" coverage.out.tmp > coverage.out
	@rm -f coverage.out.tmp
	@go tool cover -func coverage.out | grep total | awk '{print "Coverage percent: " $$3}'
	@rm -f coverage.out

gen-mock: 
	go tool mockery

gen-swag:
	go tool swag fmt	
	go tool swag init --parseDependency --parseInternal -g handlers.go -d internal/api/handlers -o internal/api/swagger

generate: gen-mock gen-swag