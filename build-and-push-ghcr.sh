#!/bin/bash
set -e

# Configuración
GITHUB_USER=${1:-"claucondor"}
IMAGE_NAME="flow-wallet-api-quave"
TAG=${2:-"latest"}
FULL_IMAGE="ghcr.io/$GITHUB_USER/$IMAGE_NAME:$TAG"

echo "🚀 Building and pushing Flow Wallet API to GitHub Container Registry..."
echo "📦 Image: $FULL_IMAGE"

# Verificar que estemos autenticados
if ! docker system info | grep -q "ghcr.io"; then
    echo "⚠️  No estás autenticado con GitHub Container Registry"
    echo "   Ejecuta primero:"
    echo "   echo 'TU_TOKEN' | docker login ghcr.io -u $GITHUB_USER --password-stdin"
    exit 1
fi

# Verificar archivos necesarios
if [ ! -f "Dockerfile.quave" ]; then
    echo "❌ Error: Dockerfile.quave no encontrado"
    exit 1
fi

if [ ! -f "openapi.yml" ]; then
    echo "❌ Error: openapi.yml no encontrado"
    exit 1
fi

echo "🔨 Construyendo imagen..."
docker build -f Dockerfile.quave -t $FULL_IMAGE .

echo "📤 Subiendo imagen a GitHub Container Registry..."
docker push $FULL_IMAGE

echo "✅ ¡Imagen subida exitosamente!"
echo ""
echo "🔗 URL de la imagen: $FULL_IMAGE"
echo ""
echo "☁️ Para usar en Quave Cloud:"
echo "   1. Ve a tu dashboard de Quave"
echo "   2. Crea nueva aplicación"
echo "   3. Usa esta imagen: $FULL_IMAGE"
echo "   4. Configura puerto: 3000"
echo "   5. Health check: /v1/info"
echo ""
echo "🎯 Para hacer pública la imagen (recomendado):"
echo "   Ve a: https://github.com/$GITHUB_USER?tab=packages"
echo "   Selecciona '$IMAGE_NAME'"
echo "   Settings → Change visibility → Public"