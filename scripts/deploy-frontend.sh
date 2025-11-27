#!/bin/bash

echo "🚀 Deploying Frontend to port 3001..."

FRONTEND_DIR="/opt/asmo-frontend"
FRONTEND_REPO="https://github.com/ASMO-team/ASMO-site-frontend.git"
PORT=3001

# Create directory
mkdir -p $FRONTEND_DIR
cd $FRONTEND_DIR

# Clone or update frontend
if [ -d ".git" ]; then
    echo "📥 Updating existing frontend repository..."
    git pull origin main
else
    echo "📥 Cloning frontend repository..."
    git clone $FRONTEND_REPO .
fi

# Install dependencies
echo "📦 Installing dependencies..."
npm install

# Build frontend
echo "🔨 Building frontend..."
npm run build

# Stop existing frontend process
echo "🛑 Stopping existing frontend..."
pkill -f "next start" || true

# Start frontend on port 3001
echo "🚀 Starting frontend on port $PORT..."
nohup npm run start -- -p $PORT > frontend.log 2>&1 &

echo "✅ Frontend deployed on port $PORT"
echo "📄 Logs: $FRONTEND_DIR/frontend.log"