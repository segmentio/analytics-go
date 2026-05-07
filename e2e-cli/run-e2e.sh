#!/bin/bash
#
# Run E2E tests for analytics-go
#
# Prerequisites: Node.js 18+ and Go 1.17+
#
# Usage:
#   ./run-e2e.sh [extra args passed to run-tests.sh]
#
# Override sdk-e2e-tests location:
#   E2E_TESTS_DIR=../my-e2e-tests ./run-e2e.sh
#

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SDK_ROOT="$SCRIPT_DIR/.."
E2E_DIR="${E2E_TESTS_DIR:-$SDK_ROOT/../sdk-e2e-tests}"

echo "=== Building analytics-go e2e-cli ==="
echo "Using Go: $(go version)"

# Build the CLI binary
cd "$SCRIPT_DIR"
go build -o e2e-cli-bin ./...

echo ""

# Run tests
cd "$E2E_DIR"
./scripts/run-tests.sh \
    --sdk-dir "$SCRIPT_DIR" \
    --cli "$SCRIPT_DIR/e2e-cli-bin" \
    "$@"
