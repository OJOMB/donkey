PHONY: test
test:
	go test -v -cover ./internal/...

.PHONY: test-verbose
test-verbose:
	go test -v -race -coverprofile=coverage.out ./internal/...
	go tool cover -html=coverage.out

.PHONY: test-unit
test-unit:
	go test -short ./internal/...