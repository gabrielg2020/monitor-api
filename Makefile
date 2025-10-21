# Go parameters
GOCMD=go
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=monitor-api

# Test parameters
TEST_TIMEOUT=30s
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

# Colors for output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
NC=\033[0m

.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: deps
deps: ## Download dependencies
	@echo "$(YELLOW)Downloading dependencies...$(NC)"
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: test-deps
test-deps: ## Install test dependencies
	@echo "$(YELLOW)Installing test dependencies...$(NC)"
	$(GOGET) -u github.com/stretchr/testify
	$(GOGET) -u github.com/vektra/mockery/v2/...
	$(GOGET) -u github.com/DATA-DOG/go-sqlmock

.PHONY: test
test: ## Run all tests
	@echo "$(YELLOW)Running all tests...$(NC)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-unit
test-unit: ## Run unit tests only
	@echo "$(YELLOW)Running unit tests...$(NC)"
	$(GOTEST) -v -short -timeout $(TEST_TIMEOUT) ./internal/...

.PHONY: test-services
test-services: ## Run service layer tests
	@echo "$(YELLOW)Running service tests...$(NC)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./internal/services/...

PHONY: test-handlers
test-handlers: ## Run handler layer tests
	@echo "$(YELLOW)Running handler tests...$(NC)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./internal/api/handlers/...

.PHONY: test-repository
test-repository: ## Run repository layer tests
	@echo "$(YELLOW)Running repository tests...$(NC)"
	$(GOTEST) -v -timeout $(TEST_TIMEOUT) ./internal/repository/...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo "$(YELLOW)Running tests with coverage...$(NC)"
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_FILE) -covermode=atomic -timeout $(TEST_TIMEOUT) ./...
	@echo "$(GREEN)Coverage report generated: $(COVERAGE_FILE)$(NC)"

.PHONY: test-coverage-html
test-coverage-html: test-coverage ## Generate HTML coverage report
	@echo "$(YELLOW)Generating HTML coverage report...$(NC)"
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "$(GREEN)HTML coverage report generated: $(COVERAGE_HTML)$(NC)"

.PHONY: test-coverage-view
test-coverage-view: test-coverage-html ## View coverage report in browser
	@echo "$(YELLOW)Opening coverage report...$(NC)"
	@if command -v open > /dev/null; then \
		open $(COVERAGE_HTML); \
	elif command -v xdg-open > /dev/null; then \
		xdg-open $(COVERAGE_HTML); \
	else \
		echo "$(RED)Please open $(COVERAGE_HTML) manually$(NC)"; \
	fi

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "$(YELLOW)Running tests with verbose output...$(NC)"
	$(GOTEST) -v -count=1 -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-bench
test-bench: ## Run benchmark tests
	@echo "$(YELLOW)Running benchmark tests...$(NC)"
	$(GOTEST) -bench=. -benchmem -timeout $(TEST_TIMEOUT) ./...

.PHONY: test-clean
test-clean: ## Clean test cache and coverage files
	@echo "$(YELLOW)Cleaning test artifacts...$(NC)"
	$(GOCMD) clean -testcache
	rm -f $(COVERAGE_FILE) $(COVERAGE_HTML)
	@echo "$(GREEN)Test artifacts cleaned$(NC)"

.PHONY: test-report
test-report: ## Generate test report
	@echo "$(YELLOW)Generating test report...$(NC)"
	@$(GOTEST) -v -json -timeout $(TEST_TIMEOUT) ./... | tee test-report.json
	@echo "$(GREEN)Test report generated: test-report.json$(NC)"

.PHONY: lint
lint: ## Run linter
	@command -v golangci-lint >/dev/null 2>&1 || { echo "$(RED)golangci-lint is required but not installed$(NC)" >&2; exit 1; }
	@echo "$(YELLOW)Running linter...$(NC)"
	golangci-lint run ./...

.PHONY: fmt
fmt: ## Format code
	@echo "$(YELLOW)Formatting code...$(NC)"
	$(GOCMD) fmt ./...
	@echo "$(GREEN)Code formatted$(NC)"

.PHONY: build
build: ## Build the application
	@echo "$(YELLOW)Building application...$(NC)"
	$(GOCMD) build -o $(BINARY_NAME) ./cmd/monitor-api
	@echo "$(GREEN)Build complete: $(BINARY_NAME)$(NC)"

.PHONY: docs
docs: ## Generate API documentation
	@command -v swag >/dev/null 2>&1 || { echo "$(RED)swag is required but not installed.$(NC)" >&2; exit 1; }
	@echo "$(YELLOW)Generating API documentation...$(NC)"
	swag init -g ./cmd/monitor-api/main.go -o ./docs
	@echo "$(GREEN)API documentation generated in ./docs$(NC)"