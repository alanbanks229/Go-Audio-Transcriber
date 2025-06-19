# PowerShell Setup Script for Windows
$AssetsDir = "./assets"
New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null

Write-Host "Downloading yt-dlp.exe..."
Invoke-WebRequest -Uri "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" -OutFile "$AssetsDir/yt-dlp.exe"

Write-Host "Downloading ffmpeg zip..."
Invoke-WebRequest -Uri "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" -OutFile "$AssetsDir/ffmpeg.zip"
Expand-Archive -Path "$AssetsDir/ffmpeg.zip" -DestinationPath "$AssetsDir" -Force
Remove-Item "$AssetsDir/ffmpeg.zip"

# Move ffmpeg.exe
$ffmpegBin = Get-ChildItem -Path $AssetsDir -Recurse -Filter "ffmpeg.exe" | Select-Object -First 1
if ($ffmpegBin) {
    Move-Item $ffmpegBin.FullName "$AssetsDir/ffmpeg.exe" -Force
}

# Move ffprobe.exe (needed for duration analysis)
$ffprobeBin = Get-ChildItem -Path $AssetsDir -Recurse -Filter "ffprobe.exe" | Select-Object -First 1
if ($ffprobeBin) {
    Move-Item $ffprobeBin.FullName "$AssetsDir/ffprobe.exe" -Force
}

# Download whisper-cli and model
try {
    Write-Host "Downloading whisper-cli.exe..."
    Invoke-WebRequest -Uri "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/bin/whisper-windows-amd64.exe" -OutFile "$AssetsDir/whisper-cli.exe"
    Write-Host "Downloading whisper model..."
    Invoke-WebRequest -Uri "https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.en.bin" -OutFile "$AssetsDir/ggml-small.en.bin"
}
catch {
    Write-Warning "Whisper binaries failed to download (Hugging Face might be temporarily down). Please try again later or download manually."
}
