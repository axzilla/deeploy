#!/bin/bash

echo "🚀 Installing deeploy v0.1.0..."

# # Check if Docker is installed
# if ! command -v docker &> /dev/null; then
#     echo "🐋 Installing Docker..."
#     curl -fsSL https://get.docker.com -o get-docker.sh
#     sudo sh get-docker.sh
# fi
#
# # Pull and run deeploy
# echo "📦 Starting deeploy..."
# docker pull ghcr.io/yourusername/deeploy:latest
# docker run -d \
#     --name deeploy \
#     -p 8090:8090 \
#     ghcr.io/yourusername/deeploy:latest
#
# echo "✨ Deeploy is running!"
# echo "📱 Visit http://localhost:8090"
