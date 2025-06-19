$AppName = "Go-Audio-Transcriber"
$OutDir = "./dist"
$MainFile = "main.go"

New-Item -ItemType Directory -Force -Path $OutDir | Out-Null

Write-Host "Building for Windows..."

if (Get-Command "fyne.exe" -ErrorAction SilentlyContinue) {
    Write-Host "Using fyne package..."
    fyne package -os windows -name $AppName -icon "icon.png" --appVersion "0.9.0" --appBuild "0"
    Move-Item "$AppName.exe" "$OutDir/$AppName.exe" -Force
} else {
    Write-Host "Fyne CLI not found. Falling back to go build..."
    go build -o "$OutDir/$AppName.exe" $MainFile
}

Write-Host "Build complete: $OutDir"
