# 🎙️ Audio Transcription Tool

A cross-platform desktop app built with **Go + Fyne** to convert **YouTube videos** or **MP3 files** into **text transcripts** using OpenAI’s Whisper (via `whisper.cpp`).

---

## 🚀 Features

- Clean UI with time-range slider
- Download & transcribe YouTube audio
- Transcribe local MP3 files
- Runs fully offline (after setup)
- Generates plain `.txt` transcripts
- Cross-platform: macOS, Linux, Windows
- Uses [`whisper.cpp`](https://github.com/ggerganov/whisper.cpp) under the hood

---

## ⚙️ Prerequisites

- **Go 1.18+**
- **Git**
- **Internet** (for setup scripts)
- **C Compiler** (required for GUI build):
  - **Windows**: [MinGW-w64](https://winlibs.com/)
  - **macOS/Linux**: Usually pre-installed

---

## 📦 Installation Guide

### 1. Install Go

[Download Go](https://go.dev/dl/) and follow installer prompts.  
Verify with:

```bash
go version
```

---

### 2. (Windows Only) Install MinGW-w64

> macOS/Linux users can skip this step.

1. Download from [winlibs.com](https://winlibs.com/)
2. Extract to `C:\mingw64`
3. Add `C:\mingw64\...mingw64\bin` to **System PATH**
4. Open **new** PowerShell and verify with:

```powershell
gcc -v
```

You should see `x86_64-w64-mingw32`.

---

### 3. Clone This Repo

```bash
git clone https://github.com/alanbanks229/Go-Audio-Transcriber.git
cd Go-Audio-Transcriber
```

---

### 4. Download Required Binaries

These are used for downloading, trimming, and transcribing audio.

#### macOS/Linux

```bash
chmod +x scripts/download-assets.sh
./scripts/download-assets.sh
```

#### Windows

```powershell
.\scripts\download-assets.ps1
```

✔️ Downloads:

- `yt-dlp`
- `ffmpeg`
- `whisper-cli`
- `ggml-small.en.bin` (Whisper model)

---

### 5. Build and Run

#### macOS/Linux

```bash
go build -o GoAudioTranscriber ./cmd
./GoAudioTranscriber
```

#### Windows

```powershell
go build -o GoAudioTranscriber.exe ./cmd
.\GoAudioTranscriber.exe
```

Or run directly:

```bash
go run ./cmd/main.go
```

---

## 🗂 Relevant Folder/File Structure

```plaintext
Go-Audio-Transcriber/
├── cmd/                # main.go entry point
├── dist/               # executable file built from scripts
│
├── internal/
│   ├── audio/          # Logic for audio processing (ffmpeg, trimming, etc.)
│   ├── transcribe/     # Core transcription logic (whisper.cpp integration)
│   ├── ui/             # Fyne UI widgets and layout logic
│   ├── util/           # General utility functions (OS detection, binary path helpers)
│   └── youtube/        # YouTube video/audio downloading (yt-dlp integration)
│
├── assets/             # Downloaded external binaries (yt-dlp, ffmpeg, whisper-cli) and Whisper AI
├── scripts/            # Helper scripts for setup and building (.sh and .ps1)
```

---

## 📝 License

Licensed under the Apache 2.0 License  
© 2025 Alan Kim Banks
