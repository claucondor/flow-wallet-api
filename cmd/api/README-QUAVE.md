# Flow Wallet API - Quave Cloud Deployment

Este directorio contiene la configuración para desplegar la Flow Wallet API en Quave Cloud usando el repositorio público de GitHub.

## Archivos de Despliegue

- **`Dockerfile.github`** - Dockerfile que clona desde GitHub público
- **`docker-compose.github.yml`** - Configuración para Quave Cloud
- **`.env`** - Variables de entorno (actualizar para producción)

## Configuración para Quave

### 1. En Quave Cloud:

```yaml
# Use este Dockerfile
dockerfile: Dockerfile.github

# O use docker-compose
docker-compose: docker-compose.github.yml
```

### 2. Variables de Entorno Importantes:

```bash
# Testnet Configuration
FLOW_WALLET_ACCESS_API_HOST=access.testnet.nodes.onflow.org:9000
FLOW_WALLET_CHAIN_ID=flow-testnet

# Admin Account (Actualizar con tus credenciales)
FLOW_WALLET_ADMIN_ADDRESS=0xe6ddfc62f8ec2020
FLOW_WALLET_ADMIN_PRIVATE_KEY=1cf49f7f4971acb060223272b903126288b7067b65ce213358d90e08caf9c3d0

# Database (Actualizar con tu Supabase)
FLOW_WALLET_DATABASE_TYPE=psql
FLOW_WALLET_DATABASE_DSN=postgresql://tu-usuario:tu-password@tu-host:5432/postgres?sslmode=require

# Tokens
FLOW_WALLET_ENABLED_TOKENS=FlowToken:0x7e60df042a9c0868:flowToken

# Mode
FLOW_WALLET_LIGHTWEIGHT_MODE=true
FLOW_WALLET_DISABLE_IDEMPOTENCY_MIDDLEWARE=true
```

### 3. URLs de Acceso:

- **Health Check**: `https://tu-app.quave.cloud/v1/health/ready`
- **Accounts**: `https://tu-app.quave.cloud/v1/accounts`
- **Scripts**: `https://tu-app.quave.cloud/v1/scripts`

## Repositorio GitHub

La imagen se construye automáticamente desde:
```
https://github.com/claucondor/flow-wallet-api
```

## Comandos de Prueba Local

```bash
# Build desde GitHub
docker build -t flow-api-github -f Dockerfile.github .

# Run
docker run -p 3000:3000 --env-file .env flow-api-github

# Test
curl http://localhost:3000/v1/health/ready
```

## Notas de Seguridad

⚠️ **IMPORTANTE**: 
- Actualiza las credenciales de base de datos
- Cambia las claves privadas para producción
- Usa secrets de Quave para datos sensibles
- No hardcodees credenciales en el código

## Arquitectura

```
Quave Cloud → GitHub Repo → Docker Build → Flow Testnet + Supabase
```