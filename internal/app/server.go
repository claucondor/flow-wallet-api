package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	flowgorm "github.com/flow-hydraulics/flow-wallet-api/internal/datastore/gorm"
	access "github.com/onflow/flow-go-sdk/access/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/flow-hydraulics/flow-wallet-api/internal/accounts"
	"github.com/flow-hydraulics/flow-wallet-api/internal/chain_events"
	"github.com/flow-hydraulics/flow-wallet-api/configs"
	"github.com/flow-hydraulics/flow-wallet-api/internal/handlers"
	"github.com/flow-hydraulics/flow-wallet-api/internal/jobs"
	"github.com/flow-hydraulics/flow-wallet-api/internal/keys"
	"github.com/flow-hydraulics/flow-wallet-api/internal/keys/basic"
	"github.com/flow-hydraulics/flow-wallet-api/internal/ops"
	"github.com/flow-hydraulics/flow-wallet-api/internal/system"
	"github.com/flow-hydraulics/flow-wallet-api/internal/templates"
	"github.com/flow-hydraulics/flow-wallet-api/internal/tokens"
	"github.com/flow-hydraulics/flow-wallet-api/internal/transactions"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"go.uber.org/ratelimit"
	"gorm.io/gorm"
)

// Server manages the HTTP server and all its dependencies
type Server struct {
	config    *configs.Config
	sha1ver   string
	buildTime string
	
	// Dependencies
	db     *gorm.DB
	fc     *access.Client
	wp     jobs.WorkerPool
	listener chain_events.Listener
	
	// Services
	systemService      system.Service
	templateService    templates.Service
	jobsService        jobs.Service
	transactionService transactions.Service
	accountService     accounts.Service
	tokenService       tokens.Service
	opsService         ops.Service
	
	// HTTP Server
	httpServer *http.Server
}

// NewServer creates a new server instance
func NewServer(cfg *configs.Config, sha1ver, buildTime string) (*Server, error) {
	s := &Server{
		config:    cfg,
		sha1ver:   sha1ver,
		buildTime: buildTime,
	}
	
	if err := s.initDependencies(); err != nil {
		return nil, err
	}
	
	if err := s.initServices(); err != nil {
		return nil, err
	}
	
	s.initHTTPServer()
	
	return s, nil
}

// initDependencies initializes core dependencies
func (s *Server) initDependencies() error {
	// Flow client
	fc, err := access.NewClient(
		s.config.AccessAPIHost,
		access.WithGRPCDialOptions(
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(s.config.GrpcMaxCallRecvMsgSize)),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create Flow client: %w", err)
	}
	s.fc = fc

	// Database
	db, err := flowgorm.New(s.config)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}
	s.db = db

	// System service (needed for worker pool)
	s.systemService = system.NewService(
		system.NewGormStore(db),
		system.WithPauseDuration(s.config.PauseDuration),
	)

	// Create worker pool
	s.wp = jobs.NewWorkerPool(
		jobs.NewGormStore(db),
		s.config.WorkerQueueCapacity,
		s.config.WorkerCount,
		jobs.WithJobStatusWebhook(s.config.JobStatusWebhookUrl, s.config.JobStatusWebhookTimeout),
		jobs.WithSystemService(s.systemService),
		jobs.WithMaxJobErrorCount(s.config.MaxJobErrorCount),
		jobs.WithDbJobPollInterval(s.config.DBJobPollInterval),
		jobs.WithAcceptedGracePeriod(s.config.AcceptedGracePeriod),
		jobs.WithReSchedulableGracePeriod(s.config.ReSchedulableGracePeriod),
	)

	return nil
}

// initServices initializes business services
func (s *Server) initServices() error {
	txRatelimiter := ratelimit.New(s.config.TransactionMaxSendRate, ratelimit.WithoutSlack)

	// Key manager
	km := basic.NewKeyManager(s.config, keys.NewGormStore(s.db), s.fc)

	// Template service
	templateService, err := templates.NewService(s.config, templates.NewGormStore(s.db))
	if err != nil {
		return fmt.Errorf("failed to create template service: %w", err)
	}
	s.templateService = templateService

	// Other services
	s.jobsService = jobs.NewService(jobs.NewGormStore(s.db))
	s.transactionService = transactions.NewService(s.config, transactions.NewGormStore(s.db), km, s.fc, s.wp, transactions.WithTxRatelimiter(txRatelimiter))
	s.accountService = accounts.NewService(s.config, accounts.NewGormStore(s.db), km, s.fc, s.wp, s.transactionService, s.templateService, accounts.WithTxRatelimiter(txRatelimiter))
	s.tokenService = tokens.NewService(s.config, tokens.NewGormStore(s.db), km, s.fc, s.wp, s.transactionService, s.templateService, s.accountService)
	s.opsService = ops.NewService(s.config, ops.NewGormStore(s.db), s.templateService, s.transactionService, s.tokenService)

	// Register event handlers
	accounts.AccountAdded.Register(&tokens.AccountAddedHandler{
		TemplateService: s.templateService,
		TokenService:    s.tokenService,
	})

	// Initialize admin account
	err = s.accountService.InitAdminAccount(context.Background())
	if err != nil {
		return fmt.Errorf("failed to initialize admin account: %w", err)
	}

	return nil
}

// initHTTPServer sets up the HTTP server and routes
func (s *Server) initHTTPServer() {
	// HTTP handlers
	systemHandler := handlers.NewSystem(s.systemService)
	templateHandler := handlers.NewTemplates(s.templateService)
	jobsHandler := handlers.NewJobs(s.jobsService)
	accountHandler := handlers.NewAccounts(s.accountService)
	transactionHandler := handlers.NewTransactions(s.transactionService)
	tokenHandler := handlers.NewTokens(s.tokenService)
	opsHandler := handlers.NewOps(s.opsService)

	r := mux.NewRouter()

	// Catch the api version
	rv := r.PathPrefix("/{apiVersion}").Subrouter()

	// Debug
	rv.Handle("/debug", handlers.Debug("https://github.com/flow-hydraulics/flow-wallet-api", s.sha1ver, s.buildTime)).Methods(http.MethodGet)

	// Swagger UI and OpenAPI spec
	rv.HandleFunc("/docs", handlers.HandleSwaggerUI).Methods(http.MethodGet)
	rv.HandleFunc("/openapi.yml", handlers.HandleOpenAPISpec(handlers.OpenAPISpec)).Methods(http.MethodGet)

	// Health
	rv.HandleFunc("/health/ready", handlers.HandleHealthReady).Methods(http.MethodGet)
	rv.Handle("/health/liveness", handlers.Liveness(func() (interface{}, error) {
		return s.wp.Status()
	})).Methods(http.MethodGet)

	// System
	rv.Handle("/system/settings", systemHandler.GetSettings()).Methods(http.MethodGet)
	rv.Handle("/system/settings", systemHandler.SetSettings()).Methods(http.MethodPost)
	rv.Handle("/system/sync-account-key-count", accountHandler.SyncAccountKeyCount()).Methods(http.MethodPost)

	// Jobs
	rv.Handle("/jobs", jobsHandler.List()).Methods(http.MethodGet)
	rv.Handle("/jobs/{jobId}", jobsHandler.Details()).Methods(http.MethodGet)

	// Token templates
	rv.Handle("/tokens", templateHandler.ListTokens(templates.NotSpecified)).Methods(http.MethodGet)
	rv.Handle("/tokens", templateHandler.AddToken()).Methods(http.MethodPost)
	rv.Handle("/tokens/{id_or_name}", templateHandler.GetToken()).Methods(http.MethodGet)
	rv.Handle("/tokens/{id}", templateHandler.RemoveToken()).Methods(http.MethodDelete)

	// List enabled tokens by type
	rv.Handle("/fungible-tokens", templateHandler.ListTokens(templates.FT)).Methods(http.MethodGet)
	rv.Handle("/non-fungible-tokens", templateHandler.ListTokens(templates.NFT)).Methods(http.MethodGet)

	// Transactions
	rv.Handle("/transactions", transactionHandler.List()).Methods(http.MethodGet)
	rv.Handle("/transactions/{transactionId}", transactionHandler.Details()).Methods(http.MethodGet)

	// Account
	rv.Handle("/accounts", accountHandler.List()).Methods(http.MethodGet)
	rv.Handle("/accounts", accountHandler.Create()).Methods(http.MethodPost)
	rv.Handle("/accounts/{address}", accountHandler.Details()).Methods(http.MethodGet)

	// Account raw transactions
	if !s.config.DisableRawTransactions {
		rv.Handle("/accounts/{address}/sign", transactionHandler.Sign()).Methods(http.MethodPost)
		rv.Handle("/accounts/{address}/transactions", transactionHandler.List()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/transactions", transactionHandler.Create()).Methods(http.MethodPost)
		rv.Handle("/accounts/{address}/transactions/{transactionId}", transactionHandler.Details()).Methods(http.MethodGet)
	} else {
		log.Info("raw transactions disabled")
	}

	// Non-custodial watchlist accounts
	rv.Handle("/watchlist/accounts", accountHandler.AddNonCustodialAccount()).Methods(http.MethodPost)
	rv.Handle("/watchlist/accounts/{address}", accountHandler.DeleteNonCustodialAccount()).Methods(http.MethodDelete)

	// Non-custodial transaction support
	rv.Handle("/accounts/{address}/transactions/prepare", transactionHandler.PrepareTransaction()).Methods(http.MethodPost)
	rv.Handle("/transactions/submit", transactionHandler.SubmitTransaction()).Methods(http.MethodPost)

	// Scripts
	rv.Handle("/scripts", transactionHandler.ExecuteScript()).Methods(http.MethodPost)

	// Fungible tokens
	if !s.config.DisableFungibleTokens {
		rv.Handle("/accounts/{address}/fungible-tokens", tokenHandler.AccountTokens(templates.FT)).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}", tokenHandler.Details()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}", tokenHandler.Setup()).Methods(http.MethodPost)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}/withdrawals", tokenHandler.ListWithdrawals()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}/withdrawals", tokenHandler.CreateWithdrawal()).Methods(http.MethodPost)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}/withdrawals/{transactionId}", tokenHandler.GetWithdrawal()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}/deposits", tokenHandler.ListDeposits()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/fungible-tokens/{tokenName}/deposits/{transactionId}", tokenHandler.GetDeposit()).Methods(http.MethodGet)
	} else {
		log.Info("fungible tokens disabled")
	}

	// Non-Fungible tokens
	if !s.config.DisableNonFungibleTokens {
		rv.Handle("/accounts/{address}/non-fungible-tokens", tokenHandler.AccountTokens(templates.NFT)).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}", tokenHandler.Details()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}", tokenHandler.Setup()).Methods(http.MethodPost)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}/withdrawals", tokenHandler.ListWithdrawals()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}/withdrawals", tokenHandler.CreateWithdrawal()).Methods(http.MethodPost)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}/withdrawals/{transactionId}", tokenHandler.GetWithdrawal()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}/deposits", tokenHandler.ListDeposits()).Methods(http.MethodGet)
		rv.Handle("/accounts/{address}/non-fungible-tokens/{tokenName}/deposits/{transactionId}", tokenHandler.GetDeposit()).Methods(http.MethodGet)
	} else {
		log.Info("non-fungible tokens disabled")
	}

	// Ops
	rv.Handle("/ops/missing-fungible-token-vaults/start", opsHandler.InitMissingFungibleVaults()).Methods(http.MethodGet)
	rv.Handle("/ops/missing-fungible-token-vaults/stats", opsHandler.GetMissingFungibleVaults()).Methods(http.MethodGet)

	// Apply middleware
	h := http.TimeoutHandler(r, s.config.ServerRequestTimeout, "request timed out")
	h = handlers.UseCors(h)
	h = handlers.UseLogging(h)
	h = handlers.UseCompress(h)

	// Setup idempotency middleware if enabled
	if !s.config.DisableIdempotencyMiddleware {
		var is handlers.IdempotencyStore
		switch s.config.IdempotencyMiddlewareDatabaseType {
		case handlers.IdempotencyStoreTypeShared.String():
			is = handlers.NewIdempotencyStoreGorm(s.db)
		case handlers.IdempotencyStoreTypeRedis.String():
			if s.config.IdempotencyMiddlewareRedisURL == "" {
				log.Fatal("idempotency middleware db set to redis but Redis URL is empty")
			}
			pool := &redis.Pool{
				MaxIdle:   80,
				MaxActive: 12000,
				Dial: func() (redis.Conn, error) {
					c, err := redis.DialURL(s.config.IdempotencyMiddlewareRedisURL)
					if err != nil {
						panic(err.Error())
					}
					return c, err
				},
			}
			client := pool.Get()
			is = handlers.NewIdempotencyStoreRedis(client)
		case handlers.IdempotencyStoreTypeLocal.String():
			is = handlers.NewIdempotencyStoreLocal()
		}

		h = handlers.UseIdempotency(h, handlers.IdempotencyHandlerOptions{
			Expiry:      1 * time.Hour,
			IgnorePaths: []string{"/v1/scripts"}, // Scripts are read-only
		}, is)
	}

	// Create HTTP server
	s.httpServer = &http.Server{
		Handler:      h,
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		WriteTimeout: 0, // Disabled, set cfg.ServerRequestTimeout instead
		ReadTimeout:  0, // Disabled, set cfg.ServerRequestTimeout instead
	}
}

// Start starts all server components
func (s *Server) Start() error {
	// Start worker pool
	s.wp.Start()
	log.Info("Started workerpool")

	// Start chain event listener if enabled
	if !s.config.DisableChainEvents {
		if err := s.startChainListener(); err != nil {
			return fmt.Errorf("failed to start chain listener: %w", err)
		}
	}

	// Start HTTP server
	go func() {
		log.WithFields(log.Fields{
			"host": s.config.Host,
			"port": s.config.Port,
		}).Info("Server listening")
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Warn(err)
		}
	}()

	return nil
}

// startChainListener starts the blockchain event listener
func (s *Server) startChainListener() error {
	store := chain_events.NewGormStore(s.db)
	getTypes := func() ([]string, error) {
		// Get all enabled tokens
		tt, err := s.templateService.ListTokens(templates.NotSpecified)
		if err != nil {
			return nil, err
		}

		token_count := len(tt)
		event_types := make([]string, token_count)

		// Listen for enabled tokens deposit events
		for i, token := range tt {
			event_types[i] = templates.DepositEventTypeFromToken(token)
		}

		return event_types, nil
	}

	s.listener = chain_events.NewListener(
		s.fc, store, getTypes,
		s.config.ChainListenerMaxBlocks,
		s.config.ChainListenerInterval,
		s.config.ChainListenerStartingHeight,
		chain_events.WithSystemService(s.systemService),
	)

	// Register chain event handler
	chain_events.ChainEvent.Register(&tokens.ChainEventHandler{
		AccountService:  s.accountService,
		ChainListener:   s.listener,
		TemplateService: s.templateService,
		TokenService:    s.tokenService,
	})

	s.listener.Start()
	log.Info("Started chain events listener")

	return nil
}

// Stop gracefully stops all server components
func (s *Server) Stop() {
	// Stop chain listener
	if s.listener != nil {
		s.listener.Stop()
		log.Info("Stopped chain events listener")
	}

	// Stop worker pool
	if s.wp != nil {
		s.wp.Stop(true)
		log.Info("Stopped workerpool")
	}

	// Close database
	if s.db != nil {
		flowgorm.Close(s.db)
	}

	// Close Flow client
	if s.fc != nil {
		if err := s.fc.Close(); err != nil {
			log.Warn(err)
		}
		log.Info("Closed Flow Client")
	}
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Warnf("Error in server shutdown: %s", err)
		return err
	}
	return nil
}