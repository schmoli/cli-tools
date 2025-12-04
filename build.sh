#!/bin/bash
set -e

# Source cargo if available
[ -f "$HOME/.cargo/env" ] && source "$HOME/.cargo/env"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Go: portainer ===${NC}"
echo "Testing..."
cd go/portainer
GO_TEST_OUTPUT=$(go test ./... -v 2>&1)
GO_TEST_EXIT=$?
if [ $GO_TEST_EXIT -eq 0 ]; then
    GO_PASSED=$(echo "$GO_TEST_OUTPUT" | grep -c -- '--- PASS')
    echo -e "${GREEN}✓ ${GO_PASSED} passed${NC}"
else
    echo "$GO_TEST_OUTPUT" | grep -E '(FAIL|---\s*FAIL|panic:)'
    echo -e "${RED}✗ Tests failed${NC}"
    exit 1
fi

echo "Building..."
go build -o portainer-cli ./cmd/portainer-cli
GO_SIZE=$(ls -lh portainer-cli | awk '{print $5}')
echo -e "${GREEN}✓ Built${NC} (${GO_SIZE})"
cd ../..

echo ""
echo -e "${BLUE}=== Rust: portainer ===${NC}"
echo "Testing..."
cd rust
RUST_TEST_OUTPUT=$(cargo test 2>&1)
RUST_TEST_EXIT=$?
if [ $RUST_TEST_EXIT -eq 0 ]; then
    RUST_PASSED=$(echo "$RUST_TEST_OUTPUT" | grep -o '[0-9]* passed' | awk '{sum+=$1} END {print sum}')
    echo -e "${GREEN}✓ ${RUST_PASSED} passed${NC}"
else
    echo "$RUST_TEST_OUTPUT" | grep -E '(FAILED|panicked|error\[)'
    echo -e "${RED}✗ Tests failed${NC}"
    exit 1
fi

echo "Building..."
cargo build --release -p portainer-cli 2>&1 | grep -E '(Compiling|Finished|error)' | tail -5
cp target/release/portainer-cli portainer/
RUST_SIZE=$(ls -lh portainer/portainer-cli | awk '{print $5}')
echo -e "${GREEN}✓ Built${NC} (${RUST_SIZE})"
cd ..

echo ""
echo -e "${BLUE}=== Summary ===${NC}"
echo "  go/portainer/portainer-cli     ${GO_SIZE}"
echo "  rust/portainer/portainer-cli   ${RUST_SIZE}"
echo -e "${GREEN}Done.${NC}"
