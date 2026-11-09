.PHONY: run swagger test dev-run

run: 
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

COVERAGE_EXCLUDE = mocks|main.go|docs|test

test: 
	go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out
	
dev-run: swagger run
