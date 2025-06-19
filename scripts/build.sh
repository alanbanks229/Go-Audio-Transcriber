#!/bin/bash
set -e

APP_NAME="Go-Audio-Transcriber"
OUT_DIR="./dist"
MAIN_FILE="./cmd/main.go"

mkdir -p "$OUT_DIR"

OS="$(uname -s)"
ARCH="$(uname -m)"

echo "Building for: $OS ($ARCH)"

if ! command -v fyne &> /dev/null; then
  echo "Fyne CLI not found. Falling back to go build..."
  go build -o "$OUT_DIR/$APP_NAME" "$MAIN_FILE"
else
  echo "Using fyne package..."
  fyne package -os "$OS" -name "$APP_NAME" -icon "./icon.png" -appVersion "1.0.0"
  mv "$APP_NAME" "$OUT_DIR/"
fi

echo "Build complete: $OUT_DIR"
