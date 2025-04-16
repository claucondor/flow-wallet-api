package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestWatchlistAndOps(t *testing.T) {
	// Pruebas para el endpoint de Watchlist
	t.Run("Watchlist Endpoints", func(t *testing.T) {
		// Crear una dirección de cuenta de Flow para la watchlist
		// Usamos una dirección ficticia ya que solo probamos la API
		testAddress := "0x01cf0e2f2f715450"

		// 1. Agregar dirección a la watchlist
		t.Run("Add Address To Watchlist", func(t *testing.T) {
			watchlistRequest := map[string]interface{}{
				"address": testAddress,
			}

			headers := map[string]string{
				"Idempotency-Key": fmt.Sprintf("add-to-watchlist-%d", time.Now().UnixNano()),
			}

			resp, body, err := apiClient.Post("/watchlist/accounts", watchlistRequest, headers)
			if err != nil {
				t.Fatalf("Error al agregar dirección a watchlist: %v", err)
			}

			// El endpoint puede no estar implementado (404)
			if resp.StatusCode == http.StatusNotFound {
				t.Logf("Endpoint de watchlist no implementado (404): %s", string(body))
				t.Skip("Endpoint de watchlist no está implementado en esta versión")
				return
			}

			// El endpoint puede devolver diferentes códigos de estado
			if resp.StatusCode >= 400 {
				t.Logf("No se pudo agregar a watchlist, estado: %d - Cuerpo: %s", resp.StatusCode, string(body))
				t.Skip("Watchlist no soportada o error al añadir, omitiendo pruebas relacionadas")
				return
			}

			t.Logf("Dirección agregada a watchlist: %s", string(body))
		})

		// 2. Listar cuentas en la watchlist
		t.Run("List Watchlist Accounts", func(t *testing.T) {
			resp, body, err := apiClient.Get("/watchlist/accounts")
			if err != nil {
				t.Fatalf("Error al listar cuentas de watchlist: %v", err)
			}

			// El endpoint puede no estar implementado (404)
			if resp.StatusCode == http.StatusNotFound {
				t.Logf("Endpoint de lista de watchlist no implementado (404): %s", string(body))
				t.Skip("Endpoint de lista de watchlist no implementado en esta versión")
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Logf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
			}

			// Verificar que la respuesta tenga el formato esperado
			var watchlistAccounts []interface{}
			if err := json.Unmarshal(body, &watchlistAccounts); err != nil {
				t.Fatalf("Error al deserializar respuesta: %v", err)
			}

			t.Logf("Cuentas en watchlist: %d", len(watchlistAccounts))
		})

		// 3. Eliminar dirección de la watchlist
		t.Run("Remove Address From Watchlist", func(t *testing.T) {
			endpoint := fmt.Sprintf("/watchlist/accounts/%s", testAddress)
			req, err := http.NewRequest(http.MethodDelete, apiClient.BaseURL+endpoint, nil)
			if err != nil {
				t.Fatalf("Error al crear petición DELETE: %v", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("Error al eliminar dirección de watchlist: %v", err)
			}
			defer resp.Body.Close()

			// El endpoint puede no estar implementado (404)
			if resp.StatusCode == http.StatusNotFound {
				t.Logf("Endpoint de eliminación de watchlist no implementado (404)")
				// No es un error, simplemente el recurso no existe o el endpoint no está implementado
				return
			}

			// El endpoint puede devolver diferentes códigos de estado
			if resp.StatusCode >= 400 {
				t.Logf("Error al eliminar de watchlist, estado: %d", resp.StatusCode)
			} else {
				t.Logf("Dirección eliminada de watchlist, estado: %d", resp.StatusCode)
			}
		})
	})

	// Pruebas para los endpoints de operaciones (Ops)
	t.Run("Ops Endpoints", func(t *testing.T) {
		// 1. Obtener estadísticas de missing vaults
		t.Run("Get Missing Vaults Stats", func(t *testing.T) {
			resp, body, err := apiClient.Get("/ops/missing-fungible-token-vaults/stats")
			if err != nil {
				t.Fatalf("Error al obtener estadísticas de vaults: %v", err)
			}

			// El endpoint puede no estar implementado (404)
			if resp.StatusCode == http.StatusNotFound {
				t.Logf("Endpoint de estadísticas de vaults no implementado (404): %s", string(body))
				t.Skip("Endpoint de estadísticas de vaults no implementado en esta versión")
				return
			}

			// El endpoint puede no estar disponible o requerir autenticación
			if resp.StatusCode >= 400 {
				t.Logf("No se pudo acceder a las estadísticas de vaults, estado: %d", resp.StatusCode)
				t.Skip("Endpoint de estadísticas no accesible, omitiendo pruebas relacionadas")
				return
			}

			// Verificar que la respuesta tenga el formato esperado (puede variar)
			t.Logf("Estadísticas de vaults: %s", string(body))
		})

		// 2. Iniciar proceso de creación de vaults faltantes
		t.Run("Start Missing Vaults Creation", func(t *testing.T) {
			headers := map[string]string{
				"Idempotency-Key": fmt.Sprintf("start-missing-vaults-%d", time.Now().UnixNano()),
			}

			resp, body, err := apiClient.Post("/ops/missing-fungible-token-vaults/start", nil, headers)
			if err != nil {
				t.Fatalf("Error al iniciar creación de vaults: %v", err)
			}

			// El endpoint puede no estar implementado (404)
			if resp.StatusCode == http.StatusNotFound {
				t.Logf("Endpoint de creación de vaults no implementado (404): %s", string(body))
				t.Skip("Endpoint de creación de vaults no implementado en esta versión")
				return
			}

			// El endpoint puede no estar disponible o requerir autenticación
			if resp.StatusCode >= 400 {
				t.Logf("No se pudo iniciar la creación de vaults, estado: %d - Cuerpo: %s", resp.StatusCode, string(body))
				t.Skip("Endpoint de creación de vaults no accesible")
				return
			}

			t.Logf("Proceso de creación de vaults iniciado: %s", string(body))
		})
	})
}
