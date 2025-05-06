include .env

LOCAL_BIN := $(CURDIR)/bin
GOLANGCI_LINT_VERSION := v2.1.5
GOLANGCI_LINT_BIN := $(LOCAL_BIN)/golangci-lint

.PHONY: install-lint
install-lint:
	@echo "Installing golangci-lint $(GOLANGCI_LINT_VERSION)..."
	@mkdir -p $(LOCAL_BIN)
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(LOCAL_BIN) $(GOLANGCI_LINT_VERSION)
	@$(GOLANGCI_LINT_BIN) --version

.PHONY: init-deps
init-deps:
	@echo "Installing developer tools..."
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest
	GOBIN=$(LOCAL_BIN) go install github.com/golang/mock/mockgen@latest

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@$(GOLANGCI_LINT_BIN) run ./... --config=.golangci.yml

.PHONY: build
build:
	@echo "Building application..."
	@go build -o $(LOCAL_BIN)/app cmd/server/main.go
	@echo "Application built at $(LOCAL_BIN)/app"

.PHONY: run
run: build
	@$(LOCAL_BIN)/app

.PHONY: clean
clean:
	@echo "Cleaning binaries..."
	@rm -rf $(LOCAL_BIN)/*

.PHONY: docker-up
docker-up:
	docker compose -f ./docker-compose.yaml up -d

.PHONY: docker-down
docker-down:
	docker compose -f ./docker-compose.yaml down

.PHONY: local-migration-status
local-migration-status:
	@$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} status -v

.PHONY: local-migration-up
local-migration-up:
	@$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} up -v

.PHONY: local-migration-down
local-migration-down:
	@$(LOCAL_BIN)/goose -dir ${MIGRATION_DIR} postgres ${PG_DSN} down -v

