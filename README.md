# 🎙️ Audio Transcription Tool

A cross-platform desktop app built with **Go + Fyne** to convert **YouTube videos** or **MP3 files** into **text transcripts** using OpenAI’s Whisper (via `whisper.cpp`).

---

## 🚀 Features/Screenshot

- Clean UI with time-range slider
- Download & transcribe YouTube audio
- Transcribe local MP3 files
- Runs fully offline (after setup)
- Generates plain `.txt` transcripts
- Cross-platform: macOS, Linux, Windows
- Uses [`whisper.cpp`](https://github.com/ggerganov/whisper.cpp) under the hood

![image](https://github.com/user-attachments/assets/dce6d0e7-3119-4922-b740-44bf69217dfa)

---

## ⚙️ Prerequisites

- **Go 1.24.4+** (matches `go.mod`)
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
2. Extract so that `C:\mingw64\bin\gcc.exe` and `C:\mingw64\bin\g++.exe` exist.
3. Add `C:\mingw64\bin` (or your actual toolchain's `bin` directory) to your **user PATH**.
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

Run these commands in **PowerShell from the repository root**. Windows setup builds whisper.cpp from source; it does not download a ready-made Whisper executable.

##### NVIDIA GPU build (strongly recommended on Windows)

**GPU acceleration is strongly recommended for Windows machines with a compatible NVIDIA GPU.** Use the CPU-only option below if CUDA is unavailable or when troubleshooting. These CUDA instructions apply to Windows only.

For **NVIDIA GPU acceleration**, install these prerequisites first:

- A supported NVIDIA driver; check your GPU with `nvidia-smi`.
- [NVIDIA CUDA Toolkit](https://developer.nvidia.com/cuda-downloads), with `nvcc.exe` on PATH. The CUDA version shown by `nvidia-smi` describes driver support; check the installed toolkit with `nvcc --version`.
- [Visual Studio Build Tools](https://visualstudio.microsoft.com/downloads/) with **Desktop development with C++**, the x64/x86 MSVC compiler, and a Windows SDK. Use a compiler supported by your CUDA version; CUDA 13.2 with Visual Studio 2019 Build Tools 16.11 was tested successfully. See [NVIDIA's compiler compatibility table](https://docs.nvidia.com/cuda/archive/13.2.0/cuda-installation-guide-microsoft-windows/index.html).
- [CMake](https://cmake.org/download/) 3.24+ on PATH, with support for your Visual Studio generator. MinGW is still required for building the Go/Fyne GUI; the CUDA Whisper build uses MSVC instead.

Reopen PowerShell after installing tools or changing PATH, then check:

```powershell
git --version
go version
gcc --version
g++ --version
cmake --version
nvcc --version
nvidia-smi
.\scripts\download-assets.ps1 -Backend CUDA
```

If PowerShell says scripts are disabled, use `Set-ExecutionPolicy -Scope Process -ExecutionPolicy RemoteSigned` for the current session. For reviewed scripts extracted from an internet-downloaded ZIP, `Get-ChildItem .\scripts\*.ps1 | Unblock-File` may also be necessary. Organization-managed policies can override the session setting. Windows setup uses `curl.exe` (7.71.0+) and does not require Git Bash.

##### CPU-only build on Windows

For **CPU-only** Whisper, install CMake, MinGW-w64 (GCC and G++), and [Ninja](https://ninja-build.org/), put them on PATH, and run:

```powershell
.\scripts\download-assets.ps1 -Backend CPU
```

##### Rebuild or switch Whisper backends

Omitting `-Backend` selects CUDA when `nvcc.exe` is available, otherwise CPU with a warning. Explicit `-Backend CUDA` fails if GPU prerequisites are missing; it never silently substitutes a CPU build.

If the other assets already exist, build just Whisper without downloading them again:

```powershell
.\scripts\setup-whisper-cli.ps1 -Backend CUDA -SkipUpdate
```

To rebuild just Whisper **without CUDA**:

```powershell
.\scripts\setup-whisper-cli.ps1 -Backend CPU -SkipUpdate
```

For an **RTX 3080**, you can specify its CUDA architecture explicitly:

```powershell
.\scripts\setup-whisper-cli.ps1 -Backend CUDA -CudaArchitectures 86 -SkipUpdate
```

The default architecture is `native`, targeting the GPU visible during configuration. CPU and CUDA build trees are separate (`.whisper.cpp\build-cpu` and `.whisper.cpp\build-cuda`). If you change Visual Studio or CUDA versions and CMake reports a cached generator/toolset mismatch, rename the affected build directory to an unused backup name and retry. `-SkipUpdate` reuses the existing Whisper checkout; without it the helper updates that checkout first.

The first CUDA source build can take several minutes because it compiles many GPU kernels. Compiler output and newly written object files under `.whisper.cpp\build-cuda` indicate progress. Keeping that build directory allows subsequent runs to reuse compiled files.

The helper reuses `assets\ggml-small.en.bin`, installs `assets\whisper-cli.exe`, and copies the CUDA runtime DLLs (`cudart64*.dll`, `cublas64*.dll`, and `cublasLt64*.dll`) alongside it. Keep those DLLs with the executable. The previous executable is preserved as `assets\whisper-cli.backup-<id>.exe` for rollback. The GUI currently uses the `small.en` model even if another model is downloaded manually.

Verify actual GPU inference, rather than checking only that the app window opens:

```powershell
.\assets\whisper-cli.exe -m .\assets\ggml-small.en.bin -f .\.whisper.cpp\samples\jfk.wav
```

Look for CUDA initialization identifying your NVIDIA card and a CUDA GPU backend/buffer. While transcription runs, `nvidia-smi` should list `whisper-cli.exe`; seeing only the Go app there indicates GUI rendering, not GPU transcription. GPU use is enabled by default in CUDA Whisper builds; no additional Go build flag is needed. A job already running on the CPU keeps that backend; start a new job to use the installed CUDA build.

If a model download is interrupted, the helper leaves a `.partial` file and retries it on the next run. If an older script left a corrupt final `.bin` file, rename that file to an unused backup name before retrying. Do not download the model again just because inference is using the CPU.

##### Remove old Windows builds after verifying CUDA

First finish any transcription jobs and close the Go app. From the repository root, the following removes the legacy Whisper build directory and the CPU build cache, if present:

```powershell
foreach ($oldBuild in @('.\.whisper.cpp\build', '.\.whisper.cpp\build-cpu')) {
    if (Test-Path -LiteralPath $oldBuild) {
        Remove-Item -LiteralPath $oldBuild -Recurse -Force
    }
}
```

Once GPU inference works and you no longer need rollback, remove the saved old Whisper executables and Go's executable backup files:

```powershell
Get-ChildItem -LiteralPath .\assets -Filter 'whisper-cli.backup-*.exe' -File | Remove-Item
Get-ChildItem -LiteralPath . -Filter 'Go-Audio-Transcriber.exe~*' -File | Remove-Item
```

Keep `Go-Audio-Transcriber.exe`, `assets\whisper-cli.exe`, the model, FFmpeg/ffprobe/yt-dlp, and the CUDA runtime DLLs. Keep `.whisper.cpp\build-cuda` for fast incremental CUDA rebuilds; it is the current GPU build, not an old CPU build. None of these cleanup commands removes audio files or transcripts in Downloads.

✔️ Downloads:

- `yt-dlp`
- `ffmpeg`
- `ffprobe` (audio duration analysis)
- `whisper-cli`
- `ggml-small.en.bin` (Whisper model)

---

### 5. Build and Run

#### macOS/Linux

```bash
go build -o GoAudioTranscriber
./GoAudioTranscriber
```

#### Windows

```powershell
go build -o Go-Audio-Transcriber.exe .
if ($LASTEXITCODE -eq 0) { .\Go-Audio-Transcriber.exe }
```

Or run directly:

```powershell
go run .
```

Keep the repository root as the current directory when running the app. Expected Windows assets are `assets\yt-dlp.exe`, `assets\ffmpeg.exe`, `assets\ffprobe.exe`, `assets\whisper-cli.exe`, `assets\ggml-small.en.bin`, plus the runtime DLLs for a CUDA build. If rebuilding reports access denied, close the app or tool holding the target file open; use a writable checkout rather than routinely running as Administrator.

---

## 🗂 Relevant Folder/File Structure

```plaintext
Go-Audio-Transcriber/
├── main.go             # application entry point
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
