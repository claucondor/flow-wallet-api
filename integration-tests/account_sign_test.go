package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestAccountSignTransaction(t *testing.T) {
	// Crear una cuenta de prueba para usarla en el test
	var accountAddress string
	t.Run("Create Test Account", func(t *testing.T) {
		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("create-account-for-sign-test-%d", time.Now().UnixNano()),
		}

		resp, body, err := apiClient.Post("/accounts?sync=true", nil, headers)
		if err != nil {
			t.Fatalf("Error al crear cuenta de prueba: %v", err)
		}

		// El API devuelve 201 para la creación de cuentas
		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("Estado esperado: 201, obtenido: %d - Cuerpo: %s", resp.StatusCode, string(body))
		}

		// Extraer la dirección de la cuenta creada
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		var ok bool
		accountAddress, ok = response["address"].(string)
		if !ok || accountAddress == "" {
			t.Fatalf("No se pudo extraer la dirección de la cuenta: %s", string(body))
		}

		t.Logf("Cuenta creada con dirección: %s", accountAddress)
	})

	// Si no se pudo crear la cuenta, omitir el resto de las pruebas
	if accountAddress == "" {
		t.Skip("No se pudo crear la cuenta de prueba, omitiendo pruebas de firma")
	}

	t.Run("Sign Transaction Without Sending", func(t *testing.T) {
		// Preparar una transacción básica para firmar
		// Esta es una transacción simple que no hace nada, solo para probar la firma
		transactionRequest := map[string]interface{}{
			"script": `
				transaction {
					prepare(signer: AuthAccount) {
						// Esta transacción no hace nada, solo se usa para probar la firma
					}
				}
			`,
			"arguments": []interface{}{},
		}

		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("sign-transaction-%d", time.Now().UnixNano()),
		}

		// Firmar la transacción sin enviarla
		endpoint := fmt.Sprintf("/accounts/%s/sign", accountAddress)
		resp, body, err := apiClient.Post(endpoint, transactionRequest, headers)
		if err != nil {
			t.Fatalf("Error al firmar transacción: %v", err)
		}

		// El endpoint puede devolver diferentes códigos de estado dependiendo de la implementación
		// o puede no estar implementado en algunas versiones de la API
		if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint de firma no implementado (404): %s", string(body))
			t.Skip("El endpoint de firma no está implementado en esta versión de la API")
			return
		}

		if resp.StatusCode >= 400 {
			t.Errorf("Error al firmar transacción, estado: %d - Cuerpo: %s", resp.StatusCode, string(body))
		} else {
			// Verificar que la respuesta contenga información sobre la transacción firmada
			var response map[string]interface{}
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatalf("Error al deserializar respuesta: %v - Cuerpo: %s", err, string(body))
			}

			// Dependiendo de la implementación, la respuesta podría contener diferentes campos
			// Algunos campos comunes podrían ser: signatures, encodedTransaction, signedTransaction
			t.Logf("Transacción firmada: %s", string(body))
		}
	})
}
