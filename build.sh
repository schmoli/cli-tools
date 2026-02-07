#!/bin/bash
set -e

cd "$(dirname "$0")"
mkdir -p bin

VERSION="${VERSION:-dev}"
LDFLAGS="-s -w -X main.version=${VERSION}"

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
    go build -ldflags "${LDFLAGS}" -o bin/${name} ./${dir}/cmd/${name}
    SIZE=$(ls -lh bin/${name} | awk '{print $5}')
    echo -e "${GREEN}✓ Built${NC} (${SIZE})"
    echo ""
}

build_tool "portainer-cli" "portainer"

if [ -f "nproxy/cmd/nproxy-cli/main.go" ]; then
    build_tool "nproxy-cli" "nproxy"
fi

if [ -f "trans/cmd/trans-cli/main.go" ]; then
    build_tool "trans-cli" "trans"
fi

if [ -f "pve/cmd/pve-cli/main.go" ]; then
    build_tool "pve-cli" "pve"
fi

if [ -f "abs/cmd/abs-cli/main.go" ]; then
    build_tool "abs-cli" "abs"
fi

if [ -f "sonarr/cmd/sonarr-cli/main.go" ]; then
    build_tool "sonarr-cli" "sonarr"
fi

if [ -f "radarr/cmd/radarr-cli/main.go" ]; then
    build_tool "radarr-cli" "radarr"
fi

echo -e "${GREEN}Done.${NC}"
