param (
    [string]$Model = "small.en"
)

$ErrorActionPreference = "Stop"

$rootDir = Resolve-Path "."
$assetsDir = "$rootDir\assets"
$clonedDir = "$rootDir\.whisper.cpp"
$whisperURL = "https://github.com/ggerganov/whisper.cpp.git"
$updatedWhisperExecPath = "$assetsDir\binaries\whisper-cli.exe"

# Clone or pull whisper.cpp
Write-Host "Cloning or updating whisper.cpp..."
if (-Not (Test-Path $clonedDir)) {
    git clone $whisperURL $clonedDir
    Set-Location $clonedDir
} else {
    Set-Location $clonedDir
    git pull
}

# Build whisper-cli with CMake (Step 2)
Write-Host "Building whisper.cpp with CMake..."
cmake -B build | Out-Null
cmake --build build --config Release | Out-Null

# Copy built binary to assets (Step 3)
Write-Host "Copying whisper-cli.exe to assets..."
$builtBinary = Get-ChildItem "$clonedDir\build" -Recurse -Filter "whisper-cli.exe" | Select-Object -First 1
if ($builtBinary) {
    Copy-Item "$clonedDir\build\bin\whisper-cli.exe" $updatedWhisperExecPath -Force
    Write-Host "whisper-cli.exe is ready: $updatedWhisperExecPath"
} else {
    Write-Warning "Build completed, but whisper.exe was not found."
}

# Run model download script (Step 4)
if (-Not (Test-Path "$assetsDir\models\ggml-$Model.bin")) {
    $downloadScript = Join-Path $clonedDir "models\download-ggml-model.cmd"
    # Directly invoke the .cmd script
    & $downloadScript $Model

    # We are currently located in the cloned Dir.
    # Move model from clonedDir/models to assets
    $sourceModel = Join-Path $clonedDir "ggml-$Model.bin"
    Copy-Item $sourceModel -Destination "$assetsDir\models\ggml-$Model.bin" -Force
}