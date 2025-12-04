#!/bin/bash
set -e

# Test install.sh in a fresh Linux container

IMAGE="${1:-ubuntu:22.04}"
echo "Testing install in ${IMAGE}..."

docker run --rm "${IMAGE}" bash -c '
set -e

# Install prerequisites
if command -v apt-get &>/dev/null; then
    apt-get update -qq && apt-get install -y -qq curl ca-certificates
elif command -v apk &>/dev/null; then
    apk add --no-cache curl bash ca-certificates
fi

echo "=== Running install script ==="
curl -fsSL https://raw.githubusercontent.com/schmoli/cli-tools/main/install.sh | bash

echo ""
echo "=== Verifying binaries ==="
export PATH="$HOME/.local/bin:$PATH"

echo -n "portainer-cli: "
portainer-cli --version

echo -n "nproxy-cli: "
nproxy-cli --version

echo ""
echo "=== Install successful ==="
'
