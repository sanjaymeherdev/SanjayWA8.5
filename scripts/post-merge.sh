#!/bin/bash
set -e

echo "Running post-merge setup..."

cd src
echo "Tidying and downloading Go dependencies..."
go mod tidy
CGO_ENABLED=1 go mod download

echo "Building whatsapp binary..."
CGO_ENABLED=1 go build -o whatsapp .

echo "Creating required directories..."
mkdir -p storages statics/qrcode statics/senditems statics/media

echo "Post-merge setup complete."
