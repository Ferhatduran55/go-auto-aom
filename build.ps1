# Oto Parca Siparis Sistemi - Build Script
# Vue.js frontend build + Go EXE build

Write-Host "Oto Parca Siparis Sistemi - Build" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan

# Frontend build
Write-Host "`nFrontend build ediliyor..." -ForegroundColor Yellow
Set-Location -Path "$PSScriptRoot\frontend"
bun run build

if ($LASTEXITCODE -ne 0) {
    Write-Host "Frontend build başarısız!" -ForegroundColor Red
    exit 1
}
Write-Host "Frontend build tamamlandi" -ForegroundColor Green

# Go build
Write-Host "`nGo EXE build ediliyor..." -ForegroundColor Yellow
Set-Location -Path $PSScriptRoot

# Read version from versioninfo.json
$verObj = Get-Content versioninfo.json | ConvertFrom-Json
$version = $verObj.StringFileInfo.ProductVersion
if (-not $version) { $version = $verObj.StringFileInfo.FileVersion }

# Architectures to build: amd64 (x64) and 386 (x86)
$arches = @("amd64","386")
$buildSuccess = @()

$origGOARCH = (& go env GOARCH).Trim()
$origCGO = (& go env CGO_ENABLED).Trim()

foreach ($arch in $arches) {
    switch ($arch) {
        "amd64" { $archLabel = "x64" }
        "386"   { $archLabel = "x86" }
        default  { $archLabel = $arch }
    }

    $exeName = "OtoParcaSiparis-$archLabel-v$version.exe"

    # Create a temporary versioninfo file with OriginalFilename set to the final exe name
    $verObj.StringFileInfo.OriginalFilename = $exeName
    $tempVerFile = Join-Path $PSScriptRoot ("versioninfo." + $arch + ".temp.json")
    $verObj | ConvertTo-Json -Depth 8 | Set-Content -Path $tempVerFile -Encoding UTF8

    # Generate resource.syso using goversioninfo if available
    if (Get-Command goversioninfo -ErrorAction SilentlyContinue) {
        Write-Host ("goversioninfo found, generating resource.syso for " + $archLabel) -ForegroundColor Yellow
        goversioninfo -icon="assets/App.ico" -o="resource.syso" $tempVerFile
        if ($LASTEXITCODE -ne 0) {
            Write-Host ("resource.syso could not be generated for " + $archLabel) -ForegroundColor Red
        } else {
            Write-Host "resource.syso generated" -ForegroundColor Green
        }
    } else {
        Write-Host "goversioninfo not found; using existing resource.syso (if any)." -ForegroundColor Yellow
    }

    # Set env for cross-build
    $env:GOARCH = $arch
    $env:CGO_ENABLED = "0"

    go build -ldflags="-H windowsgui -s -w" -o $exeName .
    if ($LASTEXITCODE -ne 0) {
        Write-Host ("Go build failed for " + $archLabel) -ForegroundColor Red
    } else {
        Write-Host ("Build succeeded: " + $exeName) -ForegroundColor Green
        $exeSize = (Get-Item $exeName).Length / 1MB
        Write-Host ('Dosya: ' + $exeName) -ForegroundColor White
        Write-Host ('Boyut: ' + [math]::Round($exeSize, 2).ToString() + ' MB') -ForegroundColor White
        $buildSuccess += $exeName
    }

    Remove-Item $tempVerFile -ErrorAction SilentlyContinue
}

# Restore environment
if ($origGOARCH) { $env:GOARCH = $origGOARCH } else { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue }
if ($origCGO) { $env:CGO_ENABLED = $origCGO } else { Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue }

if ($buildSuccess.Count -eq 0) {
    Write-Host "No build succeeded." -ForegroundColor Red
    exit 1
} else {
    Write-Host "`nBuild complete:" -ForegroundColor Cyan
    foreach ($b in $buildSuccess) { Write-Host (" - " + $b) }
}
