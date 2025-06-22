#!/usr/bin/env bash
set -e

# ----------- Config ------------
MODEL="base.en"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WHISPER_REPO="https://github.com/ggml-org/whisper.cpp.git"
WHISPER_DIR="$ROOT_DIR/.whisper.cpp"

# Allow override of output dir via WHISPER_ASSETS_BIN_DIR env var
ASSETS_BIN_DIR="$ROOT_DIR/assets/binaries"
mkdir -p "$ASSETS_BIN_DIR"
# -------------------------------

# Detect OS and ARCH
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin)   PLATFORM_OS="macos" ;;
  Linux)    PLATFORM_OS="linux" ;;
  MINGW*|MSYS*|CYGWIN*) PLATFORM_OS="windows" ;;
  *) echo "❌ Unsupported OS: $OS" && exit 1 ;;
esac

case "$ARCH" in
  arm64|aarch64)     PLATFORM_ARCH="arm64" ;;
  x86_64|amd64)      PLATFORM_ARCH="x64" ;;
  *) echo "❌ Unsupported architecture: $ARCH" && exit 1 ;;
esac

WHISPER_BUILD_NAME="whisper-${PLATFORM_OS}-${PLATFORM_ARCH}"
[[ "$PLATFORM_OS" == "windows" ]] && WHISPER_BUILD_NAME="${WHISPER_BUILD_NAME}.exe"

# Step 1: Clone whisper.cpp if not already present
#         and CD into directory
if [ ! -d "$WHISPER_DIR" ]; then
  echo "📥 Cloning whisper.cpp..."
  git clone "$WHISPER_REPO" "$WHISPER_DIR"
else
  echo "🔄 Updating whisper.cpp..."
  git -C "$WHISPER_DIR" pull
fi
cd "$WHISPER_DIR"


# Step 2: Download model if not present
if [ ! -f "models/ggml-${MODEL}.bin" ]; then
  echo "Downloading model: $MODEL..."
  sh ./models/download-ggml-model.sh "$MODEL"
  mv "./models/ggml-${MODEL}.bin" "$ROOT_DIR/assets/models"
fi

# Step 3: Build whisper-cli
echo "🔨 Building whisper-cli..."
cmake -B build
cmake --build build --config Release

# Step 4: Copy built binary to ASSETS_BIN_DIR
BIN_SOURCE="build/bin/whisper-cli"
BIN_DEST="$ASSETS_BIN_DIR/whisper-cli"

if [ ! -f "$BIN_SOURCE" ]; then
  echo "❌ Build failed: $BIN_SOURCE not found"
  exit 1
fi

mv  "$BIN_SOURCE" "$BIN_DEST"
chmod +x "$BIN_DEST"

echo "✅ whisper-cli is ready: $BIN_DEST"
