#!/usr/bin/env bash
set -euo pipefail

# --- Helper Functions ---
log() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')] $*"
}

command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# --- Check Required Commands ---
for cmd in curl tar uname; do
  if ! command_exists "$cmd"; then
    log "Error: Required command '$cmd' is not installed. Exiting."
    exit 1
  fi
done

# --- Detect OS and Architecture ---
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    arm*)
        ARCH="arm"
        ;;
    *)
        log "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# --- Determine Asset Name and Version ---
# Assume VERSION is provided as an env var (or read from a file)
if [ -z "${VERSION:-}" ]; then
  VERSION="dev"
fi
ASSET="sailfin-${OS}-${ARCH}.tar.gz"
URL="https://github.com/SailfinIO/agent/releases/download/v${VERSION}/${ASSET}"

log "Detected OS: ${OS}"
log "Detected architecture: ${ARCH}"
log "Downloading asset: ${ASSET} (version ${VERSION}) from ${URL}"

# --- Download the Asset ---
curl -LO "$URL"

# --- Extract the Tarball ---
log "Extracting ${ASSET}..."
tar -xzvf "$ASSET"

# --- Locate the Binary ---
if [ -f sailfin ]; then
    BINARY="sailfin"
elif [ -f bin/sailfin ]; then
    BINARY="bin/sailfin"
else
    log "Error: 'sailfin' binary not found in the extracted files."
    exit 1
fi

# --- Install the Binary to a version-specific directory ---
INSTALL_DIR=${INSTALL_DIR:-/usr/local/sailfin/versions}
TARGET_VERSION_DIR="${INSTALL_DIR}/${VERSION}"
sudo mkdir -p "${TARGET_VERSION_DIR}"
log "Installing sailfin binary to ${TARGET_VERSION_DIR}..."
sudo mv "$BINARY" "${TARGET_VERSION_DIR}/sailfin"
sudo chmod +x "${TARGET_VERSION_DIR}/sailfin"

# --- Create a global symlink if desired ---
GLOBAL_BIN_DIR=${GLOBAL_BIN_DIR:-/usr/local/bin}
GLOBAL_LINK="${GLOBAL_BIN_DIR}/sailfin"
# Remove the old symlink if it exists.
if [ -L "$GLOBAL_LINK" ]; then
  sudo rm "$GLOBAL_LINK"
fi
# Create new symlink pointing to the installed version.
sudo ln -s "${TARGET_VERSION_DIR}/sailfin" "$GLOBAL_LINK"
log "Global symlink updated: ${GLOBAL_LINK} -> ${TARGET_VERSION_DIR}/sailfin"

# --- Create Configuration Directory (as before) ---
CONFIG_DIR="$HOME/.sailfin"
if [ ! -d "$CONFIG_DIR" ]; then
    log "Creating configuration directory at ${CONFIG_DIR}..."
    mkdir -p "$CONFIG_DIR"
else
    log "Configuration directory already exists at ${CONFIG_DIR}."
fi

CONFIG_FILE="${CONFIG_DIR}/AgentConfig.pkl"
if [ ! -f "$CONFIG_FILE" ]; then
    log "Installing default configuration to ${CONFIG_FILE}..."
    SAMPLE_URL="https://raw.githubusercontent.com/SailfinIO/agent/main/pkl/AgentConfig.pkl.sample"
    curl -sL "$SAMPLE_URL" -o "$CONFIG_FILE"
else
    log "Configuration file already exists at ${CONFIG_FILE}. Skipping default installation."
fi

# --- Cleanup ---
log "Cleaning up downloaded asset..."
rm "$ASSET"

log "Installation complete."
log "Sailfin binary for version ${VERSION} is installed at ${TARGET_VERSION_DIR}/sailfin."
log "Global symlink set at ${GLOBAL_LINK}."
