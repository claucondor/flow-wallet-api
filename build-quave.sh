#!/bin/bash
set -e

echo "🚀 Building Flow Wallet API for Quave Cloud..."

# Verificar que existan los archivos necesarios
if [ ! -f "Dockerfile.quave" ]; then
    echo "❌ Error: Dockerfile.quave no encontrado"
    exit 1
fi

if [ ! -f "openapi.yml" ]; then
    echo "❌ Error: openapi.yml no encontrado"
    exit 1
fi

# Nombre de la imagen
IMAGE_NAME="flow-wallet-api-quave"
TAG=${1:-latest}

echo "📦 Construyendo imagen: $IMAGE_NAME:$TAG"

# Build de la imagen
docker build -f Dockerfile.quave -t $IMAGE_NAME:$TAG .

echo "✅ Imagen construida exitosamente: $IMAGE_NAME:$TAG"

# Opcional: Mostrar el tamaño de la imagen
echo "📊 Tamaño de la imagen:"
docker images $IMAGE_NAME:$TAG

echo ""
echo "🔥 Para probar localmente:"
echo "   docker run -p 3000:3000 $IMAGE_NAME:$TAG"
echo ""
echo "📚 Endpoints disponibles:"
echo "   API: http://localhost:3000"
echo "   Docs: http://localhost:3000/docs"
echo "   Health: http://localhost:3000/v1/info"
echo ""
echo "☁️ Para subir a Quave Cloud:"
echo "   1. Tag la imagen con tu registry"
echo "   2. Push al registry"
echo "   3. Deploy desde Quave Cloud dashboard"