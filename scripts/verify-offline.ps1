param(
    [switch]$SkipFrontendInstall
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$frontendDistDir = Join-Path $repoRoot 'frontend\dist'
$backendDistDir = Join-Path $repoRoot 'backend\web\dist'
$indexFile = Join-Path $frontendDistDir 'index.html'

if ($SkipFrontendInstall) {
    & (Join-Path $PSScriptRoot 'build-frontend.ps1') -SkipInstall
}
else {
    & (Join-Path $PSScriptRoot 'build-frontend.ps1')
}

if (-not (Test-Path $indexFile)) {
    throw "missing frontend dist index: $indexFile"
}
if (-not (Test-Path (Join-Path $backendDistDir 'index.html'))) {
    throw "missing backend hosted frontend index: $backendDistDir\index.html"
}

Write-Host '[verify-offline] checking for common external resource patterns in built dist'
$forbiddenPatterns = @(
    'http://127\.0\.0\.1:8080',
    'localhost:5173',
    'cdn\.',
    'fonts\.googleapis\.com',
    'fonts\.gstatic\.com',
    'unpkg\.com',
    'jsdelivr\.net'
)

$matches = Get-ChildItem -Path $frontendDistDir -Recurse -File |
    Select-String -Pattern $forbiddenPatterns
if ($matches) {
    $matches | ForEach-Object { Write-Host $_.Path ':' $_.LineNumber ':' $_.Line.Trim() }
    throw 'built dist still contains forbidden external dependency patterns'
}

Push-Location (Join-Path $repoRoot 'backend')
try {
    Write-Host '[verify-offline] running router tests for static hosting and SPA fallback'
    go test ./internal/router -count=1 -v
    if ($LASTEXITCODE -ne 0) {
        throw "go test ./internal/router failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}

Write-Host '[verify-offline] offline deployment checks passed'
