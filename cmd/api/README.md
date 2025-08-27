# Flow Wallet API - Testnet Configuration

This directory contains the Docker configuration for running the Flow Wallet API connected to Flow Testnet.

## Configuration

- **Network**: Flow Testnet
- **Database**: Supabase PostgreSQL (with connection pooling)
- **Mode**: Lightweight (no Redis dependency)
- **Admin Account**: `0xe6ddfc62f8ec2020`

## Usage

### Option 1: Docker Compose (Recommended)

```bash
cd cmd/api
docker-compose up -d
```

### Option 2: Direct Docker Build

```bash
cd cmd/api
docker build -t flow-wallet-api-testnet -f Dockerfile ../..
docker run -p 3000:3000 --env-file .env flow-wallet-api-testnet
```

### Option 3: Local Development

```bash
cd cmd/api
source .env
go run main.go
```

## API Endpoints

- **Health Check**: http://localhost:3000/v1/health/ready
- **API Documentation**: http://localhost:3000/v1/docs (if available)

## Database Migration

The API will automatically create tables on first run. If you need to reset the database, you can:

1. Connect to Supabase dashboard
2. Drop existing tables (if any)
3. Restart the API - it will recreate all tables

## Testnet Resources

- **FlowToken Contract**: `0x7e60df042a9c0868`
- **FungibleToken**: `0x9a0766d93b6608b7`
- **FlowStorageFees**: `0x8c5303eaa26202d6`
- **Access Node**: `access.testnet.nodes.onflow.org:9000`

## Logs

```bash
docker-compose logs -f api
```