# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Run/Test Commands
- **Backend**: 
  - Run all tests: `go test ./...`
  - Run specific test: `go test -run TestName ./path/to/package`
  - Run test function: `go test -run TestService_Login ./services/auth`
  - Verbose testing: `go test -v ./...`
- **Frontend**:
  - Dev server: `npm run dev`
  - Build: `npm run build`

## Code Guidelines
- **Backend**:
  - Go conventions: PascalCase for exported, camelCase for private
  - Error handling: Use custom AppError with error codes
  - Types: Strong typing with dedicated types packages
  - Testing: Use testify/gomock, follow Test{Package}_{Function} naming
- **Frontend**:
  - Vue 3 + TypeScript
  - Component names: PascalCase
  - Variables/props: camelCase
  - Use composables for shared logic
  - Type everything appropriately

## Development Plan (Mouse Comparison Website)

### Phase 1: Core Infrastructure and Mouse Comparison (2-3 weeks)
- [Partially Complete] Frontend and backend infrastructure
- [Partially Complete] Database design and initialization
- [In Progress] Mouse data API
- [Todo] SVG display and comparison components
- [Todo] Draggable ruler refinement

### Phase 2: Similarity Search and Database Browsing (3-4 weeks)
- [Todo] Similarity algorithm
- [Todo] Similarity search interface
- [Todo] Database browsing and filtering
- [Todo] Mouse detail pages

### Phase 3: User Authentication and Shopping (2-3 weeks)
- [Partially Complete] User registration and login
- [Todo] Shopping cart functionality
- [Todo] Checkout process
- [Todo] Order management

### Phase 4: Integration Testing and Optimization (1-2 weeks)
- [Todo] Site-wide integration testing
- [Todo] UI optimization
- [Todo] Performance optimization
- [Todo] Deployment preparation

## Internationalization Support (i18n)
The project implements complete internationalization:
- Frontend: Vue-i18n, supporting English and Chinese
- Backend: Custom i18n service, injected via middleware
- Language files: JSON format, organized by feature modules

## Current Priority Tasks
1. Complete mouse comparison core functionality (SVG display and comparison)
2. Implement mouse data models and API
3. Develop similarity algorithm and search functionality
4. Refine shopping cart and checkout process

## Recent Fixes (2024-04-26)
1. Fixed mouse selection issues in comparison tool
2. Fixed SVG rendering errors in comparison view 
3. Removed "difference" column from comparison table
4. Fixed errors in review list page related to score handling
5. Extracted similarity mouse finder to its own dedicated page
6. Updated home page button styles to full-width bar style
7. Enhanced error handling in SVG service for better user experience
8. Added fallbacks for missing data in components