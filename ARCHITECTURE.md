# Project Architecture

This project now follows the standard Go project layout with `cmd/` and `internal/` directories.

## Structure

```
cmd/
  api/              # Current wallet API (v1)
    main.go
  api-v2/           # Future API version
    main.go
  admin/            # Admin tools and utilities
    main.go

internal/          # Private application code
  app/             # Application setup and server
    api.go         # API application bootstrap
    server.go      # HTTP server and dependency setup
  accounts/        # Account management
  chain_events/    # Blockchain event listening
  datastore/       # Database abstraction
  errors/          # Error types
  handlers/        # HTTP handlers
  jobs/            # Background job processing
  keys/            # Key management (local, AWS KMS, Google KMS)
  migrations/      # Database migrations
  ops/             # Operations and maintenance
  system/          # System-level services
  templates/       # Transaction templates
  tokens/          # Token management (FT/NFT)
  transactions/    # Transaction handling

configs/           # Configuration management (shared)
flow/              # Flow blockchain contracts and scripts
tests/             # Integration tests
examples/          # Example implementations
```

## Benefits

1. **Multiple APIs**: You can now run different API versions:
   - `./bin/api` - Current wallet API
   - `./bin/api-v2` - Future enhanced version
   - `./bin/admin` - Admin tools

2. **Clean Separation**: Business logic is in `internal/`, preventing external imports

3. **Scalable**: Easy to add new commands or API versions

4. **Standard Go Layout**: Follows Go community conventions

## Building

```bash
# Build all versions
go build -o bin/api ./cmd/api
go build -o bin/api-v2 ./cmd/api-v2  
go build -o bin/admin ./cmd/admin

# Or build current main.go (legacy)
go build -o bin/main .
```

## Running Different APIs

```bash
# Current API
./bin/api

# Future API v2 (on port 3001)
./bin/api-v2

# Admin tools
./bin/admin -command=migrate
```

## Future Enhancements

You can now easily:
- Add new API versions with different business logic
- Create specialized tools (migration, backup, monitoring)
- Share business logic between different entry points
- Test different configurations or feature sets

This structure supports your goal of having multiple API versions that can use different modules while sharing the core business logic.