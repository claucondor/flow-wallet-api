# Installation Guide

This guide will walk you through installing and setting up Phoenix Wallet API in different environments. Choose the installation method that best fits your needs.

## 🚀 **Quick Installation (Recommended)**

The fastest way to get started with Phoenix Wallet API:

```bash
# Clone the repository
git clone https://github.com/flow-hydraulics/flow-wallet-api.git
cd flow-wallet-api

# Start in lightweight mode
make lightweight

# Your API is now running!
# API: http://localhost:3000/v1
# Documentation: http://localhost:8080
```

That's it! Phoenix Wallet API is now running with:
- ✅ Flow emulator for local development
- ✅ SQLite database for data persistence
- ✅ Complete documentation interface
- ✅ All core features enabled

## 📋 **Prerequisites**

Before installing, ensure you have the following installed:

### **Required**
- **Docker** (version 20.0 or higher)
- **Docker Compose** (version 2.0 or higher)
- **Git** for cloning the repository

### **Optional**
- **Make** (for using Makefile commands)
- **Node.js** (for running documentation locally)
- **Flow CLI** (for advanced Flow blockchain operations)

### **System Requirements**

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **RAM** | 2GB | 4GB+ |
| **CPU** | 2 cores | 4+ cores |
| **Storage** | 5GB | 20GB+ |
| **Network** | Internet connection | Stable broadband |

## 🐳 **Docker Installation**

### **Verify Docker Installation**

```bash
# Check Docker version
docker --version
# Should output: Docker version 20.x.x or higher

# Check Docker Compose version
docker compose version
# Should output: Docker Compose version 2.x.x or higher

# Test Docker is working
docker run hello-world
```

### **Install Docker (if needed)**

#### **Ubuntu/Debian**
```bash
# Update package index
sudo apt-get update

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Add user to docker group
sudo usermod -aG docker $USER

# Install Docker Compose
sudo apt-get install docker-compose-plugin
```

#### **macOS**
```bash
# Install Docker Desktop
brew install --cask docker

# Or download from: https://www.docker.com/products/docker-desktop
```

#### **Windows**
1. Download Docker Desktop from [docker.com](https://www.docker.com/products/docker-desktop)
2. Run the installer and follow the setup wizard
3. Restart your computer when prompted

## 🛠️ **Installation Methods**

### **Method 1: Lightweight Mode (Recommended for Development)**

Perfect for development, testing, and small-scale production:

```bash
# Clone repository
git clone https://github.com/flow-hydraulics/flow-wallet-api.git
cd flow-wallet-api

# Copy environment configuration
cp .env.example .env

# Start lightweight mode
make lightweight
```

**What this includes:**
- Phoenix Wallet API server
- Flow emulator (local blockchain)
- SQLite database
- Documentation interface
- Swagger UI for API testing

**Endpoints:**
- **API**: http://localhost:3000/v1
- **Documentation**: http://localhost:8080
- **Flow Emulator**: localhost:3569

### **Method 2: Standard Mode (Production)**

For production deployments with full features:

```bash
# Clone repository
git clone https://github.com/flow-hydraulics/flow-wallet-api.git
cd flow-wallet-api

# Copy and configure environment
cp .env.example .env
# Edit .env with your production settings

# Start standard mode
make dev
```

**What this includes:**
- Phoenix Wallet API server
- PostgreSQL database
- Redis cache and job queue
- Flow emulator
- pgAdmin for database management

**Endpoints:**
- **API**: http://localhost:3000/v1
- **pgAdmin**: http://localhost:5050
- **Flow Emulator**: localhost:3569

### **Method 3: Network-Specific Deployment**

Deploy directly to Flow Testnet or Mainnet:

#### **Flow Testnet**
```bash
# Configure for testnet
make lightweight-testnet

# With idempotency (recommended)
make lightweight-testnet-idempotent
```

#### **Flow Mainnet**
```bash
# Configure for mainnet (requires real FLOW tokens)
make lightweight-mainnet-idempotent
```

## ⚙️ **Configuration**

### **Environment Variables**

Phoenix Wallet API is configured through environment variables. Copy the example configuration:

```bash
cp .env.example .env
```

### **Essential Configuration**

Edit `.env` file with your settings:

```bash
# Admin account (required for testnet/mainnet)
FLOW_WALLET_ADMIN_ADDRESS=0xf8d6e0586b0a20c7
FLOW_WALLET_ADMIN_PRIVATE_KEY=your-private-key-here

# Network configuration
FLOW_WALLET_ACCESS_API_HOST=localhost:3569  # Emulator
FLOW_WALLET_CHAIN_ID=flow-emulator

# Database configuration
FLOW_WALLET_DATABASE_TYPE=sqlite
FLOW_WALLET_DATABASE_DSN=./data/wallet.db

# Security
FLOW_WALLET_ENCRYPTION_KEY=your-32-character-encryption-key
```

### **Network-Specific Configuration**

#### **For Flow Emulator (Development)**
```bash
FLOW_WALLET_ACCESS_API_HOST=localhost:3569
FLOW_WALLET_CHAIN_ID=flow-emulator
FLOW_WALLET_ADMIN_ADDRESS=0xf8d6e0586b0a20c7
FLOW_WALLET_ADMIN_PRIVATE_KEY=91a22fbd87392b019fbe332c32695c14cf2ba5b6521476a8540228bdf1987068
```

#### **For Flow Testnet**
```bash
FLOW_WALLET_ACCESS_API_HOST=access.testnet.nodes.onflow.org:9000
FLOW_WALLET_CHAIN_ID=flow-testnet
FLOW_WALLET_ADMIN_ADDRESS=your-testnet-address
FLOW_WALLET_ADMIN_PRIVATE_KEY=your-testnet-private-key
```

#### **For Flow Mainnet**
```bash
FLOW_WALLET_ACCESS_API_HOST=access.mainnet.nodes.onflow.org:9000
FLOW_WALLET_CHAIN_ID=flow-mainnet
FLOW_WALLET_ADMIN_ADDRESS=your-mainnet-address
FLOW_WALLET_ADMIN_PRIVATE_KEY=your-mainnet-private-key
```

## 🔐 **Security Configuration**

### **Key Management Options**

#### **Local Key Storage (Development)**
```bash
FLOW_WALLET_DEFAULT_KEY_TYPE=local
FLOW_WALLET_ENCRYPTION_KEY=your-32-character-encryption-key
```

#### **Google Cloud KMS (Production)**
```bash
FLOW_WALLET_DEFAULT_KEY_TYPE=google_kms
FLOW_WALLET_GOOGLE_KMS_PROJECT_ID=your-gcp-project
FLOW_WALLET_GOOGLE_KMS_LOCATION_ID=us-central1
FLOW_WALLET_GOOGLE_KMS_KEYRING_ID=your-keyring
```

#### **AWS KMS (Production)**
```bash
FLOW_WALLET_DEFAULT_KEY_TYPE=aws_kms
AWS_REGION=us-west-2
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
```

### **Idempotency Configuration**

Enable idempotency to prevent duplicate operations:

```bash
# For lightweight mode
FLOW_WALLET_LIGHTWEIGHT_IDEMPOTENCY=true

# For standard mode
FLOW_WALLET_DISABLE_IDEMPOTENCY_MIDDLEWARE=false
FLOW_WALLET_IDEMPOTENCY_MIDDLEWARE_DATABASE_TYPE=redis
FLOW_WALLET_IDEMPOTENCY_MIDDLEWARE_REDIS_URL=redis://localhost:6379/
```

## ✅ **Verification**

### **Check Installation**

1. **Verify services are running:**
   ```bash
   docker compose ps
   ```

2. **Test API connectivity:**
   ```bash
   curl http://localhost:3000/v1/health/ready
   ```

3. **Check API documentation:**
   Open http://localhost:8080 in your browser

4. **Create test account:**
   ```bash
   curl -X POST http://localhost:3000/v1/accounts \
     -H "Content-Type: application/json" \
     -d '{}'
   ```

### **Expected Response**
```json
{
  "address": "0x1234567890abcdef",
  "keys": [
    {
      "index": 0,
      "publicKey": "04a1b2c3d4e5f6...",
      "signAlgo": "ECDSA_P256",
      "hashAlgo": "SHA3_256",
      "weight": 1000
    }
  ],
  "balance": "0.00000000"
}
```

## 🔧 **Troubleshooting**

### **Common Issues**

#### **Port Already in Use**
```bash
# Error: Port 3000 is already in use
# Solution: Stop conflicting services or change port
docker compose down
# Or change port in docker-compose.yml
```

#### **Docker Permission Denied**
```bash
# Error: Permission denied while trying to connect to Docker daemon
# Solution: Add user to docker group
sudo usermod -aG docker $USER
# Then logout and login again
```

#### **Flow Emulator Connection Failed**
```bash
# Error: Failed to connect to Flow emulator
# Solution: Wait for emulator to start completely
docker compose logs emulator
# Wait for "Flow Emulator started" message
```

#### **Database Connection Error**
```bash
# Error: Database connection failed
# Solution: Check database configuration
docker compose logs api
# Verify DATABASE_DSN in .env file
```

### **Logs and Debugging**

```bash
# View all service logs
docker compose logs

# View specific service logs
docker compose logs api
docker compose logs emulator

# Follow logs in real-time
docker compose logs -f api

# Check service status
docker compose ps
```

## 🚀 **Next Steps**

Now that Phoenix Wallet API is installed and running:

1. **[Quick Start Guide](./quick-start)** - Learn basic operations
2. **[Core Concepts](../concepts/architecture)** - Understand the architecture
3. **[API Reference](../api-reference/overview)** - Explore all endpoints
4. **[Production Deployment](../deployment/production-setup)** - Deploy to production

## 📚 **Additional Resources**

### **Development Tools**
- **[Flow CLI](https://docs.onflow.org/flow-cli/)** - Command-line tools for Flow
- **[Flow Emulator](https://docs.onflow.org/emulator/)** - Local Flow blockchain
- **[Cadence](https://docs.onflow.org/cadence/)** - Flow's smart contract language

### **Monitoring and Management**
- **Docker Dashboard** - Visual container management
- **Portainer** - Web-based Docker management
- **Grafana + Prometheus** - Advanced monitoring (for production)

### **IDE Extensions**
- **Cadence VS Code Extension** - Syntax highlighting for Cadence
- **Flow VS Code Extension** - Flow development tools
- **Docker VS Code Extension** - Container management

Ready to start building? Let's move on to the [Quick Start Guide](./quick-start)! 🚀