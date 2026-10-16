param (
    [ValidateSet("Auto", "CPU", "CUDA")]
    [string]$Backend = "Auto",
    [string]$CudaArchitectures = "native"
)

# PowerShell Setup Script for Windows
$ErrorActionPreference = "Stop"

$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
$AssetsDir = Join-Path $ProjectRoot "assets"
New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null

function DownloadFile {
    param (
        [string]$url,
        [string]$outFile
    )
    Write-Host "`nDownloading $outFile..."
    & curl.exe -fL "$url" -o "$outFile"
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to download $url"
    }
}

# --- yt-dlp ---
& "$PSScriptRoot/update-ytdlp-windows.ps1"

# --- ffmpeg ---
$ffmpegZip = "$AssetsDir/ffmpeg.zip"
Write-Host "Next Download takes a while..."
DownloadFile "https://www.gyan.dev/ffmpeg/builds/ffmpeg-release-essentials.zip" $ffmpegZip

Write-Host "Extracting ffmpeg..."
$ffmpegExtractDir = Join-Path $AssetsDir ("ffmpeg-extract-" + [guid]::NewGuid().ToString("N"))
New-Item -ItemType Directory -Path $ffmpegExtractDir | Out-Null
try {
    Expand-Archive -LiteralPath $ffmpegZip -DestinationPath $ffmpegExtractDir
    # Search only this archive, never existing executables from a previous run.
    foreach ($name in @("ffmpeg.exe", "ffprobe.exe")) {
        $binary = Get-ChildItem -LiteralPath $ffmpegExtractDir -Recurse -Filter $name -File | Select-Object -First 1
        if (-not $binary) { throw "$name was not found in the downloaded archive." }
        Move-Item -LiteralPath $binary.FullName -Destination (Join-Path $AssetsDir $name) -Force
    }
    Remove-Item -LiteralPath $ffmpegZip
} finally {
    $resolvedExtractDir = (Resolve-Path -LiteralPath $ffmpegExtractDir).Path
    $expectedPrefix = (Resolve-Path -LiteralPath $AssetsDir).Path + "\ffmpeg-extract-"
    if (-not $resolvedExtractDir.StartsWith($expectedPrefix, [StringComparison]::OrdinalIgnoreCase)) {
        throw "Refusing to clean unexpected extraction path: $resolvedExtractDir"
    }
    Remove-Item -LiteralPath $resolvedExtractDir -Recurse -Force
}

# --- whisper-cli setup ---
try {
    Write-Host "Setting up whisper-cli from source..."
    & (Join-Path $PSScriptRoot "setup-whisper-cli.ps1") -Backend $Backend -CudaArchitectures $CudaArchitectures
}
catch {
    throw "Whisper setup failed: $($_.Exception.Message)"
}

Write-Host "Assets are ready. Build and run the app using the README instructions."
