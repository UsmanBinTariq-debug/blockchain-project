#!/bin/bash

# Setup script for development environment

echo "🚀 Setting up Crypto Wallet System..."

# Create necessary directories
mkdir -p backend backend/bin frontend/dist database/migrations

# Backend setup
echo "\n📦 Setting up backend..."
cd backend

# Copy environment file
if [ ! -f .env ]; then
    cp .env.example .env
    echo "✅ Created .env file (update with your values)"
fi

# Download Go dependencies
go mod download
go mod tidy

# Build backend
echo "\n🔨 Building backend..."
go build -o bin/wallet-server ./cmd/server

if [ -f bin/wallet-server ]; then
    echo "✅ Backend built successfully"
else
    echo "❌ Backend build failed"
    exit 1
fi

cd ..

# Frontend setup
echo "\n📦 Setting up frontend..."
cd frontend

# Copy environment file
if [ ! -f .env.development ]; then
    cp .env.example .env.development
    echo "✅ Created .env.development file"
fi

# Install dependencies
npm install

cd ..

# Database setup
echo "\n💾 Database setup instructions:"
echo "1. Create a Supabase project at https://supabase.com"
echo "2. Run the SQL schema in database/schema.sql via Supabase SQL editor"
echo "3. Get your DATABASE_URL from Project Settings → Database"
echo "4. Update backend/.env with your DATABASE_URL"

echo "\n✅ Setup complete!"
echo "\n📝 Next steps:"
echo "1. Update backend/.env with your database URL"
echo "2. Update frontend/.env.development with API URL"
echo "3. Run: npm run dev (from frontend directory)"
echo "4. Run: go run ./cmd/server/main.go (from backend directory)"
echo "5. Open http://localhost:5173 in your browser"
