#!/bin/bash

# Examples script for Confused Dependency Confusion Scanner
# This script demonstrates various usage patterns and features

set -e

echo "Confused Dependency Confusion Scanner - Examples"
echo "================================================"

# Check if confused binary exists
if [ ! -f "./confused" ]; then
    echo "Error: confused binary not found. Please run 'make build' first."
    exit 1
fi

echo ""
echo "1. Basic Local File Scanning"
echo "----------------------------"

# Create sample package.json
cat > sample-package.json << 'EOF'
{
  "name": "sample-app",
  "version": "1.0.0",
  "dependencies": {
    "lodash": "^4.17.21",
    "express": "^4.18.2",
    "internal-package": "^1.0.0",
    "@company/private-package": "^2.0.0"
  },
  "devDependencies": {
    "jest": "^29.0.0",
    "internal-dev-package": "^1.0.0"
  }
}
EOF

echo "Scanning sample package.json..."
./confused scan sample-package.json -v

echo ""
echo "2. Scan with Safe Spaces"
echo "------------------------"

echo "Scanning with safe spaces for @company/*..."
./confused scan sample-package.json --safe-spaces "@company/*" -v

echo ""
echo "3. Different Package Managers"
echo "-----------------------------"

# Create sample requirements.txt
cat > sample-requirements.txt << 'EOF'
requests==2.31.0
flask==2.3.3
internal-python-package==1.0.0
company-utils==2.1.0
EOF

echo "Scanning Python requirements.txt..."
./confused scan -l pip sample-requirements.txt -v

# Create sample composer.json
cat > sample-composer.json << 'EOF'
{
  "name": "sample/php-app",
  "require": {
    "monolog/monolog": "^3.0",
    "internal/php-package": "^1.0",
    "company/utilities": "^2.0"
  }
}
EOF

echo "Scanning PHP composer.json..."
./confused scan -l composer sample-composer.json -v

echo ""
echo "4. JSON Output Format"
echo "---------------------"

echo "Generating JSON report..."
./confused scan sample-package.json -f json -o sample-report.json
echo "JSON report saved to: sample-report.json"

echo ""
echo "5. Configuration Management"
echo "---------------------------"

echo "Generating sample configuration..."
./confused config generate -o confused.example.yaml
echo "Sample configuration saved to: confused.example.yaml"

echo "Validating configuration..."
./confused config validate

echo ""
echo "6. Worker Pool Testing"
echo "----------------------"

echo "Testing with different worker counts..."
echo "Workers: 1"
./confused scan sample-package.json -w 1 --format json -o test-1-worker.json

echo "Workers: 5"
./confused scan sample-package.json -w 5 --format json -o test-5-workers.json

echo "Workers: 20"
./confused scan sample-package.json -w 20 --format json -o test-20-workers.json

echo ""
echo "7. Multiple Language Scanning"
echo "-----------------------------"

echo "Scanning multiple files with different languages..."
./confused scan sample-package.json -l npm -f json -o npm-results.json &
./confused scan sample-requirements.txt -l pip -f json -o pip-results.json &
./confused scan sample-composer.json -l composer -f json -o composer-results.json &
wait

echo "All scans completed!"

echo ""
echo "8. Batch Processing Example"
echo "---------------------------"

echo "Processing multiple files in batch..."
for file in sample-*.json sample-*.txt; do
    if [ -f "$file" ]; then
        echo "Processing: $file"
        ./confused scan "$file" -f json -o "batch-${file%.*}.json" 2>/dev/null || true
    fi
done

echo ""
echo "9. Performance Testing"
echo "----------------------"

echo "Testing performance with timing..."
time ./confused scan sample-package.json -w 10 -f json -o performance-test.json

echo ""
echo "10. Error Handling"
echo "------------------"

echo "Testing error handling with non-existent file..."
./confused scan non-existent-file.json 2>/dev/null || echo "Error handled gracefully"

echo "Testing error handling with invalid language..."
./confused scan sample-package.json -l invalid 2>/dev/null || echo "Invalid language error handled gracefully"

echo ""
echo "11. Cleanup"
echo "-----------"

echo "Cleaning up example files..."
rm -f sample-*.json sample-*.txt
rm -f sample-report.json
rm -f test-*-workers.json
rm -f *-results.json
rm -f batch-*.json
rm -f performance-test.json
rm -f confused.example.yaml

echo "Cleanup complete!"

echo ""
echo "Examples completed successfully!"
echo "==============================="
echo ""
echo "For more advanced examples, see:"
echo "- GitHub scanning: ./confused github --help"
echo "- Web scanning: ./confused web --help"
echo "- Configuration: ./confused config --help"
echo ""
echo "For integration with AutoAR framework, see the README.md"
