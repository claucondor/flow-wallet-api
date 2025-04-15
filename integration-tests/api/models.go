package api

import "time"

// SystemSettings representa la configuración del sistema
type SystemSettings struct {
	MaintenanceMode bool `json:"maintenanceMode"`
}

// HealthStatus representa el estado de salud del sistema
type HealthStatus struct {
	JobsInit        int `json:"jobsInit"`
	JobsNotAccepted int `json:"jobsNotAccepted"`
	JobsAccepted    int `json:"jobsAccepted"`
	JobsErrored     int `json:"jobsErrored"`
	JobsFailed      int `json:"jobsFailed"`
	JobsCompleted   int `json:"jobsCompleted"`
	PoolCapacity    int `json:"poolCapacity"`
	WorkerCount     int `json:"workerCount"`
}

// Account representa una cuenta en la API
type Account struct {
	Address   string    `json:"address"`
	Keys      []Key     `json:"keys"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Key representa una clave criptográfica
type Key struct {
	Index     int       `json:"index"`
	Type      string    `json:"type"`
	PublicKey string    `json:"publicKey"`
	SignAlgo  string    `json:"signAlgo"`
	HashAlgo  string    `json:"hashAlgo"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Job representa un trabajo asíncrono
type Job struct {
	JobID         string    `json:"jobId"`
	Type          string    `json:"type"`
	State         string    `json:"state"`
	Error         string    `json:"error"`
	Errors        []string  `json:"errors"`
	Result        string    `json:"result"`
	TransactionID string    `json:"transactionId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// Token representa un token habilitado
type Token struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Type    string `json:"type"`
}

// Script representa un script Cadence para ejecutar
type Script struct {
	Code      string           `json:"code"`
	Arguments []ScriptArgument `json:"arguments"`
}

// ScriptArgument representa un argumento para un script Cadence
type ScriptArgument struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Transaction representa una transacción
type Transaction struct {
	TransactionID   string    `json:"transactionId"`
	TransactionType string    `json:"transactionType"`
	Events          []Event   `json:"events,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Event representa un evento de transacción
type Event struct {
	Type             string `json:"Type"`
	TransactionID    string `json:"TransactionId"`
	TransactionIndex int    `json:"TransactionIndex"`
	EventIndex       int    `json:"EventIndex"`
	Value            string `json:"Value"`
}
