#!/usr/bin/env bash
set -e

# ----------- Config ------------
MODEL="base.en"
ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ASSETS_DIR="$ROOT_DIR/assets"
WHISPER_REPO="https://github.com/ggml-org/whisper.cpp.git"
WHISPER_DIR="$ROOT_DIR/.whisper.cpp"
# -------------------------------

mkdir -p "$ASSETS_DIR"

# Detect OS and ARCH
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin)
    PLATFORM_OS="macos"
    ;;
  Linux)
    PLATFORM_OS="linux"
    ;;
  MINGW*|MSYS*|CYGWIN*)
    PLATFORM_OS="windows"
    ;;
  *)
    echo "❌ Unsupported OS: $OS"
    exit 1
    ;;
esac

case "$ARCH" in
  arm64|aarch64)
    PLATFORM_ARCH="arm64"
    ;;
  x86_64|amd64)
    PLATFORM_ARCH="x64"
    ;;
  *)
    echo "❌ Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

WHISPER_BUILD_NAME="whisper-${PLATFORM_OS}-${PLATFORM_ARCH}"
[[ "$PLATFORM_OS" == "windows" ]] && WHISPER_BUILD_NAME="${WHISPER_BUILD_NAME}.exe"

# Step 1: Clone whisper.cpp
if [ ! -d "$WHISPER_DIR" ]; then
  echo "📥 Cloning whisper.cpp..."
  git clone "$WHISPER_REPO" "$WHISPER_DIR"
else
  echo "🔄 Updating whisper.cpp..."
  git -C "$WHISPER_DIR" pull
fi

cd "$WHISPER_DIR"

# Step 2: Download model if needed
if [ ! -f "models/ggml-${MODEL}.bin" ]; then
  echo "⬇️ Downloading model: $MODEL..."
  sh ./models/download-ggml-model.sh "$MODEL"
fi

# Step 3: Build whisper-cli
echo "🔨 Building whisper-cli..."
cmake -B build
cmake --build build --config Release

# Step 4: Copy to bin/
echo "📦 Moving binary to $ASSETS_DIR/$WHISPER_BUILD_NAME"
cp build/bin/whisper-cli "$ASSETS_DIR/$WHISPER_BUILD_NAME"
mv "$ASSETS_DIR/$WHISPER_BUILD_NAME" "$ASSETS_DIR/whisper-cli"
chmod +x "$ASSETS_DIR/whisper-cli"

echo "✅ whisper-cli is ready: $ASSETS_DIR/whisper-cli"
