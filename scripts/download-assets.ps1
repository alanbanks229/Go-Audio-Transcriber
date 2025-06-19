# PowerShell Setup Script for Windows
$AssetsDir = "./assets"
New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null

function DownloadFile {
    param (
        [string]$url,
        [string]$outFile
    )
    Write-Host "`nDownloading $outFile..."
    & curl.exe -L "$url" -o "$outFile"
}

# --- yt-dlp ---
DownloadFile "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe" "$AssetsDir/yt-dlp.exe"

# --- ffmpeg ---
$ffmpegZip = "$AssetsDir/ffmpeg.zip"
Write-Host "Next Download takes a while..."
DownloadFile "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" $ffmpegZip

Write-Host "Extracting ffmpeg..."
Expand-Archive -Path $ffmpegZip -DestinationPath $AssetsDir -Force
Remove-Item $ffmpegZip

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

# --- Clean up extracted ffmpeg folder ---
$ffmpegExtractedFolder = Get-ChildItem -Path $AssetsDir -Directory | Where-Object {
    $_.Name -like "ffmpeg-*-essentials_build"
} | Select-Object -First 1

if ($ffmpegExtractedFolder) {
    Remove-Item -Path $ffmpegExtractedFolder.FullName -Recurse -Force
    Write-Host "🧹 Removed leftover folder: $($ffmpegExtractedFolder.Name)"
}

# --- whisper-cli setup ---
try {
    Write-Host "Setting up whisper-cli from source..."
    & "./scripts/setup-whisper-cli.ps1"
}
catch {
    Write-Warning "Whisper binaries failed to download or build. Please try again later or check manually."
}

Write-Host "If no errors occurred you can now Run/Build the APP! See README."