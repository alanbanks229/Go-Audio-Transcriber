#!/bin/bash
set -euo pipefail

ASSETS_DIR="./assets"
mkdir -p "$ASSETS_DIR"

OS="$(uname -s)"
ARCH="$(uname -m)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

### -------- yt-dlp -------- ###
if [[ "$OS" == "Darwin" ]]; then
    "$SCRIPT_DIR/update-ytdlp-macos.sh"
else
    YTDLP_URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
    echo "⬇️ Downloading yt-dlp..."
    curl -fLo "$ASSETS_DIR/yt-dlp" "$YTDLP_URL"
    chmod +x "$ASSETS_DIR/yt-dlp"
    echo "✅ yt-dlp updated: $("$ASSETS_DIR/yt-dlp" --version)"
fi

### -------- ffmpeg -------- ###
if [[ "$OS" == "Darwin" ]]; then
    if ! command -v ffmpeg >/dev/null 2>&1; then
        echo "⚠️  ffmpeg not found. Installing via Homebrew..."
        if command -v brew >/dev/null 2>&1; then
            brew install ffmpeg
        else
            echo "❌ Homebrew not found. Please install it from https://brew.sh"
            exit 1
        fi
    else
        echo "✅ ffmpeg already installed"
    fi

    echo "🔗 Linking ffmpeg into $ASSETS_DIR"
    ln -sf "$(which ffmpeg)" "$ASSETS_DIR/ffmpeg"

elif [[ "$OS" == "Linux" ]]; then
    if ! command -v ffmpeg >/dev/null 2>&1; then
        echo "⚠️  Please install ffmpeg using your package manager (e.g., sudo apt install ffmpeg)"
        exit 1
    else
        echo "✅ ffmpeg already installed"
    fi

    echo "🔗 Linking ffmpeg into $ASSETS_DIR"
    ln -sf "$(which ffmpeg)" "$ASSETS_DIR/ffmpeg"
fi

### -------- whisper-cli (build from source) -------- ###
echo "🧱 Building whisper-cli from source..."
"$SCRIPT_DIR/setup-whisper-cli.sh"

### -------- whisper model -------- ###
MODEL_URL="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.en.bin"
echo "⬇️ Downloading Whisper model (ggml-small.en.bin)..."
curl -fLo "$ASSETS_DIR/ggml-small.en.bin" "$MODEL_URL"

echo "If no errors occurred you can now Run/Build the APP! See README."
