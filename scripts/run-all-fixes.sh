#!/bin/bash

# Master script to run all TypeScript error fixes in optimal order

echo "========================================================"
echo "Starting comprehensive TypeScript error fixing process..."
echo "========================================================"

# Create backup before any changes
BACKUP_DIR="./frontend/src_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR
cp -r ./frontend/src $BACKUP_DIR
echo "Created backup at: $BACKUP_DIR"

# Step 1: Install ts-morph if needed
echo -e "\n\n--- Step 1: Setting up ts-morph ---"
./scripts/install-ts-morph.sh

# Step 2: Run basic TypeScript error fixes
echo -e "\n\n--- Step 2: Running basic TypeScript error fixes ---"
./scripts/fix-ts-errors.sh

# Step 3: Run advanced ts-morph fixes
echo -e "\n\n--- Step 3: Running advanced ts-morph fixes ---"
cd frontend
sudo npx node ../scripts/fix-with-ts-morph.js
cd ..

# Step 4: Set up enhanced ESLint
echo -e "\n\n--- Step 4: Setting up enhanced ESLint ---"
./scripts/enhanced-eslint.sh

# Step 5: Run ESLint auto-fixes
echo -e "\n\n--- Step 5: Running ESLint auto-fixes ---"
cd frontend
sudo npm run lint:fix || echo "ESLint fix had some issues, continuing..."
cd ..

# Step 6: Verify build
echo -e "\n\n--- Step 6: Verifying build ---"
cd frontend
sudo npm run build || echo "Build still has some issues that need manual fixing"
cd ..

echo -e "\n========================================================"
echo "TypeScript error fixing process complete!"
echo "If build errors remain, check the output above for details."
echo "Manual fixes may still be required for complex issues."
echo "========================================================"
