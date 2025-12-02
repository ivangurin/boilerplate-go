up:
	docker-compose up -d

down:
	docker-compose down

run:
	go run ./cmd/main.go

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
	./...g

test:
	go test -v -count=1 ./... -coverprofile coverage.out.tmp 
	@grep -vE "mock.go|pb.go" coverage.out.tmp > coverage.out
	@rm -f coverage.out.tmp
	@go tool cover -func coverage.out | grep total | awk '{print "Coverage percent: " $$3}'
	@rm -f coverage.out