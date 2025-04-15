package integrationtests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/flow-hydraulics/flow-wallet-api/integration-tests/utils"
)

const (
	apiBaseURL = "http://localhost:3001/v1"
	timeout    = 30 * time.Second
)

var (
	adminAddress string
	apiClient    *utils.APIClient
)

func TestMain(m *testing.M) {
	// Leer variables de entorno
	adminAddress = os.Getenv("FLOW_WALLET_ADMIN_ADDRESS")
	if adminAddress == "" {
		fmt.Println("ERROR: La variable de entorno FLOW_WALLET_ADMIN_ADDRESS es requerida")
		os.Exit(1)
	}

	// Verificar que la API esté disponible antes de comenzar las pruebas
	fmt.Println("Verificando que la API esté disponible...")
	err := waitForAPIReadiness(apiBaseURL, timeout)
	if err != nil {
		fmt.Printf("ERROR: La API no está disponible: %s\n", err)
		fmt.Println("Asegúrate de ejecutar 'make -f test.mk test' antes de las pruebas")
		os.Exit(1)
	}

	// Crear cliente API para todas las pruebas
	apiClient = utils.NewAPIClient(apiBaseURL)

	fmt.Println("Iniciando pruebas de integración...")
	exitCode := m.Run()
	fmt.Println("Pruebas completadas")

	os.Exit(exitCode)
}

// Función para verificar que la API esté disponible
func waitForAPIReadiness(baseURL string, timeout time.Duration) error {
	return utils.CheckAPIReadiness(baseURL+"/health/ready", timeout)
}
