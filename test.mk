# Makefile para pruebas
TEST_COMPOSE = docker compose -f docker-compose.test.yml -p flow-wallet-api-test

# Cargar variables de entorno desde .env
include .env
export

.PHONY: test-setup
test-setup:
	@echo "Configurando entorno de pruebas..."
	@$(TEST_COMPOSE) down -v || true
	@$(TEST_COMPOSE) build
	@$(TEST_COMPOSE) up -d

.PHONY: test-wait
test-wait:
	@echo "Esperando a que los servicios estén listos..."
	@sleep 10
	@curl -s http://localhost:3001/v1/health/ready || (echo "API no está lista aún, esperando 10 segundos más..." && sleep 10)
	@echo "Servicios listos para pruebas"

.PHONY: test-run
test-run:
	@echo "Ejecutando pruebas..."
	@echo "NOTA: Para ejecutar las pruebas HTTP, abre el archivo test-suite/run-tests.http en un cliente HTTP como REST Client para VSCode"
	@echo "O ejecuta manualmente los archivos .http con un cliente como Insomnia o Postman"
	@echo "Servicios corriendo en:"
	@echo "  API: http://localhost:3001/v1"
	@echo "  Base de datos: localhost:5433"
	@echo "  Emulador Flow: localhost:3570"
	@echo "  Redis: localhost:6380"

.PHONY: test-clean
test-clean:
	@echo "Limpiando entorno de pruebas..."
	@$(TEST_COMPOSE) down -v

.PHONY: test-logs
test-logs:
	@$(TEST_COMPOSE) logs -f

.PHONY: test
test: test-setup test-wait test-run
	@echo "Entorno de pruebas listo. Usa 'make test-clean' cuando hayas terminado."

.PHONY: test-full
test-full: test
	@echo "Pruebas completadas. Limpiando entorno..."
	@$(TEST_COMPOSE) down -v

# Comandos para pruebas de integración en Go
.PHONY: test-go
test-go:
	@echo "Ejecutando pruebas de integración en Go..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v

.PHONY: test-go-health
test-go-health:
	@echo "Ejecutando pruebas de salud..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestHealthCheck

.PHONY: test-go-system
test-go-system:
	@echo "Ejecutando pruebas de sistema..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestSystemEndpoints

.PHONY: test-go-accounts
test-go-accounts:
	@echo "Ejecutando pruebas de cuentas..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestAccountsFlow

.PHONY: test-go-account-sign
test-go-account-sign:
	@echo "Ejecutando pruebas de firma de transacciones..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestAccountSignTransaction

.PHONY: test-go-transactions
test-go-transactions:
	@echo "Ejecutando pruebas de transacciones..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestTransactionEndpoints

.PHONY: test-go-scripts
test-go-scripts:
	@echo "Ejecutando pruebas de scripts..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestScripts

.PHONY: test-go-tokens
test-go-tokens:
	@echo "Ejecutando pruebas de tokens..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestTokens

.PHONY: test-go-token-operations
test-go-token-operations:
	@echo "Ejecutando pruebas de operaciones de tokens..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestTokenOperations

.PHONY: test-go-watchlist-ops
test-go-watchlist-ops:
	@echo "Ejecutando pruebas de watchlist y ops..."
	@cd integration-tests && FLOW_WALLET_ADMIN_ADDRESS=$(FLOW_WALLET_ADMIN_ADDRESS) go test -v -run TestWatchlistAndOps

.PHONY: test-integration
test-integration: test-setup test-wait test-go
	@echo "Pruebas de integración completadas."

.PHONY: test-integration-full
test-integration-full: test-setup test-wait test-go test-clean
	@echo "Pruebas de integración completadas y entorno limpiado." 