#!/usr/bin/env bash
# WIP Don't run this yet
set -e

echo "→ Installing Go deps"
go get fyne.io/fyne/v2@latest
go install fyne.io/tools/cmd/fyne@latest
go mod tidy

echo "→ Installing yt-dlp & ffmpeg"
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
sudo apt update
sudo apt install ffmpeg python3 python3-pip -y
pip3 install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu118
pip3 install git+https://github.com/openai/whisper.git


echo "✓ Done."