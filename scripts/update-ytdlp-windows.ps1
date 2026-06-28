# PowerShell updater for the bundled Windows yt-dlp binary.
$ErrorActionPreference = "Stop"

$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$AssetsDir = Join-Path $ProjectRoot "assets"
$YtDlpPath = Join-Path $AssetsDir "yt-dlp.exe"
$YtDlpUrl = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.exe"

New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null

Write-Host "Downloading latest Windows yt-dlp..."
& curl.exe -fL $YtDlpUrl -o $YtDlpPath

if ($LASTEXITCODE -ne 0) {
    throw "Failed to download yt-dlp.exe from $YtDlpUrl"
}

$Version = & $YtDlpPath --version
if ($LASTEXITCODE -ne 0) {
    throw "Downloaded yt-dlp.exe, but version check failed."
}

Write-Host "yt-dlp updated: $Version"
