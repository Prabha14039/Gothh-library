#!/bin/bash

set -e

echo "📦 Updating system packages..."
sudo apt-get update -y

echo "🔧 Installing make..."
sudo apt-get install -y make

echo "🐳 Installing Docker using get.docker.com script..."
curl -fsSL https://get.docker.com | sh

echo "✅ Docker installed successfully."

echo "➕ Adding current user to docker group (optional)..."
sudo usermod -aG docker $USER
echo "⚠️ You may need to log out and log back in for group changes to take effect."

echo "🔁 Installing Docker Compose v2 (already included in Docker >= v20.10)..."
sudo apt install docker-compose

echo "🎉 Done. Verifying installation:"
docker --version
docker compose version

