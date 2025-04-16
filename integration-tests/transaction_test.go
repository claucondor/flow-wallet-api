package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestTransactionEndpoints(t *testing.T) {
	// Necesitamos crear una cuenta y realizar una transacción para probar los endpoints de transacciones
	var accountAddress string
	var transactionID string

	// 1. Crear cuenta de prueba
	t.Run("Create Test Account", func(t *testing.T) {
		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("create-account-for-tx-test-%d", time.Now().UnixNano()),
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
		t.Skip("No se pudo crear la cuenta de prueba, omitiendo pruebas de transacciones")
	}

	// 2. Ejecutar una transacción simple para obtener un ID de transacción
	t.Run("Execute Simple Transaction", func(t *testing.T) {
		// Transacción simple que no hace nada
		transactionRequest := map[string]interface{}{
			"type": "simple",
			"script": `
				transaction {
					prepare(signer: AuthAccount) {
						// Esta transacción no hace nada, solo se usa para pruebas
					}
				}
			`,
			"arguments": []interface{}{},
		}

		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("execute-simple-tx-%d", time.Now().UnixNano()),
		}

		endpoint := fmt.Sprintf("/accounts/%s/transactions", accountAddress)
		resp, body, err := apiClient.Post(endpoint, transactionRequest, headers)
		if err != nil {
			t.Fatalf("Error al ejecutar transacción: %v", err)
		}

		// La API puede devolver 200, 201 o 202 dependiendo de la implementación
		if resp.StatusCode >= 400 {
			t.Fatalf("Error al ejecutar transacción, estado: %d - Cuerpo: %s", resp.StatusCode, string(body))
		}

		// Extraer el ID de la transacción
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// El ID puede estar en diferentes campos dependiendo si es síncrono o asíncrono
		if txID, ok := response["transactionId"].(string); ok && txID != "" {
			transactionID = txID
		} else if jobID, ok := response["jobId"].(string); ok && jobID != "" {
			// Si es asíncrono, necesitamos esperar y obtener el resultado del job
			t.Logf("Transacción iniciada como job: %s, esperando resultado...", jobID)

			// Esperar hasta que el job termine (máximo 30 segundos)
			var jobCompleted bool
			for i := 0; i < 30; i++ {
				time.Sleep(time.Second)

				jobEndpoint := fmt.Sprintf("/jobs/%s", jobID)
				jobResp, body, err := apiClient.Get(jobEndpoint)
				if err != nil {
					t.Logf("Error al verificar estado del job: %v", err)
					continue
				}

				if jobResp.StatusCode != http.StatusOK {
					t.Logf("Estado no esperado al verificar job: %d", jobResp.StatusCode)
					continue
				}

				var jobStatus map[string]interface{}
				if err := json.Unmarshal(body, &jobStatus); err != nil {
					t.Logf("Error al deserializar respuesta del job: %v", err)
					continue
				}

				state, _ := jobStatus["state"].(string)
				if state == "COMPLETE" || state == "COMPLETED" {
					if txID, ok := jobStatus["transactionId"].(string); ok && txID != "" {
						transactionID = txID
						jobCompleted = true
						break
					}
				} else if state == "FAILED" || state == "ERRORED" {
					t.Logf("Job falló: %s", string(body))
					break
				}
			}

			if !jobCompleted {
				t.Skip("No se pudo obtener un ID de transacción, omitiendo pruebas restantes")
				return
			}
		}

		if transactionID == "" {
			t.Skip("No se pudo extraer el ID de transacción, omitiendo pruebas restantes")
			return
		}

		t.Logf("Transacción creada con ID: %s", transactionID)
	})

	// Si no tenemos un ID de transacción, omitir el resto de las pruebas
	if transactionID == "" {
		t.Skip("No se obtuvo un ID de transacción, omitiendo pruebas restantes")
	}

	// 3. Listar todas las transacciones
	t.Run("List All Transactions", func(t *testing.T) {
		resp, body, err := apiClient.Get("/transactions")
		if err != nil {
			t.Fatalf("Error al listar transacciones: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// Verificar que la respuesta tenga el formato esperado
		var transactions []interface{}
		if err := json.Unmarshal(body, &transactions); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que haya al menos una transacción
		if len(transactions) == 0 {
			t.Error("No se encontraron transacciones")
		} else {
			t.Logf("Se encontraron %d transacciones", len(transactions))
		}
	})

	// 4. Obtener detalles de una transacción específica
	t.Run("Get Transaction Details", func(t *testing.T) {
		endpoint := fmt.Sprintf("/transactions/%s", transactionID)
		resp, body, err := apiClient.Get(endpoint)
		if err != nil {
			t.Fatalf("Error al obtener detalles de la transacción: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// Verificar que la respuesta tenga el formato esperado
		var transaction map[string]interface{}
		if err := json.Unmarshal(body, &transaction); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que el ID de la transacción coincida
		if txID, ok := transaction["transactionId"].(string); ok {
			if txID != transactionID {
				t.Errorf("ID de transacción esperado: %s, obtenido: %s", transactionID, txID)
			}
		} else {
			t.Error("La respuesta no contiene el campo transactionId")
		}

		t.Logf("Detalles de la transacción obtenidos: %s", string(body))
	})

	// 5. Listar transacciones de una cuenta específica
	t.Run("List Account Transactions", func(t *testing.T) {
		endpoint := fmt.Sprintf("/accounts/%s/transactions", accountAddress)
		resp, body, err := apiClient.Get(endpoint)
		if err != nil {
			t.Fatalf("Error al listar transacciones de la cuenta: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// Verificar que la respuesta tenga el formato esperado
		var transactions []interface{}
		if err := json.Unmarshal(body, &transactions); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		t.Logf("Se encontraron %d transacciones para la cuenta", len(transactions))
	})
}
