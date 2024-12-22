#!/bin/bash

VERSION=${1:-latest}

echo "🚀 Installing deeploy version: $VERSION"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "🐋 Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
fi

# Pull and run deeploy with specific version
echo "📦 Starting deeploy..."
docker pull ghcr.io/axzilla/deeploy:$VERSION
docker run -d \
    --name deeploy \
    -p 8090:8090 \
    ghcr.io/axzilla/deeploy:$VERSION

echo "✨ Deeploy ($VERSION) is running!"
