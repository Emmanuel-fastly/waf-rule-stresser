#!/bin/bash
set -e

echo "Building WAF Tester single binary..."

# Step 1: Build frontend
echo "1/4 Building frontend..."
cd frontend
npm run build
cd ..

# Step 2: Copy frontend dist to backend
echo "2/4 Copying frontend to backend..."
mkdir -p backend/frontend
cp -r frontend/dist backend/frontend/

# Step 3: Copy payloads to backend
echo "3/4 Copying payloads to backend..."
mkdir -p backend/payloads
cp -r payloads/*.txt backend/payloads/

# Step 4: Build Go binary
echo "4/4 Building Go binary..."
cd backend
go build -o ../waf-tester
cd ..

echo "✓ Build complete! Binary: waf-tester"
echo ""
echo "To run: ./waf-tester"
