$RootDir = Resolve-Path "."
$AppName = "Go-Audio-Transcriber"
$BuildDir = "$RootDir\build"
$MainFile = "main.go"

Write-Host "Cleaning old build..."
Remove-Item -Path $BuildDir -Recurse -Force -ErrorAction SilentlyContinue

# Recreate build directory
New-Item -ItemType Directory -Force -Path $BuildDir | Out-Null

Write-Host "Building for Windows..."

if (Get-Command "fyne.exe" -ErrorAction SilentlyContinue) {
    Write-Host "Using fyne package..."
    fyne package -os windows -icon "icon.png" -name $AppName --appVersion "0.9.0" --appBuild "0"
    Move-Item "$AppName.exe" "$BuildDir/$AppName.exe" -Force
} else {
    Write-Host "Fyne CLI not found. Falling back to go build..."
    go build -o "$BuildDir/$AppName.exe" $MainFile
}

Write-Host "Build complete: $BuildDir"
