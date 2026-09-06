$ErrorActionPreference = "Stop"
$BinaryPath = Join-Path $HOME ".local\bin\lapip.exe"

if (Test-Path $BinaryPath) {
    Remove-Item -Force $BinaryPath
    Write-Host "removed lapip.exe from $BinaryPath"
} else {
    Write-Host "lapip.exe was not found at $BinaryPath"
}
