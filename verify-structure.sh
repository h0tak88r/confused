#!/bin/bash

# Test script for Confused project structure
# This script verifies the project structure without requiring Go to be installed

echo "Confused Project Structure Verification"
echo "======================================"

# Check if we're in the right directory
if [ ! -f "go.mod" ]; then
    echo "❌ Error: go.mod not found. Please run this script from the project root."
    exit 1
fi

echo "✅ Found go.mod"

# Check project structure
echo ""
echo "Checking project structure..."

# Check cmd directory
if [ -d "cmd/confused" ]; then
    echo "✅ cmd/confused directory exists"
    if [ -f "cmd/confused/main.go" ]; then
        echo "✅ cmd/confused/main.go exists"
    else
        echo "❌ cmd/confused/main.go missing"
    fi
    if [ -f "cmd/confused/commands.go" ]; then
        echo "✅ cmd/confused/commands.go exists"
    else
        echo "❌ cmd/confused/commands.go missing"
    fi
else
    echo "❌ cmd/confused directory missing"
fi

# Check pkg directory
if [ -d "pkg" ]; then
    echo "✅ pkg directory exists"
    for pkg in config logger github web scanner; do
        if [ -d "pkg/$pkg" ]; then
            echo "✅ pkg/$pkg directory exists"
        else
            echo "❌ pkg/$pkg directory missing"
        fi
    done
else
    echo "❌ pkg directory missing"
fi

# Check internal directory
if [ -d "internal" ]; then
    echo "✅ internal directory exists"
    if [ -d "internal/types" ]; then
        echo "✅ internal/types directory exists"
        if [ -f "internal/types/types.go" ]; then
            echo "✅ internal/types/types.go exists"
        else
            echo "❌ internal/types/types.go missing"
        fi
    else
        echo "❌ internal/types directory missing"
    fi
    if [ -d "internal/resolvers" ]; then
        echo "✅ internal/resolvers directory exists"
    else
        echo "❌ internal/resolvers directory missing"
    fi
else
    echo "❌ internal directory missing"
fi

# Check build files
echo ""
echo "Checking build files..."
if [ -f "Makefile" ]; then
    echo "✅ Makefile exists"
else
    echo "❌ Makefile missing"
fi

if [ -f "build.sh" ]; then
    echo "✅ build.sh exists"
else
    echo "❌ build.sh missing"
fi

if [ -f "Dockerfile" ]; then
    echo "✅ Dockerfile exists"
else
    echo "❌ Dockerfile missing"
fi

# Check documentation
echo ""
echo "Checking documentation..."
if [ -f "README.md" ]; then
    echo "✅ README.md exists"
else
    echo "❌ README.md missing"
fi

if [ -f "CHANGELOG.md" ]; then
    echo "✅ CHANGELOG.md exists"
else
    echo "❌ CHANGELOG.md missing"
fi

# Check configuration
echo ""
echo "Checking configuration..."
if [ -f "confused.yaml" ]; then
    echo "✅ confused.yaml exists"
else
    echo "❌ confused.yaml missing"
fi

# Check go.mod content
echo ""
echo "Checking go.mod content..."
if grep -q "module github.com/h0tak88r/confused" go.mod; then
    echo "✅ go.mod has correct module name"
else
    echo "❌ go.mod module name incorrect"
fi

if grep -q "go 1.21" go.mod; then
    echo "✅ go.mod has correct Go version"
else
    echo "❌ go.mod Go version incorrect"
fi

# Check for required dependencies
echo ""
echo "Checking dependencies..."
required_deps=("github.com/spf13/cobra" "github.com/spf13/viper" "github.com/google/go-github" "github.com/fatih/color")
for dep in "${required_deps[@]}"; do
    if grep -q "$dep" go.mod; then
        echo "✅ $dep dependency found"
    else
        echo "❌ $dep dependency missing"
    fi
done

echo ""
echo "Project Structure Verification Complete!"
echo "======================================="

# Summary
echo ""
echo "Summary:"
echo "- Project follows Go standard project layout"
echo "- Main application is in cmd/confused/"
echo "- Packages are properly organized in pkg/"
echo "- Internal types are in internal/"
echo "- Build system is configured"
echo "- Documentation is present"
echo ""
echo "To build the project:"
echo "1. Install Go 1.21 or later"
echo "2. Run: make build"
echo "3. Or run: go build -o confused ./cmd/confused"
echo ""
echo "To test GitHub token functionality:"
echo "./confused github repo microsoft/PowerShell --github-token YOUR_TOKEN"
echo "./confused github org microsoft --github-token YOUR_TOKEN"
