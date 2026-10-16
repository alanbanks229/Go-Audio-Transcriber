param (
    [string]$Model = "small.en",
    [ValidateSet("Auto", "CPU", "CUDA")]
    [string]$Backend = "Auto",
    [string]$CudaArchitectures = "native",
    [ValidateRange(1, 64)]
    [int]$BuildJobs = 4,
    [switch]$SkipUpdate
)

$ErrorActionPreference = "Stop"
$RootDir = (Resolve-Path (Join-Path $PSScriptRoot "..")).Path
$AssetsDir = Join-Path $RootDir "assets"
$WhisperDir = Join-Path $RootDir ".whisper.cpp"
$WhisperExePath = Join-Path $AssetsDir "whisper-cli.exe"

function Invoke-Checked {
    param ([string]$Command, [string[]]$Arguments)
    & $Command @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "$Command failed with exit code $LASTEXITCODE. See the output above."
    }
}

foreach ($tool in @("git", "cmake")) {
    if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
        throw "Required tool '$tool' is missing from PATH. See the Windows prerequisites in README.md."
    }
}
$Nvcc = Get-Command nvcc.exe -ErrorAction SilentlyContinue
if ($Backend -eq "Auto") {
    $Backend = if ($Nvcc) { "CUDA" } else { "CPU" }
    Write-Host "Selected Whisper backend: $Backend"
    if ($Backend -eq "CPU") {
        Write-Warning "CUDA Toolkit was not found. Building CPU-only Whisper. For NVIDIA acceleration install CUDA and Visual C++ tools, then run with -Backend CUDA."
    }
}
$ConfigureArgs = @(
    "-S", $WhisperDir, "-DBUILD_SHARED_LIBS=OFF",
    "-DWHISPER_BUILD_TESTS=OFF", "-DWHISPER_BUILD_SERVER=OFF", "-DGGML_CCACHE=OFF"
)
if ($Backend -eq "CUDA") {
    if (-not $Nvcc) {
        throw "CUDA requires NVIDIA CUDA Toolkit (nvcc.exe on PATH), a supported NVIDIA driver, and compatible Visual C++ build tools. MinGW alone cannot build this CUDA configuration."
    }
    $VsWhere = Join-Path ${env:ProgramFiles(x86)} "Microsoft Visual Studio\Installer\vswhere.exe"
    if (-not (Test-Path $VsWhere)) {
        throw "Install Visual C++ Build Tools with the C++ desktop workload and Windows SDK."
    }
    $VsJson = & $VsWhere -latest -products * -requires Microsoft.VisualStudio.Component.VC.Tools.x86.x64 -format json
    if ($LASTEXITCODE -ne 0) { throw "vswhere failed to locate the Visual C++ compiler." }
    $Vs = @($VsJson | ConvertFrom-Json) | Select-Object -First 1
    if (-not $Vs) {
        throw "No Visual C++ compiler was found. Add the C++ desktop workload and Windows SDK to Visual Studio Build Tools. A Build Tools installation without C++ is insufficient."
    }
    $VsMajor = ([version]$Vs.installationVersion).Major
    $VsYear = switch ($VsMajor) {
        16 { "2019" }
        17 { "2022" }
        18 { "2026" }
        default { throw "Unsupported Visual Studio version: $($Vs.installationVersion)" }
    }
    $CudaBin = Split-Path $Nvcc.Source -Parent
    if ((Split-Path $CudaBin -Leaf) -eq "x64") { $CudaBin = Split-Path $CudaBin -Parent }
    $CudaRoot = Split-Path $CudaBin -Parent
    $BuildDir = Join-Path $WhisperDir "build-cuda"
    $ConfigureArgs += @(
        "-B", $BuildDir, "-G", "Visual Studio $VsMajor $VsYear", "-A", "x64",
        "-T", "cuda=$CudaRoot", "-DCMAKE_GENERATOR_INSTANCE=$($Vs.installationPath)",
        "-DGGML_CUDA=ON", "-DCMAKE_CUDA_ARCHITECTURES=$CudaArchitectures"
    )
} else {
    foreach ($tool in @("gcc", "g++", "ninja")) {
        if (-not (Get-Command $tool -ErrorAction SilentlyContinue)) {
            throw "CPU build requires '$tool' on PATH. Install MinGW-w64 and Ninja, then reopen PowerShell."
        }
    }
    $BuildDir = Join-Path $WhisperDir "build-cpu"
    $ConfigureArgs += @(
        "-B", $BuildDir, "-G", "Ninja", "-DCMAKE_BUILD_TYPE=Release",
        "-DCMAKE_C_COMPILER=gcc", "-DCMAKE_CXX_COMPILER=g++", "-DGGML_CUDA=OFF"
    )
}
New-Item -ItemType Directory -Force -Path $AssetsDir | Out-Null
if (-not (Test-Path $WhisperDir)) {
    Invoke-Checked -Command git -Arguments @("clone", "https://github.com/ggml-org/whisper.cpp.git", $WhisperDir)
} elseif (-not $SkipUpdate) {
    Invoke-Checked -Command git -Arguments @("-C", $WhisperDir, "pull", "--ff-only")
}

# Download directly with curl.exe; Windows setup does not require a Bash helper.
$ModelPath = Join-Path $AssetsDir "ggml-$Model.bin"
if (-not (Test-Path $ModelPath)) {
    if (-not (Get-Command curl.exe -ErrorAction SilentlyContinue)) { throw "curl.exe is missing from PATH." }
    $ModelRepository = if ($Model -like "*tdrz*") {
        "https://huggingface.co/akashmjn/tinydiarize-whisper.cpp"
    } else {
        "https://huggingface.co/ggerganov/whisper.cpp"
    }
    $PartialModel = "$ModelPath.partial"
    Invoke-Checked -Command curl.exe -Arguments @(
        "-fL", "--retry", "5", "--retry-delay", "5", "--retry-all-errors",
        "--output", $PartialModel, "$ModelRepository/resolve/main/ggml-$Model.bin"
    )
    Move-Item -LiteralPath $PartialModel -Destination $ModelPath -Force
}
Write-Host "Building $Backend Whisper in $BuildDir..."
Invoke-Checked -Command cmake -Arguments $ConfigureArgs
Invoke-Checked -Command cmake -Arguments @("--build", $BuildDir, "--config", "Release", "--target", "whisper-cli", "--parallel", "$BuildJobs")
$BinarySubpath = if ($Backend -eq "CUDA") { "bin\Release\whisper-cli.exe" } else { "bin\whisper-cli.exe" }
$BuiltBinary = Join-Path $BuildDir $BinarySubpath
if (-not (Test-Path $BuiltBinary)) { throw "Build did not produce $BuiltBinary" }

# Bundle runtime DLLs so launching the Go app does not depend on CUDA PATH.
if ($Backend -eq "CUDA") {
    foreach ($pattern in @("cudart64*.dll", "cublas64*.dll", "cublasLt64*.dll")) {
        $RuntimeDll = @(
            Get-ChildItem -Path (Join-Path $CudaRoot "bin"), (Join-Path $CudaRoot "bin\x64") -Filter $pattern -File -ErrorAction SilentlyContinue
        ) | Select-Object -First 1
        if (-not $RuntimeDll) { throw "Required CUDA runtime $pattern was not found under $CudaRoot\bin." }
        $Destination = Join-Path $AssetsDir $RuntimeDll.Name
        if ((Test-Path $Destination) -and ((Get-FileHash $Destination).Hash -eq (Get-FileHash $RuntimeDll.FullName).Hash)) { continue }
        Copy-Item -LiteralPath $RuntimeDll.FullName -Destination $Destination -Force
    }
}
$PendingBinary = Join-Path $AssetsDir "whisper-cli.pending.exe"
Copy-Item -LiteralPath $BuiltBinary -Destination $PendingBinary -Force
Invoke-Checked -Command $PendingBinary -Arguments @("--version")

# Preserve the previous executable, including one still used by a running job.
$Backup = $null
if (Test-Path $WhisperExePath) {
    $Backup = Join-Path $AssetsDir ("whisper-cli.backup-" + [guid]::NewGuid().ToString("N") + ".exe")
    Move-Item -LiteralPath $WhisperExePath -Destination $Backup
}
try {
    Move-Item -LiteralPath $PendingBinary -Destination $WhisperExePath
} catch {
    if ($Backup -and -not (Test-Path $WhisperExePath)) { Move-Item -LiteralPath $Backup -Destination $WhisperExePath }
    throw
}
Write-Host "Installed $Backend Whisper: $WhisperExePath"
if ($Backup) { Write-Host "Previous executable preserved at: $Backup" }
Write-Host "Existing transcription jobs keep their original backend. New jobs use this build."
