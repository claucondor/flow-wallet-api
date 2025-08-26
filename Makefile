# Simplify Docker Compose detection - just use the command directly
# This works with both Docker Compose V2 and V1
DOCKER_COMPOSE_CMD = docker compose

# Detect OS for local commands
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    OPEN_CMD = xdg-open
    KILL_CMD = pkill -f
else ifeq ($(UNAME_S),Darwin)
    OPEN_CMD = open
    KILL_CMD = pkill -f
else
    OPEN_CMD = echo "Unsupported OS for open command:"
    KILL_CMD = echo "Unsupported OS for kill command:"
endif

# Only check for flow and go when running commands that need them
check-flow:
ifeq (, $(shell which flow))
	$(error "No flow in PATH. Install Flow CLI: https://docs.onflow.org/flow-cli/install/")
endif

check-go:
ifeq (, $(shell which go))
	$(error "No go in PATH")
endif

dev = $(DOCKER_COMPOSE_CMD) -f docker-compose.dev.yml -p flow-wallet-api-dev
lightweight = $(DOCKER_COMPOSE_CMD) -f docker-compose.lightweight.yml -p flow-wallet-api-lightweight
lightweight-testnet = $(DOCKER_COMPOSE_CMD) -f docker-compose.lightweight-testnet.yml -p flow-wallet-api-lightweight-testnet
lightweight-mainnet = $(DOCKER_COMPOSE_CMD) -f docker-compose.lightweight-mainnet.yml -p flow-wallet-api-lightweight-mainnet
test-suite = $(DOCKER_COMPOSE_CMD) -f docker-compose.test-suite.yml -p flow-wallet-api-test

.PHONY: dev
dev:
	@$(dev) up --remove-orphans -d db pgadmin emulator redis

.PHONY: stop
stop:
	@$(dev) stop

.PHONY: down
down:
	@$(dev) down --remove-orphans

.PHONY: reset
reset: down dev

.PHONY: lightweight
lightweight:
	@$(lightweight) down --remove-orphans || true
	@$(lightweight) build api
	@$(lightweight) up -d

.PHONY: lightweight-logs
lightweight-logs:
	@$(lightweight) logs -f

.PHONY: lightweight-stop
lightweight-stop:
	@$(lightweight) stop

.PHONY: lightweight-down
lightweight-down:
	@$(lightweight) down --remove-orphans

.PHONY: lightweight-reset
lightweight-reset: lightweight-down lightweight

.PHONY: lightweight-idempotent
lightweight-idempotent:
	@echo "🚀 Starting Phoenix Wallet API in lightweight mode with idempotency enabled"
	@$(lightweight) down --remove-orphans || true
	@$(lightweight) build api
	@FLOW_WALLET_LIGHTWEIGHT_IDEMPOTENCY=true $(lightweight) up -d
	@echo "✅ Services started with idempotency enabled"
	@echo "📋 Available endpoints:"
	@echo "   🌐 API: http://localhost:3000/v1"
	@echo "   📚 Documentation: http://localhost:8080"
	@echo "🔑 Idempotency is ENABLED - use 'Idempotency-Key' header in POST requests"

# Testnet commands
.PHONY: lightweight-testnet
lightweight-testnet:
	@echo "🌐 Starting Phoenix Wallet API in lightweight mode (Flow Testnet)"
	@$(lightweight-testnet) down --remove-orphans || true
	@$(lightweight-testnet) build api
	@$(lightweight-testnet) up -d
	@echo "✅ Services started on Flow Testnet"
	@echo "📋 Available endpoints:"
	@echo "   🌐 API: http://localhost:3000/v1"
	@echo "   📚 Documentation: http://localhost:8080"
	@echo "   🔗 Network: Flow Testnet"

.PHONY: lightweight-testnet-idempotent
lightweight-testnet-idempotent:
	@echo "🌐 Starting Phoenix Wallet API in lightweight mode (Flow Testnet) with idempotency"
	@$(lightweight-testnet) down --remove-orphans || true
	@$(lightweight-testnet) build api
	@FLOW_WALLET_LIGHTWEIGHT_IDEMPOTENCY=true $(lightweight-testnet) up -d
	@echo "✅ Services started on Flow Testnet with idempotency enabled"
	@echo "📋 Available endpoints:"
	@echo "   🌐 API: http://localhost:3000/v1"
	@echo "   📚 Documentation: http://localhost:8080"
	@echo "   🔗 Network: Flow Testnet"
	@echo "🔑 Idempotency is ENABLED - use 'Idempotency-Key' header in POST requests"

.PHONY: lightweight-testnet-stop
lightweight-testnet-stop:
	@$(lightweight-testnet) stop

.PHONY: lightweight-testnet-down
lightweight-testnet-down:
	@$(lightweight-testnet) down --remove-orphans

.PHONY: lightweight-testnet-logs
lightweight-testnet-logs:
	@$(lightweight-testnet) logs -f

# Mainnet commands
.PHONY: lightweight-mainnet
lightweight-mainnet:
	@echo "🏦 Starting Phoenix Wallet API in lightweight mode (Flow Mainnet)"
	@echo "⚠️  WARNING: You are connecting to MAINNET - real funds will be used!"
	@$(lightweight-mainnet) down --remove-orphans || true
	@$(lightweight-mainnet) build api
	@$(lightweight-mainnet) up -d
	@echo "✅ Services started on Flow Mainnet"
	@echo "📋 Available endpoints:"
	@echo "   🌐 API: http://localhost:3000/v1"
	@echo "   📚 Documentation: http://localhost:8080"
	@echo "   🔗 Network: Flow Mainnet"
	@echo "⚠️  CAUTION: This is MAINNET - real transactions will cost real FLOW!"

.PHONY: lightweight-mainnet-idempotent
lightweight-mainnet-idempotent:
	@echo "🏦 Starting Phoenix Wallet API in lightweight mode (Flow Mainnet) with idempotency"
	@echo "⚠️  WARNING: You are connecting to MAINNET - real funds will be used!"
	@$(lightweight-mainnet) down --remove-orphans || true
	@$(lightweight-mainnet) build api
	@FLOW_WALLET_LIGHTWEIGHT_IDEMPOTENCY=true $(lightweight-mainnet) up -d
	@echo "✅ Services started on Flow Mainnet with idempotency enabled"
	@echo "📋 Available endpoints:"
	@echo "   🌐 API: http://localhost:3000/v1"
	@echo "   📚 Documentation: http://localhost:8080"
	@echo "   🔗 Network: Flow Mainnet"
	@echo "🔑 Idempotency is ENABLED - use 'Idempotency-Key' header in POST requests"
	@echo "⚠️  CAUTION: This is MAINNET - real transactions will cost real FLOW!"

.PHONY: lightweight-mainnet-stop
lightweight-mainnet-stop:
	@$(lightweight-mainnet) stop

.PHONY: lightweight-mainnet-down
lightweight-mainnet-down:
	@$(lightweight-mainnet) down --remove-orphans

.PHONY: lightweight-mainnet-logs
lightweight-mainnet-logs:
	@$(lightweight-mainnet) logs -f

# Help command
.PHONY: help
help:
	@echo "🚀 Phoenix Wallet API - Available Commands"
	@echo "=========================================="
	@echo ""
	@echo "📦 Standard Development:"
	@echo "  make dev                    - Start full stack (PostgreSQL + Redis + Emulator)"
	@echo "  make stop                   - Stop development services"
	@echo "  make down                   - Stop and remove development containers"
	@echo "  make reset                  - Reset development environment"
	@echo ""
	@echo "🪶 Lightweight Mode (Local Emulator):"
	@echo "  make lightweight            - Start lightweight mode (SQLite, no idempotency)"
	@echo "  make lightweight-idempotent - Start lightweight mode with idempotency"
	@echo "  make lightweight-stop       - Stop lightweight services"
	@echo "  make lightweight-down       - Stop and remove lightweight containers"
	@echo "  make lightweight-logs       - View lightweight logs"
	@echo ""
	@echo "🌐 Lightweight Mode (Flow Testnet):"
	@echo "  make lightweight-testnet            - Connect to Flow Testnet"
	@echo "  make lightweight-testnet-idempotent - Connect to Flow Testnet with idempotency"
	@echo "  make lightweight-testnet-stop       - Stop testnet services"
	@echo "  make lightweight-testnet-down       - Stop and remove testnet containers"
	@echo "  make lightweight-testnet-logs       - View testnet logs"
	@echo ""
	@echo "🏦 Lightweight Mode (Flow Mainnet):"
	@echo "  make lightweight-mainnet            - Connect to Flow Mainnet ⚠️"
	@echo "  make lightweight-mainnet-idempotent - Connect to Flow Mainnet with idempotency ⚠️"
	@echo "  make lightweight-mainnet-stop       - Stop mainnet services"
	@echo "  make lightweight-mainnet-down       - Stop and remove mainnet containers"
	@echo "  make lightweight-mainnet-logs       - View mainnet logs"
	@echo ""
	@echo "🧪 Testing:"
	@echo "  make test                   - Run tests with emulator"
	@echo "  make run-tests              - Run tests only"
	@echo "  make test-clean             - Clean test cache and run tests"
	@echo "  make run-test-suite         - Run full dockerized test suite"
	@echo "  make stop-test-suite        - Stop test suite containers"
	@echo ""
	@echo "🔧 Local Development (No Docker):"
	@echo "  make local-start            - Start services locally"
	@echo "  make local-stop             - Stop local services"
	@echo "  make local-status           - Check local services status"
	@echo ""
	@echo "⚠️  Important Notes:"
	@echo "  - For testnet/mainnet: Configure FLOW_WALLET_ADMIN_ADDRESS and FLOW_WALLET_ADMIN_PRIVATE_KEY"
	@echo "  - Mainnet commands use real FLOW tokens - test on testnet first!"
	@echo "  - Idempotency prevents duplicate operations using 'Idempotency-Key' header"
	@echo ""
	@echo "📚 Endpoints (when running):"
	@echo "  - API: http://localhost:3000/v1"
	@echo "  - Documentation: http://localhost:8080"
	@echo "  - Flow Emulator: localhost:3569 (local only)"
	@echo ""
	@echo "📖 Documentation Commands:"
	@echo "  make docs-dev             - Start documentation in development mode"
	@echo "  make docs-build           - Build documentation for production"
	@echo "  make docs-serve           - Serve built documentation locally"

# Documentation commands
.PHONY: docs-dev
docs-dev:
	@echo "🚀 Starting documentation in development mode..."
	@cd docs && npm install && npm start

.PHONY: docs-build
docs-build:
	@echo "🏗️  Building documentation for production..."
	@cd docs && npm install && npm run build

.PHONY: docs-serve
docs-serve:
	@echo "📚 Serving built documentation..."
	@cd docs && npm run serve

.PHONY: run-tests
run-tests: check-go
	@go test ./... -p 1

.PHONY: test
test: check-flow check-go start-emulator deploy run-tests

.PHONY: test-clean
test-clean: check-go clean-test-cache test

.PHONY: clean-test-cache
clean-test-cache: check-go
	@go clean -testcache

.PHONY: deploy
deploy: check-flow
	@cd flow && flow project deploy --update

.PHONY: start-emulator
start-emulator: check-flow emulator.pid
	@sleep 1

.PHONY: stop-emulator
stop-emulator: check-flow emulator.pid
	@kill `cat $<` && rm $<

emulator.pid: check-flow
	@cd flow && { flow emulator -b 100ms & echo $$! > ../$@; }

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: run-test-suite
run-test-suite:
	@$(test-suite) build flow api
	@$(test-suite) up --remove-orphans -d db redis flow
	@$(test-suite) unpause \
	; echo "\nRunning tests, hang on...\n" \
	; $(test-suite) run --rm api go test ./... -p 1 \
	; echo "\nRunning linter, hang on...\n" \
	; $(test-suite) run --rm lint golangci-lint run \
	; $(test-suite) pause

.PHONY: stop-test-suite
stop-test-suite:
	@$(test-suite) down --remove-orphans

.PHONY: clean-test-suite
clean-test-suite:
	@$(test-suite) run --rm api go clean -testcache

# Local commands without Docker
.PHONY: local-start
local-start: check-flow check-go
	@echo "Starting Phoenix Wallet API services locally..."
	@echo "Creating necessary directories..."
	@mkdir -p ./data
	
	@echo "Starting Flow Emulator..."
	@cd flow && { flow emulator --persist --storage-dir="../data/flowdb" --port=3569 & echo $$! > ../emulator.pid; }
	@echo "Flow Emulator started with PID: $$(cat emulator.pid)"
	@sleep 3
	
	@echo "Starting Swagger UI (using npx)..."
	@npx swagger-ui-dist@latest serve ./openapi.yml -p 8081 & echo $$! > swagger.pid
	@echo "Swagger UI started with PID: $$(cat swagger.pid)"
	@sleep 2
	
	@echo "Building and starting Phoenix Wallet API..."
	@go build -o ./phoenix-wallet-api
	@FLOW_WALLET_LIGHTWEIGHT_MODE=true \
	FLOW_WALLET_ACCESS_API_HOST=localhost:3569 \
	FLOW_WALLET_CHAIN_ID=flow-emulator \
	FLOW_WALLET_DATABASE_DSN=./data/wallet.db \
	./phoenix-wallet-api & echo $$! > api.pid
	@echo "Phoenix Wallet API started with PID: $$(cat api.pid)"
	
	@echo "All services started successfully!"
	@echo "API: http://localhost:3000"
	@echo "Swagger UI: http://localhost:8081"
	@echo "Flow Emulator: localhost:3569"
	@$(OPEN_CMD) http://localhost:8081

.PHONY: local-stop
local-stop:
	@echo "Stopping all Phoenix Wallet API services..."
	@if [ -f emulator.pid ]; then kill $$(cat emulator.pid) && rm emulator.pid && echo "Flow Emulator stopped"; else echo "Flow Emulator not running"; fi
	@if [ -f swagger.pid ]; then kill $$(cat swagger.pid) && rm swagger.pid && echo "Swagger UI stopped"; else echo "Swagger UI not running"; fi
	@if [ -f api.pid ]; then kill $$(cat api.pid) && rm api.pid && echo "Phoenix Wallet API stopped"; else echo "Phoenix Wallet API not running"; fi
	@echo "All services stopped successfully!"

.PHONY: local-status
local-status:
	@echo "Checking status of Phoenix Wallet API services..."
	@if [ -f emulator.pid ] && ps -p $$(cat emulator.pid) > /dev/null; then echo "Flow Emulator: Running (PID: $$(cat emulator.pid))"; else echo "Flow Emulator: Not running"; fi
	@if [ -f swagger.pid ] && ps -p $$(cat swagger.pid) > /dev/null; then echo "Swagger UI: Running (PID: $$(cat swagger.pid))"; else echo "Swagger UI: Not running"; fi
	@if [ -f api.pid ] && ps -p $$(cat api.pid) > /dev/null; then echo "Phoenix Wallet API: Running (PID: $$(cat api.pid))"; else echo "Phoenix Wallet API: Not running"; fi
