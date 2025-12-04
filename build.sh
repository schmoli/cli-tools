#!/bin/bash
set -e

cd "$(dirname "$0")"
mkdir -p bin

RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

build_tool() {
    local name=$1
    local dir=$2

    echo -e "${BLUE}=== ${name} ===${NC}"

    echo "Testing..."
    TEST_OUTPUT=$(go test ./${dir}/... -v 2>&1)
    if [ $? -eq 0 ]; then
        PASSED=$(echo "$TEST_OUTPUT" | grep -c -- '--- PASS' || echo 0)
        echo -e "${GREEN}✓ ${PASSED} passed${NC}"
    else
        echo "$TEST_OUTPUT" | grep -E '(FAIL|---\s*FAIL|panic:)'
        echo -e "${RED}✗ Tests failed${NC}"
        exit 1
    fi

    echo "Building..."
    go build -o bin/${name} ./${dir}/cmd/${name}
    SIZE=$(ls -lh bin/${name} | awk '{print $5}')
    echo -e "${GREEN}✓ Built${NC} (${SIZE})"
    echo ""
}

build_tool "portainer-cli" "portainer"

if [ -f "nproxy/cmd/nproxy-cli/main.go" ]; then
    build_tool "nproxy-cli" "nproxy"
fi

echo -e "${GREEN}Done.${NC}"
