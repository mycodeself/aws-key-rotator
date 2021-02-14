TESTS?=./...
MAIN=./main.go

run: 
	@echo "==> Running $(MAIN)..."
	@go run $(MAIN)

lint:
	@echo "==> Running linter..."
	@golangci-lint run

fmt:
	@echo "==> Running go fmt"
	@find . -name '*.go' | xargs gofmt -s -w

test:
	@echo "==> Running tests"
	@go test $(TESTS)