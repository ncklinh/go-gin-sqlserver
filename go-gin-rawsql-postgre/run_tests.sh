#!/bin/bash

echo "Running Authentication Middleware Tests..."
echo "=========================================="

# Run middleware tests
echo "1. Testing Auth Middleware..."
go test -v ./middleware -run TestAuthMiddleware

echo ""
echo "2. Testing Permission Middleware..."
go test -v ./middleware -run TestRequirePermission

echo ""
echo "3. Testing All Middleware Tests..."
go test -v ./middleware

echo ""
echo "4. Testing Model Permission Functions..."
go test -v ./model

echo ""
echo "5. Testing Token Package..."
go test -v ./token

echo ""
echo "=========================================="
echo "Test Summary:"
echo "- Auth middleware tests verify JWT token validation"
echo "- Permission middleware tests verify role-based access control"
echo "- Model tests verify permission checking logic"
echo "- Token tests verify JWT creation and validation"
echo ""
echo "Note: Handler tests require database connection for full functionality"
echo "but middleware tests work independently." 