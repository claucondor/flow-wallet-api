package integrationtests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/api"
)

func TestHealthCheck(t *testing.T) {
	t.Run("Health Ready", func(t *testing.T) {
		resp, _, err := apiClient.Get("/health/ready")
		if err != nil {
			t.Fatalf("Error al verificar health/ready: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}
	})

	t.Run("Health Liveness", func(t *testing.T) {
		resp, body, err := apiClient.Get("/health/liveness")
		if err != nil {
			t.Fatalf("Error al verificar health/liveness: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var healthStatus api.HealthStatus
		if err := json.Unmarshal(body, &healthStatus); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que los campos relevantes existan
		if healthStatus.WorkerCount <= 0 {
			t.Errorf("Se esperaba un número positivo de workers, obtenido: %d", healthStatus.WorkerCount)
		}

		if healthStatus.PoolCapacity <= 0 {
			t.Errorf("Se esperaba una capacidad positiva del pool, obtenido: %d", healthStatus.PoolCapacity)
		}
	})
}
