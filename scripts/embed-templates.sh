#!/bin/bash

# Script to embed templates into the binary
set -e

echo "🔧 Embedding templates into binary..."

# Create bindata.go file with embedded templates
go-bindata -o internal/bindata.go -pkg internal templates/...

echo "✅ Templates embedded successfully!"
echo "📁 Generated: internal/bindata.go" 