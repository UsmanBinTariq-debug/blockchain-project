#!/bin/bash

# Deployment script for production

set -e

echo "🚀 Deploying Crypto Wallet System to Production..."

# Check if .env exists
if [ ! -f backend/.env ]; then
    echo "❌ Error: backend/.env file not found"
    exit 1
fi

echo "\n📦 Building backend..."
cd backend
go mod tidy
go build -o bin/wallet-server ./cmd/server

if [ ! -f bin/wallet-server ]; then
    echo "❌ Backend build failed"
    exit 1
fi

cd ..

echo "✅ Backend build successful"

echo "\n📦 Building frontend..."
cd frontend
npm install
npm run build

if [ ! -d dist ]; then
    echo "❌ Frontend build failed"
    exit 1
fi

cd ..

echo "✅ Frontend build successful"

echo "\n🚀 Deploying to Fly.io (backend)..."
flyctl deploy

echo "\n🚀 Deploying to Vercel (frontend)..."
cd frontend
vercel deploy --prod
cd ..

echo "\n✅ Deployment complete!"
echo "💡 Get your application URLs from Fly.io and Vercel dashboards"
