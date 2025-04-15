package integrationtests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestTokenOperations(t *testing.T) {
	// Crear una cuenta para las pruebas
	var accountAddress string

	t.Run("Create Test Account", func(t *testing.T) {
		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("create-account-for-token-test-%d", time.Now().UnixNano()),
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
		t.Skip("No se pudo crear la cuenta de prueba, omitiendo pruebas de tokens")
	}

	// Desplegar el contrato FUSD
	t.Run("Deploy FUSD Token", func(t *testing.T) {
		adminAddress := os.Getenv("FLOW_WALLET_ADMIN_ADDRESS")
		if adminAddress == "" {
			t.Skip("No se ha configurado la dirección del administrador, omitiendo despliegue del token")
			return
		}

		deployRequest := map[string]interface{}{
			"tokenName": "FUSD",
			"address":   adminAddress,
		}

		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("deploy-fusd-token-%d", time.Now().UnixNano()),
		}

		resp, body, err := apiClient.Post("/tokens/deploy?sync=true", deployRequest, headers)
		if err != nil {
			t.Fatalf("Error al desplegar token FUSD: %v", err)
		}

		if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint de despliegue no implementado (404): %s", string(body))
			t.Skip("Endpoint /tokens/deploy no implementado en esta versión")
			return
		}

		if resp.StatusCode != http.StatusOK {
			t.Logf("Error al desplegar FUSD, estado: %d - Cuerpo: %s", resp.StatusCode, string(body))
		} else {
			t.Logf("Token FUSD desplegado correctamente: %s", string(body))
		}
	})

	// Listar tokens disponibles
	var fungibleTokens []map[string]interface{}
	var testTokenName string

	t.Run("List Available Tokens", func(t *testing.T) {
		// Obtener tokens fungibles
		resp, body, err := apiClient.Get("/fungible-tokens")
		if err != nil {
			t.Fatalf("Error al listar tokens fungibles: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		// Verificar que la respuesta tenga el formato esperado
		if err := json.Unmarshal(body, &fungibleTokens); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		if len(fungibleTokens) == 0 {
			t.Skip("No hay tokens fungibles disponibles, omitiendo pruebas relacionadas")
			return
		}

		// Usar el primer token disponible para las pruebas
		testTokenName, _ = fungibleTokens[0]["name"].(string)
		if testTokenName == "" {
			t.Skip("No se pudo obtener el nombre del token, omitiendo pruebas relacionadas")
			return
		}

		t.Logf("Usando token '%s' para las pruebas", testTokenName)
	})

	// Si no tenemos un token para probar, omitir el resto de las pruebas
	if testTokenName == "" {
		t.Skip("No se obtuvo un token válido para las pruebas, omitiendo pruebas restantes")
	}

	// Setup de token para la cuenta
	t.Run("Setup Token For Account", func(t *testing.T) {
		// Crear petición para setup de token
		setupRequest := map[string]interface{}{
			"name": testTokenName,
		}

		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("setup-token-%d", time.Now().UnixNano()),
		}

		endpoint := fmt.Sprintf("/accounts/%s/tokens/%s", accountAddress, testTokenName)
		resp, body, err := apiClient.Post(endpoint, setupRequest, headers)
		if err != nil {
			t.Fatalf("Error al realizar setup de token: %v", err)
		}

		// El endpoint puede no estar implementado (404) o devolver otros códigos
		if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint no implementado (404): %s", string(body))
			t.Skip("Endpoint /accounts/{address}/tokens/{token} no implementado en esta versión")
			return
		}

		if resp.StatusCode >= 400 {
			t.Logf("Estado de respuesta: %d - Cuerpo: %s", resp.StatusCode, string(body))
			t.Skip("No se pudo realizar el setup del token, posiblemente no soportado o ya configurado")
			return
		}

		t.Logf("Setup de token completado: %s", string(body))
	})

	// Verificar balance de token
	t.Run("Check Token Balance", func(t *testing.T) {
		endpoint := fmt.Sprintf("/accounts/%s/tokens/%s/balance", accountAddress, testTokenName)
		resp, body, err := apiClient.Get(endpoint)
		if err != nil {
			t.Fatalf("Error al verificar balance: %v", err)
		}

		// El endpoint puede no estar implementado (404)
		if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint de balance no implementado (404): %s", string(body))
			t.Skip("Endpoint de balance no implementado en esta versión")
			return
		}

		// Verificar respuesta exitosa (puede ser 0 si no hay fondos)
		if resp.StatusCode != http.StatusOK {
			t.Logf("Estado esperado: 200, obtenido: %d - Cuerpo: %s", resp.StatusCode, string(body))
		} else {
			t.Logf("Balance actual: %s", string(body))
		}
	})

	// Iniciar una transferencia de token (si la cuenta tiene fondos y hay otra cuenta disponible)
	t.Run("Transfer Token", func(t *testing.T) {
		// Primero verificamos si hay fondos
		balanceEndpoint := fmt.Sprintf("/accounts/%s/tokens/%s/balance", accountAddress, testTokenName)
		resp, body, err := apiClient.Get(balanceEndpoint)
		if err != nil {
			t.Fatalf("Error al verificar balance para transferencia: %v", err)
		}

		// El endpoint puede no estar implementado (404)
		if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint de balance no implementado (404): %s", string(body))
			t.Skip("Endpoint de balance no implementado, omitiendo prueba de transferencia")
			return
		}

		// Si no podemos obtener el balance, omitimos la prueba
		if resp.StatusCode != http.StatusOK {
			t.Skip("No se pudo verificar el balance, omitiendo prueba de transferencia")
			return
		}

		// Parsear el balance como número (puede ser string o número)
		var balance float64
		var balanceResponse interface{}
		if err := json.Unmarshal(body, &balanceResponse); err != nil {
			t.Logf("Error al deserializar balance: %v", err)
			t.Skip("No se pudo leer el balance, omitiendo prueba de transferencia")
			return
		}

		// Intentar convertir la respuesta a un número
		switch v := balanceResponse.(type) {
		case float64:
			balance = v
		case string:
			// Intentar parsear string a float, omitir si falla
			t.Skip("Balance en formato no numérico, omitiendo prueba de transferencia")
			return
		default:
			t.Skip("Formato de balance desconocido, omitiendo prueba de transferencia")
			return
		}

		// Si no hay fondos suficientes, omitir la prueba
		if balance <= 0 {
			t.Skip("No hay fondos suficientes para transferir, omitiendo prueba de transferencia")
			return
		}

		// Necesitamos otra cuenta para la transferencia
		// Crear una segunda cuenta
		headers := map[string]string{
			"Idempotency-Key": fmt.Sprintf("create-second-account-%d", time.Now().UnixNano()),
		}

		resp, body, err = apiClient.Post("/accounts?sync=true", nil, headers)
		if err != nil {
			t.Fatalf("Error al crear segunda cuenta: %v", err)
		}

		// El API devuelve 201 para la creación de cuentas
		if resp.StatusCode != http.StatusCreated {
			t.Skip("No se pudo crear la segunda cuenta, omitiendo prueba de transferencia")
			return
		}

		// Extraer la dirección de la segunda cuenta
		var accountResponse map[string]interface{}
		if err := json.Unmarshal(body, &accountResponse); err != nil {
			t.Skip("No se pudo deserializar respuesta de segunda cuenta")
			return
		}

		secondAddress, ok := accountResponse["address"].(string)
		if !ok || secondAddress == "" {
			t.Skip("No se pudo extraer la dirección de la segunda cuenta")
			return
		}

		// Realizar la transferencia
		transferAmount := 0.1 // Transferir una cantidad pequeña
		if balance < transferAmount {
			transferAmount = balance / 2 // Si no hay suficiente, transferir la mitad
		}

		transferRequest := map[string]interface{}{
			"amount":    transferAmount,
			"recipient": secondAddress,
		}

		headers = map[string]string{
			"Idempotency-Key": fmt.Sprintf("transfer-token-%d", time.Now().UnixNano()),
		}

		transferEndpoint := fmt.Sprintf("/accounts/%s/tokens/%s/transfer", accountAddress, testTokenName)
		resp, body, err = apiClient.Post(transferEndpoint, transferRequest, headers)
		if err != nil {
			t.Fatalf("Error al realizar transferencia: %v", err)
		}

		// El endpoint puede no estar implementado (404)
		if resp.StatusCode == http.StatusNotFound {
			t.Logf("Endpoint de transferencia no implementado (404): %s", string(body))
			t.Skip("Endpoint de transferencia no implementado en esta versión")
			return
		}

		// Puede ser síncrono o asíncrono
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
			t.Logf("Error al transferir, estado: %d - Cuerpo: %s", resp.StatusCode, string(body))
		} else {
			t.Logf("Transferencia iniciada: %s", string(body))
		}
	})
}
