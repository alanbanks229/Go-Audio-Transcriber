#!/bin/bash
set -euo pipefail

if [[ "$(uname -s)" != "Darwin" ]]; then
    echo "❌ update-ytdlp-macos.sh is only for macOS."
    echo "   Use scripts/download-assets.sh for cross-platform setup."
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
ASSETS_DIR="$PROJECT_ROOT/assets"
YTDLP_PATH="$ASSETS_DIR/yt-dlp"
YTDLP_URL="https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_macos"

mkdir -p "$ASSETS_DIR"

echo "⬇️ Downloading latest macOS yt-dlp..."
curl -fL -o "$YTDLP_PATH" "$YTDLP_URL"
chmod +x "$YTDLP_PATH"

echo "✅ yt-dlp updated: $("$YTDLP_PATH" --version)"
