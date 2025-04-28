#!/bin/bash

# Script to fix DeviceListResponse import cycle and type conflicts
echo "Creating backup of original files..."
BACKUP_DIR="./frontend/src_backup_$(date +%Y%m%d_%H%M%S)"
mkdir -p $BACKUP_DIR
cp -r ./frontend/src/api ./frontend/src/types $BACKUP_DIR

# Fix cyclic import in device.ts
echo "Fixing cyclic import in device.ts..."
sed -i '1d' ./frontend/src/api/device.ts

# Remove duplicate DeviceListResponse from mouse.ts
echo "Removing duplicate DeviceListResponse from mouse.ts..."
sed -i '/export interface DeviceListResponse {/,/}/d' ./frontend/src/types/mouse.ts

# Update imports in mouse.ts
echo "Updating imports in mouse.ts..."
sed -i '1s|.*|import { DeviceListResponse } from "@/api/device";|' ./frontend/src/types/mouse.ts

# Check components that use DeviceListResponse and fix imports if needed
echo "Checking and fixing components that use DeviceListResponse..."
find ./frontend/src -type f -name "*.vue" -o -name "*.ts" | xargs grep -l "DeviceListResponse" | while read file; do
  # Skip the api/device.ts file since we already fixed it
  if [[ "$file" != "./frontend/src/api/device.ts" && "$file" != "./frontend/src/types/mouse.ts" ]]; then
    # Check if it imports DeviceListResponse
    if ! grep -q "import.*DeviceListResponse.*from" "$file"; then
      # Add import if missing
      if [[ "$file" == *".vue" ]]; then
        # For Vue files, add after script tag
        sed -i '/<script.*>/a\
import { DeviceListResponse } from "@/api/device";' "$file"
        echo "Added import to $file"
      else
        # For TS files, add at the top
        sed -i '1i\
import { DeviceListResponse } from "@/api/device";' "$file"
        echo "Added import to $file"
      fi
    fi
  fi
done

echo "DeviceListResponse conflicts resolved successfully!"