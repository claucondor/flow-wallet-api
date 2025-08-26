package app

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/flow-hydraulics/flow-wallet-api/configs"
	log "github.com/sirupsen/logrus"
)

// API represents the wallet API application
type API struct {
	config    *configs.Config
	sha1ver   string
	buildTime string
	server    *Server
}

// NewAPI creates a new API application instance
func NewAPI(cfg *configs.Config, sha1ver, buildTime string) *API {
	return &API{
		config:    cfg,
		sha1ver:   sha1ver,
		buildTime: buildTime,
	}
}

// Run starts the API application
func (a *API) Run() error {
	// Apply deployment mode configurations
	a.configureDeploymentModes()
	
	configs.ConfigureLogger(a.config.LogLevel)
	log.Info("Starting server")

	// Create and configure the server
	server, err := NewServer(a.config, a.sha1ver, a.buildTime)
	if err != nil {
		return err
	}
	a.server = server

	// Start the server
	if err := server.Start(); err != nil {
		return err
	}
	defer server.Stop()

	// Wait for shutdown signal
	return a.waitForShutdown()
}

// configureDeploymentModes applies configuration based on deployment mode
func (a *API) configureDeploymentModes() {
	if a.config.LightweightMode {
		log.Info("Running in lightweight mode with simplified dependencies")
		
		// Force SQLite as database
		a.config.DatabaseType = "sqlite"
		if a.config.DatabaseDSN == "" || a.config.DatabaseDSN == "wallet.db" {
			// Use a more explicit path for the SQLite database in lightweight mode
			a.config.DatabaseDSN = "./data/wallet-lightweight.db"
		}
		
		// Configure idempotency for lightweight mode
		if a.config.LightweightIdempotency {
			log.Info("Lightweight mode: Enabling idempotency with SQLite storage")
			a.config.DisableIdempotencyMiddleware = false
			a.config.IdempotencyMiddlewareDatabaseType = "shared" // Use same SQLite DB
		} else {
			log.Info("Lightweight mode: Disabling idempotency middleware")
			a.config.DisableIdempotencyMiddleware = true
		}
		
		// Optimize worker settings for lighter usage
		if a.config.WorkerCount > 4 {
			a.config.WorkerCount = 4
		}
		if a.config.WorkerQueueCapacity > 500 {
			a.config.WorkerQueueCapacity = 500
		}
	} else if a.config.QuaveMode {
		log.Info("Running in Quave Cloud mode with optimized settings")
		
		// Keep PostgreSQL but disable Redis-dependent features
		// Database configuration should be provided via environment variables
		
		// Disable idempotency middleware to remove Redis dependency
		a.config.DisableIdempotencyMiddleware = true
		
		// Optimize worker settings for cloud deployment
		if a.config.WorkerCount > 4 {
			a.config.WorkerCount = 4
		}
		if a.config.WorkerQueueCapacity > 500 {
			a.config.WorkerQueueCapacity = 500
		}
		
		log.Info("Quave mode: Redis disabled, PostgreSQL enabled, workers optimized")
	}
}

// waitForShutdown waits for an interrupt signal and gracefully shuts down the server
func (a *API) waitForShutdown() error {
	// Trap interrupt or sigterm and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	sig := <-c

	log.Infof("Got signal: %s. Shutting down..", sig)

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	return a.server.Shutdown(ctx)
}