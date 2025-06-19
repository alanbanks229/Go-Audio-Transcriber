#!/bin/bash
set -e

APP_NAME="Go-Audio-Transcriber"
OUT_DIR="./dist"
MAIN_FILE="main.go"

mkdir -p "$OUT_DIR"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

echo "Building for: $OS ($ARCH)"

if ! command -v fyne &> /dev/null; then
  echo "Fyne CLI not found. Falling back to go build..."
  go build -o "$OUT_DIR/$APP_NAME" "$MAIN_FILE"
else
  echo "Using fyne package..."
  # use icon.png
  fyne package --source-dir . \
    --os "$OS" \
    --icon Look_Its_Me_Alan.png \
    --name "$APP_NAME" \
    --app-version 0.9.0 \
    --app-build 0

  if [ "$OS" = "darwin" ]; then
    # Ensure assets and bin get bundled inside the macOS app package
    cp -R assets "$APP_NAME.app/Contents/Resources/"
    mv "$APP_NAME.app" "$OUT_DIR/"
  else
    mv "$APP_NAME" "$OUT_DIR/"
  fi
fi

echo "✅ Build complete: $OUT_DIR/$APP_NAME"
