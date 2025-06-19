#!/bin/bash
set -e

ASSETS_DIR="./assets"
mkdir -p "$ASSETS_DIR"

OS="$(uname -s)"
ARCH="$(uname -m)"

# yt-dlp
YTDLP_URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp"
if [[ "$OS" == "Darwin" ]]; then
    cp_cmd="cp"
    YTDLP_URL="${YTDLP_URL}_macos"
elif [[ "$OS" == "Linux" ]]; then
    cp_cmd="cp"
    YTDLP_URL="${YTDLP_URL}"
else
    echo "Unsupported OS: $OS"
    exit 1
fi
curl -Lo "$ASSETS_DIR/yt-dlp" "$YTDLP_URL"
chmod +x "$ASSETS_DIR/yt-dlp"

# ffmpeg (macOS/Intel)
if [[ "$OS" == "Darwin" ]]; then
    FFMPEG_URL="https://evermeet.cx/ffmpeg/ffmpeg.zip"
    curl -Lo "$ASSETS_DIR/ffmpeg.zip" "$FFMPEG_URL"
    unzip -o "$ASSETS_DIR/ffmpeg.zip" -d "$ASSETS_DIR"
    rm "$ASSETS_DIR/ffmpeg.zip"
elif [[ "$OS" == "Linux" ]]; then
    echo "Please install ffmpeg via your package manager (e.g., apt install ffmpeg)"
fi

# whisper-cli binary and model
WHISPER_URL="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/bin/whisper-${OS,,}-${ARCH}"
MODEL_URL="https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.en.bin"

curl -Lo "$ASSETS_DIR/whisper-cli" "$WHISPER_URL"
chmod +x "$ASSETS_DIR/whisper-cli"

curl -Lo "$ASSETS_DIR/ggml-small.en.bin" "$MODEL_URL"