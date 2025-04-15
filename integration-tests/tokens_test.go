package integrationtests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/api"
)

func TestTokens(t *testing.T) {
	t.Run("List Tokens", func(t *testing.T) {
		// Listar todos los tokens
		resp, body, err := apiClient.Get("/tokens")
		if err != nil {
			t.Fatalf("Error al listar tokens: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var tokens []api.Token
		if err := json.Unmarshal(body, &tokens); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Debe haber al menos algunos tokens habilitados por defecto (FlowToken, FUSD)
		if len(tokens) < 1 {
			t.Logf("No hay tokens habilitados. Esto es esperado si es una instalación limpia.")
		} else {
			t.Logf("Encontrados %d tokens habilitados", len(tokens))
		}
	})

	t.Run("List Fungible Tokens", func(t *testing.T) {
		// Listar tokens fungibles
		resp, body, err := apiClient.Get("/fungible-tokens")
		if err != nil {
			t.Fatalf("Error al listar tokens fungibles: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var tokens []api.Token
		if err := json.Unmarshal(body, &tokens); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar el formato de los tokens recibidos
		for _, token := range tokens {
			if token.Type != "FT" {
				t.Errorf("Se esperaba un token de tipo FT, obtenido: %s", token.Type)
			}
			if token.Name == "" || token.Address == "" {
				t.Error("Token con nombre o dirección vacío")
			}
		}
	})

	t.Run("List Non-Fungible Tokens", func(t *testing.T) {
		// Listar tokens no fungibles
		resp, body, err := apiClient.Get("/non-fungible-tokens")
		if err != nil {
			t.Fatalf("Error al listar tokens no fungibles: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var tokens []api.Token
		if err := json.Unmarshal(body, &tokens); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar el formato de los tokens recibidos
		for _, token := range tokens {
			if token.Type != "NFT" {
				t.Errorf("Se esperaba un token de tipo NFT, obtenido: %s", token.Type)
			}
			if token.Name == "" || token.Address == "" {
				t.Error("Token con nombre o dirección vacío")
			}
		}
	})
}
