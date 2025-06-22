#!/bin/bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# --- Setup target directories ---
WHISPER_REPO="$ROOT_DIR/.whisper.cpp"
BIN_DIR="$ROOT_DIR/assets/binaries"
MODEL_DIR="$ROOT_DIR/assets/models"

rm -rf $WHISPER_REPO
rm -rf $BIN_DIR
rm -rf $MODEL_DIR

mkdir -p "$BIN_DIR"
mkdir -p "$MODEL_DIR"

OS="$(uname -s)"
ARCH="$(uname -m)"

### -------- yt-dlp -------- ###
YTDLP_URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
if [[ "$OS" == "Darwin" ]]; then
    YTDLP_URL="${YTDLP_URL}_macos"
fi

echo "⬇️ Downloading yt-dlp..."
curl -fLo "$BIN_DIR/yt-dlp" "$YTDLP_URL"
chmod +x "$BIN_DIR/yt-dlp"

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

    echo "🔗 Linking ffmpeg into $BIN_DIR"
    ln -sf "$(which ffmpeg)" "$BIN_DIR/ffmpeg"

elif [[ "$OS" == "Linux" ]]; then
    if ! command -v ffmpeg >/dev/null 2>&1; then
        echo "⚠️  Please install ffmpeg using your package manager (e.g., sudo apt install ffmpeg)"
        exit 1
    else
        echo "✅ ffmpeg already installed"
    fi

    echo "🔗 Linking ffmpeg into $BIN_DIR"
    ln -sf "$(which ffmpeg)" "$BIN_DIR/ffmpeg"
fi

### -------- whisper-cli (build from source) -------- ###
"$ROOT_DIR/scripts/setup-whisper-cli.sh"

echo "✅ All assets downloaded. You're ready to Run/Build the app!"
