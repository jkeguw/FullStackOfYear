#!/bin/bash

# Add enhanced ESLint configuration for better type checking and code style enforcement

echo "Setting up enhanced ESLint configuration..."
cd frontend

# Install additional ESLint plugins for Vue and TypeScript
npm install --save-dev \
  @typescript-eslint/eslint-plugin \
  @typescript-eslint/parser \
  eslint-plugin-vue \
  eslint-config-prettier \
  eslint-plugin-prettier \
  typescript

# Create enhanced ESLint configuration file
cat > ./.eslintrc.js << 'EOF'
module.exports = {
  root: true,
  env: {
    node: true,
    browser: true,
    es2021: true
  },
  extends: [
    'plugin:vue/vue3-recommended',
    'eslint:recommended',
    '@vue/typescript/recommended',
    'prettier'
  ],
  parserOptions: {
    ecmaVersion: 2021,
    parser: '@typescript-eslint/parser'
  },
  plugins: ['@typescript-eslint', 'prettier'],
  rules: {
    // Vue specific rules
    'vue/component-name-in-template-casing': ['error', 'PascalCase'],
    'vue/component-definition-name-casing': ['error', 'PascalCase'],
    
    // TypeScript rules
    '@typescript-eslint/no-explicit-any': 'warn',
    '@typescript-eslint/explicit-module-boundary-types': 'off',
    '@typescript-eslint/no-unused-vars': ['warn', { 
      argsIgnorePattern: '^_',
      varsIgnorePattern: '^_' 
    }],
    '@typescript-eslint/naming-convention': [
      'error',
      { 
        selector: 'default', 
        format: ['camelCase'],
        leadingUnderscore: 'allow'
      },
      { 
        selector: 'variable', 
        format: ['camelCase', 'UPPER_CASE', 'PascalCase'] 
      },
      { 
        selector: 'parameter', 
        format: ['camelCase'], 
        leadingUnderscore: 'allow' 
      },
      { 
        selector: 'memberLike', 
        format: ['camelCase'] 
      },
      { 
        selector: 'typeLike', 
        format: ['PascalCase'] 
      },
      { 
        selector: 'property',
        format: ['camelCase'],
        leadingUnderscore: 'allow'
      },
      { 
        selector: 'enum', 
        format: ['PascalCase'] 
      },
      { 
        selector: 'enumMember', 
        format: ['UPPER_CASE'] 
      }
    ],
    
    // Style rules
    'prettier/prettier': ['error', {
      singleQuote: true,
      semi: true,
      tabWidth: 2,
      trailingComma: 'none',
      printWidth: 100,
      endOfLine: 'auto'
    }],
    'no-console': process.env.NODE_ENV === 'production' ? 'warn' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'warn' : 'off'
  }
};
EOF

# Create a .prettierrc file for consistency
cat > ./.prettierrc << 'EOF'
{
  "singleQuote": true,
  "semi": true,
  "tabWidth": 2,
  "trailingComma": "none",
  "printWidth": 100,
  "endOfLine": "auto"
}
EOF

# Add ESLint script to package.json
if ! grep -q '"lint:fix"' package.json; then
  # Use temporary file to ensure proper JSON formatting
  cat > ./temp-package.json << 'EOF'
  "scripts": {
    "lint": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts --ignore-path .gitignore",
    "lint:fix": "eslint . --ext .vue,.js,.jsx,.cjs,.mjs,.ts,.tsx,.cts,.mts --fix --ignore-path .gitignore"
EOF

  # Insert the new scripts into package.json
  sed -i '/\"scripts\": {/r ./temp-package.json' package.json
  rm ./temp-package.json
fi

# Create a script for running the ESLint fixes
cat > ../scripts/run-eslint-fix.sh << 'EOF'
#!/bin/bash

# Run ESLint with --fix option to automatically fix issues
echo "Running ESLint auto-fix..."
cd frontend
npm run lint:fix

echo "ESLint auto-fix complete!"
EOF

chmod +x ../scripts/run-eslint-fix.sh

echo "Enhanced ESLint configuration has been set up!"
cd ..