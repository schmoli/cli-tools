#!/bin/bash
set -e

# Source cargo if available
[ -f "$HOME/.cargo/env" ] && source "$HOME/.cargo/env"

echo "Building Go..."
cd go && go build -o portainer-cli ./cmd/portainer-cli && cd ..

echo "Building Rust..."
cd rust && cargo build --release && cp target/release/portainer-cli . && cd ..

echo "Done."
echo "  go/portainer-cli"
echo "  rust/portainer-cli"
