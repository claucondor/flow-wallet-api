# Makefile para pruebas
TEST_COMPOSE = docker compose -f docker-compose.test.yml -p flow-wallet-api-test

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
	@echo "NOTA: Para ejecutar las pruebas, abre el archivo test-suite/run-tests.http en un cliente HTTP como REST Client para VSCode"
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