$ErrorActionPreference = "Stop"

$RepoUrl = "https://github.com/ESHAYAT102/lapip.git"
$BinaryName = "lapip.exe"
$InstallDir = Join-Path $HOME ".local\bin"
$CloneDir = $null

if (-not (Get-Command go -ErrorAction SilentlyContinue)) { throw "Go is required to install lapip" }
if (-not (Get-Command git -ErrorAction SilentlyContinue)) { throw "git is required to install lapip" }

try {
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null
    $CloneDir = Join-Path ([System.IO.Path]::GetTempPath()) "lapip-$([System.Guid]::NewGuid())"
    Write-Host "cloning $RepoUrl"
    git clone --depth 1 $RepoUrl $CloneDir
    Write-Host "building lapip"
    Push-Location $CloneDir
    try { go build -o (Join-Path $InstallDir $BinaryName) . }
    finally { Pop-Location }
    Write-Host "installed $BinaryName to $(Join-Path $InstallDir $BinaryName)"
}
finally {
    if ($CloneDir -and (Test-Path $CloneDir)) { Remove-Item -Recurse -Force $CloneDir }
}
