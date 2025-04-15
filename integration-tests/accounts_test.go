package integrationtests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/api"
	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/utils"
)

// TestAccountsFlow prueba el flujo completo de cuentas: lista, crea, y obtiene detalles
func TestAccountsFlow(t *testing.T) {
	// Variable para almacenar la dirección de la cuenta creada
	var createdAccountAddress string

	t.Run("List Accounts", func(t *testing.T) {
		resp, body, err := apiClient.Get("/accounts")
		if err != nil {
			t.Fatalf("Error al listar cuentas: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var accounts []api.Account
		if err := json.Unmarshal(body, &accounts); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que al menos exista la cuenta de admin
		if len(accounts) < 1 {
			t.Errorf("Se esperaba al menos una cuenta, obtenido: %d", len(accounts))
		}
	})

	t.Run("Create Account Sync", func(t *testing.T) {
		// Crear una cuenta de forma síncrona
		headers := map[string]string{
			"Idempotency-Key": utils.GenerateIdempotencyKey(),
		}

		resp, body, err := apiClient.Post("/accounts?sync=true", nil, headers)
		if err != nil {
			t.Fatalf("Error al crear cuenta: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Estado esperado: 201, obtenido: %d", resp.StatusCode)
		}

		var account api.Account
		if err := json.Unmarshal(body, &account); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que la cuenta tenga una dirección válida
		if account.Address == "" {
			t.Error("La cuenta creada no tiene dirección")
		}

		// Verificar que la cuenta tenga al menos una clave
		if len(account.Keys) == 0 {
			t.Error("La cuenta creada no tiene claves")
		}

		// Guardar la dirección para pruebas posteriores
		createdAccountAddress = account.Address
		t.Logf("Cuenta creada con dirección: %s", createdAccountAddress)
	})

	t.Run("Get Account Details", func(t *testing.T) {
		// Verificar que la cuenta exista y tenga los detalles correctos
		if createdAccountAddress == "" {
			t.Skip("Prueba de creación de cuenta falló, saltando esta prueba")
		}

		resp, body, err := apiClient.Get("/accounts/" + createdAccountAddress)
		if err != nil {
			t.Fatalf("Error al obtener detalles de cuenta: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Estado esperado: 200, obtenido: %d", resp.StatusCode)
		}

		var account api.Account
		if err := json.Unmarshal(body, &account); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que sea la misma cuenta
		if account.Address != createdAccountAddress {
			t.Errorf("Dirección esperada: %s, obtenida: %s", createdAccountAddress, account.Address)
		}
	})

	t.Run("Create Account Async", func(t *testing.T) {
		// Crear una cuenta de forma asíncrona
		headers := map[string]string{
			"Idempotency-Key": utils.GenerateIdempotencyKey(),
		}

		resp, body, err := apiClient.Post("/accounts", nil, headers)
		if err != nil {
			t.Fatalf("Error al crear cuenta asíncrona: %v", err)
		}

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Estado esperado: 201, obtenido: %d", resp.StatusCode)
		}

		var job api.Job
		if err := json.Unmarshal(body, &job); err != nil {
			t.Fatalf("Error al deserializar respuesta: %v", err)
		}

		// Verificar que se haya creado un job
		if job.JobID == "" {
			t.Error("No se creó un job para la cuenta asíncrona")
		}

		// Esperar a que el job se complete
		t.Logf("Esperando a que el job %s se complete...", job.JobID)
		completedJob, err := waitForJobCompletion(job.JobID, 30*time.Second)
		if err != nil {
			t.Fatalf("Error esperando la finalización del job: %v", err)
		}

		// Verificar que el job se haya completado correctamente
		if completedJob.State != "COMPLETE" {
			t.Errorf("Estado del job esperado: COMPLETE, obtenido: %s", completedJob.State)
		}

		if completedJob.Result == "" {
			t.Error("El job se completó, pero no tiene resultado (dirección de cuenta)")
		}
	})
}

// waitForJobCompletion espera a que un job se complete o falle
func waitForJobCompletion(jobID string, timeout time.Duration) (*api.Job, error) {
	startTime := time.Now()

	for {
		resp, body, err := apiClient.Get("/jobs/" + jobID)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			continue
		}

		var job api.Job
		if err := json.Unmarshal(body, &job); err != nil {
			return nil, err
		}

		// El job está completo o ha fallado
		if job.State == "COMPLETE" || job.State == "FAILED" || job.State == "ERROR" {
			return &job, nil
		}

		// Verificar timeout
		if time.Since(startTime) > timeout {
			return &job, nil
		}

		// Esperar antes de verificar de nuevo
		time.Sleep(1 * time.Second)
	}
}
