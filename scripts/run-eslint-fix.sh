#!/bin/bash

# Run ESLint with --fix option to automatically fix issues
echo "Running ESLint auto-fix..."
cd frontend
npm run lint:fix

echo "ESLint auto-fix complete!"
