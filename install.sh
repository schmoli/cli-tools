#!/bin/bash
set -e

REPO="schmoli/cli-tools"
INSTALL_DIR="${HOME}/.local/bin"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "${ARCH}" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm64)   ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac

case "${OS}" in
    linux|darwin) ;;
    *)
        echo "Unsupported OS: ${OS}"
        exit 1
        ;;
esac

echo "Detected: ${OS}-${ARCH}"

# Get latest release
echo "Fetching latest release..."
LATEST=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "${LATEST}" ]; then
    echo "Failed to fetch latest release. Check https://github.com/${REPO}/releases"
    exit 1
fi

echo "Latest version: ${LATEST}"

# Download
TARBALL="cli-tools-${OS}-${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${LATEST}/${TARBALL}"

echo "Downloading ${URL}..."
TMPDIR=$(mktemp -d)
curl -fsSL "${URL}" -o "${TMPDIR}/${TARBALL}"

# Extract
echo "Extracting..."
tar -xzf "${TMPDIR}/${TARBALL}" -C "${TMPDIR}"

# Install
mkdir -p "${INSTALL_DIR}"
mv "${TMPDIR}/portainer-cli" "${INSTALL_DIR}/"
mv "${TMPDIR}/nproxy-cli" "${INSTALL_DIR}/"
chmod +x "${INSTALL_DIR}/portainer-cli" "${INSTALL_DIR}/nproxy-cli"

# Cleanup
rm -rf "${TMPDIR}"

echo ""
echo "Installed to ${INSTALL_DIR}:"
echo "  - portainer-cli"
echo "  - nproxy-cli"

# Check PATH
if [[ ":${PATH}:" != *":${INSTALL_DIR}:"* ]]; then
    echo ""
    echo "Add to your PATH:"
    echo "  export PATH=\"${INSTALL_DIR}:\${PATH}\""
fi
