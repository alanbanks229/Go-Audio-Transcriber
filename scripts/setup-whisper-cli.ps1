param (
    [string]$Model = "small.en"
)

$ErrorActionPreference = "Stop"

$RootDir = Resolve-Path "$PSScriptRoot\.."
$AssetsDir = "$RootDir\assets"
$WhisperDir = "$RootDir\.whisper.cpp"
$WhisperRepo = "https://github.com/ggerganov/whisper.cpp.git"
$WhisperExePath = "$AssetsDir\whisper-cli.exe"

# Ensure assets dir exists
New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null

# Clone or pull whisper.cpp
Write-Host "Cloning or updating whisper.cpp..."
if (-Not (Test-Path $WhisperDir)) {
    git clone $WhisperRepo $WhisperDir
} else {
    Push-Location $WhisperDir
    git pull
    Pop-Location
}

# Run model download script (Step 2)
Write-Host "⬇Downloading model ($Model) if needed..."
$downloadScript = Join-Path $WhisperDir "models\download-ggml-model.sh"
if (-Not (Test-Path "$AssetsDir\ggml-$Model.bin")) {
    # Run bash from PowerShell
    bash $downloadScript $Model
    # Move model from whisper.cpp/models to assets
    $sourceModel = Join-Path $WhisperDir "models\ggml-$Model.bin"
    Copy-Item $sourceModel -Destination "$AssetsDir\ggml-$Model.bin" -Force
}

# Build whisper-cli with CMake (Step 3)
Write-Host "Building whisper.cpp with CMake..."
Push-Location $WhisperDir

if (-Not (Test-Path "build")) {
    New-Item -ItemType Directory -Path "build" | Out-Null
}
Set-Location "build"

cmake .. | Out-Null
cmake --build . --config Release | Out-Null

Pop-Location

# Copy built binary to assets (Step 4)
Write-Host "Copying whisper-cli.exe to assets..."
$builtBinary = Get-ChildItem "$WhisperDir\build" -Recurse -Filter "whisper-cli.exe" | Select-Object -First 1
if ($builtBinary) {
    Copy-Item $builtBinary.FullName $WhisperExePath -Force
    Write-Host "whisper-cli.exe is ready: $WhisperExePath"
} else {
    Write-Warning "Build completed, but whisper.exe was not found."
}
