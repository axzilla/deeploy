#!/bin/bash
# Check if script is running as root by checking user ID
if [ "$(id -u)" != "0" ]; then   # id -u returns user ID, root is always 0
   echo "This script must be run as root" >&2   # Print to stderr (error output)
   exit 1   # Exit with error code 1 (non-zero = error)
fi

# Get version from first argument or use 'latest' as default
VERSION=${1:-latest}
echo "🚀 Installing deeploy version: $VERSION"

# Check if Docker is installed by looking for docker command
if command -v docker > /dev/null 2>&1; then
   echo "✅ Docker already installed"
else
   # Install Docker using official install script
   echo "🐋 Installing Docker..."
   curl -sSL https://get.docker.com | sh
fi

# Ensure volume exists (create if not present)
VOLUME_NAME="deeploy_data"
if ! docker volume inspect $VOLUME_NAME > /dev/null 2>&1; then
   echo "📂 Creating Docker volume: $VOLUME_NAME"
   docker volume create $VOLUME_NAME
else
   echo "📂 Docker volume $VOLUME_NAME already exists"
fi

echo "📦 Starting deeploy..."
# Pull specified version of image (will check if update available)
docker pull ghcr.io/axzilla/deeploy:$VERSION

# Remove existing container if it exists
# -f forces removal even if running
# 2>/dev/null hides error if container doesn't exist
# || true continues even if command fails (no container found)
docker rm -f deeploy 2>/dev/null || true

# Start new container
# -d runs in detached (background) mode
# --name gives container a name for easier management
# -p maps host port 8090 to container port 8090 
docker run -d \
   --name deeploy \
   -p 8090:8090 \
   -v $VOLUME_NAME:/app/data \
   ghcr.io/axzilla/deeploy:$VERSION

# Get IP address for URL display
# Checks if hostname command exists and -I option works (Linux)
# Falls back to localhost for macOS/other systems
if command -v hostname &>/dev/null && hostname -I &>/dev/null; then
   IP=$(hostname -I | awk '{print $1}')
else
   IP="localhost"
fi

echo "✨ Deeploy ($VERSION) is running on http://$IP:8090"
