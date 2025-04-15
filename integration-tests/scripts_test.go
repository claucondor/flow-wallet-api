package integrationtests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/api"
)

func TestScripts(t *testing.T) {
	t.Run("Execute Simple Script", func(t *testing.T) {
		// Crear un script simple que devuelve un entero
		script := api.Script{
			Code:      "pub fun main(): Int { return 42 }",
			Arguments: []api.ScriptArgument{},
		}

		resp, body, err := apiClient.Post("/scripts", script, nil)
		if err != nil {
			t.Fatalf("Error al ejecutar script: %v", err)
		}

		// Imprimir la respuesta para depuración
		t.Logf("Respuesta (código %d): %s", resp.StatusCode, string(body))

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// La respuesta debe ser 42
		var result int
		if err := json.Unmarshal(body, &result); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v - Cuerpo: %s", err, string(body))
		}

		if result != 42 {
			t.Errorf("Resultado esperado: 42, obtenido: %d", result)
		}
	})

	t.Run("Execute Script With Arguments", func(t *testing.T) {
		// Crear un script que acepta y devuelve un string
		script := api.Script{
			Code: "pub fun main(message: String): String { return message }",
			Arguments: []api.ScriptArgument{
				{
					Type:  "String",
					Value: "Hello, Flow!",
				},
			},
		}

		resp, body, err := apiClient.Post("/scripts", script, nil)
		if err != nil {
			t.Fatalf("Error al ejecutar script: %v", err)
		}

		// Imprimir la respuesta para depuración
		t.Logf("Respuesta (código %d): %s", resp.StatusCode, string(body))

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// La respuesta debe contener el mensaje original
		var cadenceValue struct {
			Value string `json:"Value"`
		}

		if err := json.Unmarshal(body, &cadenceValue); err != nil {
			// Puede que la respuesta sea un valor simple, intentar con string
			var stringValue string
			if err := json.Unmarshal(body, &stringValue); err != nil {
				t.Fatalf("Error al deserializar respuesta: %v - Cuerpo: %s", err, string(body))
			}

			if stringValue != "Hello, Flow!" {
				t.Errorf("Resultado esperado: 'Hello, Flow!', obtenido: '%s'", stringValue)
			}
		} else {
			if cadenceValue.Value != "Hello, Flow!" {
				t.Errorf("Resultado esperado: 'Hello, Flow!', obtenido: '%s'", cadenceValue.Value)
			}
		}
	})
}
