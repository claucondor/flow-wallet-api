# 🐳 Docker Architecture

¡Ya no tienes 13 archivos Docker! Todo ha sido consolidado en una arquitectura limpia y escalable.

## 📦 Antes vs Después

### ❌ **ANTES (Caos total):**
```
docker-compose.yml
docker-compose.dev.yml
docker-compose.lightweight.yml
docker-compose.quave.lightweight.yml
docker-compose.test-suite.yml
Dockerfile.quave
Dockerfile.quave-final
docker/wallet/Dockerfile
docker/flow-cli/Dockerfile
... 8 compose + 5 dockerfiles = 13 archivos!
```

### ✅ **AHORA (Limpio y organizado):**
```
Dockerfile                 # UN solo Dockerfile multi-stage
docker-compose.yml         # UN solo compose con profiles
.env.example              # Configuración centralizada
Makefile.docker           # Comandos simplificados
docker/legacy/            # Archivos viejos (backup)
```

## 🎯 Arquitectura Unificada

### **1. Dockerfile Multi-Stage**
Un solo archivo que genera múltiples versiones:

```dockerfile
# Targets disponibles:
api-production      # API principal
api-v2-production   # API futura versión
admin-production    # Herramientas admin  
legacy-production   # Compatibilidad
development         # Desarrollo con hot-reload
lightweight         # Minimal (scratch)
all-in-one         # Todos los binarios
```

### **2. Docker Compose Profiles**
Un solo archivo, múltiples configuraciones:

```yaml
# Profiles disponibles:
dev          # Stack completo de desarrollo
prod         # Producción con Nginx
lightweight  # Solo SQLite + API + Emulator
test         # Testing environment
api          # Solo API con dependencias
api-v2       # API V2 con dependencias  
admin        # Herramientas admin
docs         # Solo documentación
```

## 🚀 Uso Súper Simple

### **Comandos Directos:**
```bash
# Desarrollo completo
docker compose --profile dev up

# Modo ligero (SQLite)
docker compose --profile lightweight up

# Producción
docker compose --profile prod up -d

# API V2 (futuro)
docker compose --profile api-v2 up

# Solo admin tools
docker compose --profile admin run admin -command=migrate
```

### **Con Makefile (Aún más fácil):**
```bash
# Desarrollo
make dev

# Producción  
make prod

# Lightweight
make lightweight

# Ver logs
make logs SERVICE=api

# Shell en container
make shell SERVICE=postgres

# Limpiar todo
make clean
```

## 🎮 Ejemplos Prácticos

### **Desarrollo Local:**
```bash
# Setup inicial
make env-setup          # Crea .env desde .env.example
make dev               # Levanta todo (Postgres + Redis + API + Emulator)

# URLs disponibles:
# API: http://localhost:3000
# PgAdmin: http://localhost:5050  
# Docs: http://localhost:8080
```

### **Modo Súper Ligero:**
```bash
make lightweight       # Solo SQLite + API + Emulator
# Perfecto para desarrollo rápido o demos
```

### **Múltiples APIs:**
```bash
# API actual en puerto 3000
docker compose --profile api up -d

# API V2 en puerto 3001  
docker compose --profile api-v2 up -d

# Ambas corriendo simultáneamente!
```

### **Testing:**
```bash
make test              # Entorno de testing
make test-integration  # Tests de integración
```

## 🗂️ Estructura de Archivos

```
docker/
├── legacy/            # Archivos viejos (backup)
│   ├── docker-compose.dev.yml
│   ├── docker-compose.lightweight.yml  
│   ├── Dockerfile.quave
│   └── ... (todos los viejos)
├── nginx/            # Configs de Nginx (futuro)
└── pgadmin/         # Configs de PgAdmin (futuro)

Dockerfile           # Multi-stage unificado
docker-compose.yml   # Profiles unificado
.env.example         # Configuración template
Makefile.docker      # Comandos simplificados
```

## 🔧 Variables de Entorno

Todo centralizado en `.env`:

```bash
# Puertos
API_PORT=3000
API_V2_PORT=3001

# Database  
DATABASE_TYPE=postgres
POSTGRES_PASSWORD=wallet

# Modos
LIGHTWEIGHT_MODE=false
QUAVE_MODE=false

# Flow
CHAIN_ID=flow-emulator
```

## 🎉 Beneficios

1. **13 → 4 archivos**: Reducción masiva de complejidad
2. **Un comando, múltiples entornos**: `make dev`, `make prod`, `make lightweight`
3. **Escalable**: Fácil agregar nuevas APIs o servicios
4. **Mantenible**: Un solo lugar para cambios
5. **Testeable**: Perfiles específicos para testing
6. **Producción-ready**: Con Nginx, healthchecks, etc.

## 🚨 Migration Path

Los archivos viejos están en `docker/legacy/` por si necesitas referenciar algo. ¡Pero ya no los necesitas!

**¡De 13 archivos caóticos a 4 archivos organizados!** 🎊