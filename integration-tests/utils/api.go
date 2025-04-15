package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// CheckAPIReadiness verifica si la API está disponible
func CheckAPIReadiness(url string, timeout time.Duration) error {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	startTime := time.Now()
	for {
		resp, err := client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}

		if resp != nil {
			resp.Body.Close()
		}

		if time.Since(startTime) > timeout {
			return fmt.Errorf("timeout esperando a que la API esté disponible")
		}

		time.Sleep(1 * time.Second)
	}
}

// APIClient es un cliente HTTP para interactuar con la API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPIClient crea un nuevo cliente API
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Get realiza una petición GET
func (c *APIClient) Get(path string) (*http.Response, []byte, error) {
	return c.Request("GET", path, nil, nil)
}

// Post realiza una petición POST
func (c *APIClient) Post(path string, body interface{}, headers map[string]string) (*http.Response, []byte, error) {
	return c.Request("POST", path, body, headers)
}

// Request realiza una petición HTTP genérica
func (c *APIClient) Request(method, path string, body interface{}, headers map[string]string) (*http.Response, []byte, error) {
	url := c.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, nil, err
	}

	// Establecer headers predeterminados
	req.Header.Set("Content-Type", "application/json")

	// Agregar headers personalizados
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, respBody, nil
}

// GenerateIdempotencyKey genera una clave de idempotencia única
func GenerateIdempotencyKey() string {
	return fmt.Sprintf("test-%d", time.Now().UnixNano())
}
