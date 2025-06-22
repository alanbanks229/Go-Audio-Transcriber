$rootDir = Resolve-Path "."
$whisperRepo = "$rootDir\.whisper.cpp"
$binDir = "$rootDir\assets\binaries"
$modelDir = "$rootDir\assets\models"

# Clean up old binaries/models/repo
Write-Host "Cleaning old binaries/models/repo..."
Remove-Item -Path $whisperRepo -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path $binDir -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path $modelDir -Recurse -Force -ErrorAction SilentlyContinue

# Recreate target directories
New-Item -ItemType Directory -Path $binDir -Force | Out-Null
New-Item -ItemType Directory -Path $modelDir -Force | Out-Null

function DownloadFile {
    param (
        [string]$url,
        [string]$outFile
    )
    Write-Host "`nDownloading $outFile..."
    & curl.exe -L "$url" -o "$outFile"
}

# --- yt-dlp ---
$ytDlpExe = Join-Path $binDir "yt-dlp.exe"
DownloadFile "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" $ytDlpExe

# --- ffmpeg ---
$ffmpegZip = Join-Path $binDir "ffmpeg.zip"
Write-Host "Downloading ffmpeg... this may take a while."
DownloadFile "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" $ffmpegZip

Write-Host "Extracting ffmpeg..."
Expand-Archive -Path $ffmpegZip -DestinationPath $binDir -Force
Remove-Item $ffmpegZip

# Move ffmpeg.exe and ffprobe.exe to binaries
$ffmpegBin = Get-ChildItem -Path $binDir -Recurse -Filter "ffmpeg.exe" | Select-Object -First 1
if ($ffmpegBin) {
    Move-Item $ffmpegBin.FullName (Join-Path $binDir "ffmpeg.exe") -Force
}

$ffprobeBin = Get-ChildItem -Path $binDir -Recurse -Filter "ffprobe.exe" | Select-Object -First 1
if ($ffprobeBin) {
    Move-Item $ffprobeBin.FullName (Join-Path $binDir "ffprobe.exe") -Force
}

# Clean up extracted folder
$ffmpegExtractedFolder = Get-ChildItem -Path $binDir -Directory | Where-Object {
    $_.Name -like "ffmpeg-*-essentials_build"
} | Select-Object -First 1

if ($ffmpegExtractedFolder) {
    Remove-Item -Path $ffmpegExtractedFolder.FullName -Recurse -Force
    Write-Host "Removed leftover folder: $($ffmpegExtractedFolder.Name)"
}

# --- whisper-cli setup ---
try {
    Write-Host "Setting up whisper-cli from source..."
    # Call operator — actually invokes the helper script
    & "$rootDir\scripts\setup-whisper-cli.ps1"
    Set-Location $rootDir
    Write-Host "All assets downloaded. You can now Build and/or run the app!"
}
catch {
    Write-Warning "Whisper CLI setup failed. Please troubleshoot manually."
}
