param(
    [switch]$SkipInstall
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$frontendDir = Join-Path $repoRoot 'frontend'
$frontendDistDir = Join-Path $frontendDir 'dist'
$backendWebDir = Join-Path $repoRoot 'backend\web\dist'

Write-Host '[build-frontend] repo root:' $repoRoot
Push-Location $frontendDir
try {
    if (-not $SkipInstall) {
        Write-Host '[build-frontend] installing frontend dependencies with npm ci'
        npm ci
        if ($LASTEXITCODE -ne 0) {
            throw "npm ci failed with exit code $LASTEXITCODE"
        }
    }

    Write-Host '[build-frontend] building frontend dist'
    npm run build:offline
    if ($LASTEXITCODE -ne 0) {
        throw "npm run build:offline failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}

if (-not (Test-Path $frontendDistDir)) {
    throw "frontend dist directory not found: $frontendDistDir"
}

if (Test-Path $backendWebDir) {
    Remove-Item -Recurse -Force $backendWebDir
}
New-Item -ItemType Directory -Force -Path $backendWebDir | Out-Null
Copy-Item -Recurse -Force (Join-Path $frontendDistDir '*') $backendWebDir

Write-Host '[build-frontend] synced frontend dist to' $backendWebDir
