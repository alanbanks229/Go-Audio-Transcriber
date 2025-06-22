#!/bin/bash
set -e

APP_NAME="Go-Audio-Transcriber"
BUILD_DIR="$PWD/build"
ASSETS_DIR="$PWD/assets"
ICON="$PWD/icon.png"
MAIN_DIR="$PWD"

mkdir -p "$BUILD_DIR"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

echo "🔨 Building for: $OS ($ARCH)"

if ! command -v fyne &> /dev/null; then
  echo "⚠️  Fyne CLI not found. Falling back to 'go build'..."
  go build -o "$BUILD_DIR/$APP_NAME" "$MAIN_DIR/main.go"
  exit 0
fi

# Use Fyne CLI to build from root
pushd "$MAIN_DIR" > /dev/null

fyne package \
  --os "$OS" \
  --icon "$ICON" \
  --name "$APP_NAME" \
  --app-version 0.9.0 \
  --app-build 0

popd > /dev/null

# Handle output format per OS
if [[ "$OS" == "darwin" ]]; then
  echo "📁 Copying .app bundle and assets..."
  cp -R "$ASSETS_DIR" "$APP_NAME.app/Contents/Resources/"
  mv "$APP_NAME.app" "$BUILD_DIR/"

elif [[ "$OS" == "linux" ]]; then
  echo "📦 Extracting .tar.xz package to build..."
  tar -xf "$APP_NAME.tar.xz" -C "$BUILD_DIR"
  rm -rf "$APP_NAME.tar.xz"

elif [[ "$OS" == "windows_nt" || "$OS" == "mingw"* || "$OS" == "msys"* ]]; then
  echo "📦 Moving .zip package..."
  mv "$APP_NAME.zip" "$BUILD_DIR/"
fi

echo "✅ Build complete! Output is in $BUILD_DIR"
