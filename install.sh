
set -e

REPO="carlosarraes/taildog"
BINARY_NAME="taildog"
BIN_DIR="${BIN_DIR:-$HOME/.local/bin}"
GITHUB_LATEST="https://api.github.com/repos/${REPO}/releases/latest"

get_arch() {
  ARCH=$(uname -m)
  case $ARCH in
  x86_64) ARCH="x86_64" ;;
  aarch64) ARCH="aarch64" ;;
  arm64) ARCH="aarch64" ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
  esac
}

get_os() {
  OS=$(uname -s)
  case $OS in
  Linux) OS="linux" ;;
  Darwin) OS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
  esac
}

download_binary() {
  echo "Fetching latest release..."
  VERSION=$(curl -s $GITHUB_LATEST | grep -o '"tag_name": "[^"]*' | cut -d'"' -f4)
  if [ -z "$VERSION" ]; then
    echo "Failed to fetch latest version"
    exit 1
  fi

  echo "Latest version: $VERSION"

  TMP_DIR=$(mktemp -d)
  trap 'rm -rf "$TMP_DIR"' EXIT

  echo "Downloading ${BINARY_NAME} ${VERSION}..."

  BINARY_SUFFIX="${BINARY_NAME}-${OS}-${ARCH}"
  
  if [ "$OS" = "darwin" ] && [ "$ARCH" = "x86_64" ]; then
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}"
    echo "Downloading from: $DOWNLOAD_URL"
    if ! curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${BINARY_NAME}" 2>/dev/null; then
      echo "Generic binary not found, trying platform-specific binary..."
      DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_SUFFIX}"
      echo "Downloading from: $DOWNLOAD_URL"
      curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${BINARY_NAME}" || {
        echo "Download failed. Check URL/permissions/network."
        exit 1
      }
    fi
  else
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_SUFFIX}"
    echo "Downloading from: $DOWNLOAD_URL"
    curl -fsSL "$DOWNLOAD_URL" -o "${TMP_DIR}/${BINARY_NAME}" || {
      echo "Download failed. Check URL/permissions/network."
      exit 1
    }
  fi

  chmod +x "${TMP_DIR}/${BINARY_NAME}"

  CREATED_DIR_MSG=""
  if [ ! -d "$BIN_DIR" ]; then
    echo "Installation directory '$BIN_DIR' not found."
    echo "Creating directory: $BIN_DIR"
    mkdir -p "$BIN_DIR"
    CREATED_DIR_MSG="Note: Created directory '$BIN_DIR'. You might need to add it to your system's PATH."
  fi

  echo "Installing to $BIN_DIR..."
  install -m 755 "${TMP_DIR}/${BINARY_NAME}" "$BIN_DIR"


  echo "${BINARY_NAME} ${VERSION} installed successfully to $BIN_DIR"

  if [ -n "$CREATED_DIR_MSG" ]; then
    echo ""
    echo "$CREATED_DIR_MSG"
  fi
}

get_arch
get_os
download_binary

echo ""
echo "Installation complete! Run '${BINARY_NAME} --help' to get started."
echo "Example usage: ${BINARY_NAME} \"service:my-app\""
