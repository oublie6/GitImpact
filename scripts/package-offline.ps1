param(
    [string]$OutputDir = 'artifacts\offline',
    [switch]$SkipFrontendInstall
)

$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$timestamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$bundleRoot = Join-Path $repoRoot $OutputDir
$stagingDir = Join-Path $bundleRoot "gitimpact-offline-$timestamp"
$archivePath = "$stagingDir.zip"

New-Item -ItemType Directory -Force -Path $bundleRoot | Out-Null

if ($SkipFrontendInstall) {
    & (Join-Path $PSScriptRoot 'build-release.ps1') -OutputDir (Join-Path $OutputDir "gitimpact-offline-$timestamp") -SkipFrontendInstall
}
else {
    & (Join-Path $PSScriptRoot 'build-release.ps1') -OutputDir (Join-Path $OutputDir "gitimpact-offline-$timestamp")
}

$runScript = @'
$ErrorActionPreference = "Stop"
$bundleDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$configPath = Join-Path $bundleDir "config.yaml"
if (-not (Test-Path $configPath)) {
    Copy-Item -Force (Join-Path $bundleDir "config.example.yaml") $configPath
    Write-Host "[run-offline] created config.yaml from config.example.yaml"
}
Write-Host "[run-offline] IMPORTANT: initialize database before startup."
Write-Host "[run-offline] mysql:  sql/mysql/init.sql"
Write-Host "[run-offline] dameng: sql/dameng/init.sql"
$env:GITIMPACT_CONFIG = $configPath
$binaryPath = Join-Path $bundleDir "gitimpact-backend.exe"
if (-not (Test-Path $binaryPath)) {
    $binaryPath = Join-Path $bundleDir "gitimpact-backend"
}
& $binaryPath
'@

Set-Content -Path (Join-Path $stagingDir 'run-offline.ps1') -Value $runScript -Encoding UTF8

if (Test-Path $archivePath) {
    Remove-Item -Force $archivePath
}
Compress-Archive -Path (Join-Path $stagingDir '*') -DestinationPath $archivePath

Write-Host '[package-offline] offline bundle:' $archivePath
