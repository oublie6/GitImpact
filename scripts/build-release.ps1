param(
    [string]$OutputDir = 'artifacts\release',
    [switch]$SkipFrontendInstall
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$outputPath = Join-Path $repoRoot $OutputDir
$binaryName = if ($env:OS -eq 'Windows_NT') { 'gitimpact-backend.exe' } else { 'gitimpact-backend' }

Write-Host '[build-release] building frontend and syncing dist'
if ($SkipFrontendInstall) {
    & (Join-Path $PSScriptRoot 'build-frontend.ps1') -SkipInstall
}
else {
    & (Join-Path $PSScriptRoot 'build-frontend.ps1')
}

if (Test-Path $outputPath) {
    Remove-Item -Recurse -Force $outputPath
}
New-Item -ItemType Directory -Force -Path $outputPath | Out-Null

Push-Location (Join-Path $repoRoot 'backend')
try {
    Write-Host '[build-release] building backend binary'
    $env:GOFLAGS = '-mod=vendor'
    go build -trimpath -ldflags "-s -w" -o (Join-Path $outputPath $binaryName) ./cmd/server
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}

Copy-Item -Force (Join-Path $repoRoot 'backend\config.example.yaml') (Join-Path $outputPath 'config.example.yaml')
Copy-Item -Recurse -Force (Join-Path $repoRoot 'backend\web') (Join-Path $outputPath 'web')

Write-Host '[build-release] release output ready at' $outputPath
