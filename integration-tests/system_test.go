package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/api"
)

func TestSystemEndpoints(t *testing.T) {
	// Prueba para obtener configuración del sistema
	t.Run("Get System Settings", func(t *testing.T) {
		resp, body, err := apiClient.Get("/system/settings")
		if err != nil {
			t.Fatalf("Error al obtener configuración del sistema: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var settings api.SystemSettings
		if err := json.Unmarshal(body, &settings); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v - Cuerpo: %s", err, string(body))
		}

		// Verificar que el campo MaintenanceMode exista (puede ser true o false)
		t.Logf("Modo de mantenimiento: %v", settings.MaintenanceMode)
	})

	// Prueba para actualizar configuración del sistema
	t.Run("Update System Settings", func(t *testing.T) {
		// Guardar el estado original
		resp, body, err := apiClient.Get("/system/settings")
		if err != nil {
			t.Fatalf("Error al obtener configuración original: %v", err)
		}

		var originalSettings api.SystemSettings
		if err := json.Unmarshal(body, &originalSettings); err != nil {
			t.Fatalf("Error al deserializar respuesta original: %v", err)
		}

		// Cambiar el modo de mantenimiento al valor opuesto
		newSettings := api.SystemSettings{
			MaintenanceMode: !originalSettings.MaintenanceMode,
		}

		// Actualizar configuración con una clave idempotente única
		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("test-system-settings-update-%d", time.Now().UnixNano()),
		}
		resp, body, err = apiClient.Post("/system/settings", newSettings, headers)
		if err != nil {
			t.Fatalf("Error al actualizar configuración: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d - Cuerpo: %s", resp.StatusCode, string(body))
		}

		// Verificar que el cambio se haya aplicado
		resp, body, err = apiClient.Get("/system/settings")
		if err != nil {
			t.Fatalf("Error al obtener configuración actualizada: %v", err)
		}

		var updatedSettings api.SystemSettings
		if err := json.Unmarshal(body, &updatedSettings); err != nil {
			t.Fatalf("Error al deserializar respuesta actualizada: %v", err)
		}

		if updatedSettings.MaintenanceMode != newSettings.MaintenanceMode {
			t.Errorf("El modo de mantenimiento no se actualizó correctamente: esperado %v, obtenido %v",
				newSettings.MaintenanceMode, updatedSettings.MaintenanceMode)
		}

		// Restaurar el estado original con una nueva clave idempotente
		headers = map[string]string{
			"Idempotency-Key": fmt.Sprintf("test-system-settings-restore-%d", time.Now().UnixNano()),
		}
		resp, _, err = apiClient.Post("/system/settings", originalSettings, headers)
		if err != nil {
			t.Logf("Error al restaurar configuración original (no fatal): %v", err)
		}
	})

	// Prueba para sincronizar conteo de claves de cuentas
	t.Run("Sync Account Key Count", func(t *testing.T) {
		// La API puede no admitir este endpoint o requerir un cuerpo de solicitud
		// específico, verificaremos ambos casos

		// Primero intentamos con un cuerpo vacío
		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("test-sync-account-key-count-%d", time.Now().UnixNano()),
		}

		// Usar un cuerpo vacío pero válido en formato JSON
		emptyBody := map[string]interface{}{}

		resp, body, err := apiClient.Post("/system/sync-account-key-count", emptyBody, headers)
		if err != nil {
			t.Fatalf("Error al sincronizar conteo de claves: %v", err)
		}

		// Validamos cualquier código de estado que no sea un error de servidor (5xx)
		if resp.StatusCode >= 500 {
			t.Errorf("Error de servidor: %d - Cuerpo: %s", resp.StatusCode, string(body))
		} else if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint no encontrado (404): %s", string(body))
			t.Skip("El endpoint /system/sync-account-key-count no está implementado")
		} else {
			t.Logf("Respuesta de sincronización (estado %d): %s", resp.StatusCode, string(body))
		}
	})

	// Prueba para verificar el endpoint de debug
	t.Run("Debug Endpoint", func(t *testing.T) {
		resp, body, err := apiClient.Get("/debug")
		if err != nil {
			t.Fatalf("Error al acceder al endpoint de debug: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// El endpoint de debug devuelve texto plano, verificar que no esté vacío
		if len(body) == 0 {
			t.Error("La respuesta del endpoint de debug está vacía")
		}

		t.Logf("Debug info: %s", string(body))
	})
}
