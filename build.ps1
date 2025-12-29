#==============================================================================
# AutoManagement - Build Script
# Cross-platform desktop application build tool
# Vue.js 3 Frontend + Go Backend + WebView2
#==============================================================================

param(
    [switch]$SkipFrontend,
    [switch]$x64Only,
    [switch]$x86Only,
    [switch]$Verbose
)

$ErrorActionPreference = "Stop"
$ProgressPreference = "SilentlyContinue"

# Colors
function Write-Info    { param($msg) Write-Host $msg -ForegroundColor Cyan }
function Write-Success { param($msg) Write-Host $msg -ForegroundColor Green }
function Write-Warn    { param($msg) Write-Host $msg -ForegroundColor Yellow }
function Write-Err     { param($msg) Write-Host $msg -ForegroundColor Red }
function Write-Muted   { param($msg) Write-Host $msg -ForegroundColor Gray }

#==============================================================================
# Header
#==============================================================================
Write-Host ""
Write-Host "  ╔══════════════════════════════════════════════════════════════╗" -ForegroundColor Cyan
Write-Host "  ║                    AutoManagement Builder                    ║" -ForegroundColor Cyan
Write-Host "  ║              Vue.js 3 + Go + WebView2 Desktop App            ║" -ForegroundColor Cyan
Write-Host "  ╚══════════════════════════════════════════════════════════════╝" -ForegroundColor Cyan
Write-Host ""

#==============================================================================
# Prerequisites Check
#==============================================================================
Write-Info "[1/4] Checking prerequisites..."

# Check Go
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Err "ERROR: Go is not installed or not in PATH"
    Write-Muted "       Install from: https://go.dev/dl/"
    exit 1
}
$goVersion = (go version) -replace "go version ", ""
Write-Success "      Go: $goVersion"

# Check Node.js or Bun
$hasBun = Get-Command bun -ErrorAction SilentlyContinue
$hasNpm = Get-Command npm -ErrorAction SilentlyContinue
if ($hasBun) {
    $bunVersion = (bun --version)
    Write-Success "      Bun: v$bunVersion"
    $packageManager = "bun"
} elseif ($hasNpm) {
    $npmVersion = (npm --version)
    Write-Success "      npm: v$npmVersion"
    $packageManager = "npm"
} else {
    Write-Err "ERROR: Neither Bun nor npm is installed"
    Write-Muted "       Install Node.js from: https://nodejs.org/"
    exit 1
}

# Check goversioninfo
$hasGoversioninfo = $false
if (Get-Command goversioninfo -ErrorAction SilentlyContinue) {
    $hasGoversioninfo = $true
    Write-Success "      goversioninfo: installed"
} else {
    Write-Warn "      goversioninfo: not found"
    Write-Host ""
    Write-Warn "goversioninfo is required for EXE metadata (icon, version info)."
    Write-Host ""
    Write-Info "The following command will be executed to install it:"
    Write-Host ""
    Write-Host "  go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest" -ForegroundColor White
    Write-Host ""
    
    $confirm = Read-Host "Do you want to proceed with installation? (Y/N)"
    if ($confirm -eq "Y" -or $confirm -eq "y") {
        Write-Info "Installing goversioninfo..."
        go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
        if ($LASTEXITCODE -eq 0) {
            $hasGoversioninfo = $true
            Write-Success "goversioninfo installed successfully"
        } else {
            Write-Err "Failed to install goversioninfo"
        }
    } else {
        Write-Muted "Installation skipped."
        
        if (Test-Path "resource.syso") {
            Write-Warn "Using existing resource.syso file."
        } else {
            Write-Err "ERROR: resource.syso not found. Cannot build without EXE metadata."
            Write-Muted "       Run: go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest"
            exit 1
        }
    }
}

Write-Host ""

#==============================================================================
# Read Version Info
#==============================================================================
Set-Location -Path $PSScriptRoot

$verObj = Get-Content versioninfo.json | ConvertFrom-Json
$version = $verObj.StringFileInfo.ProductVersion
if (-not $version) { $version = $verObj.StringFileInfo.FileVersion }

$appName = $verObj.StringFileInfo.InternalName
if (-not $appName) { $appName = "AutoManagement" }

Write-Info "[2/4] Project Information"
Write-Host "      App Name: $appName" -ForegroundColor White
Write-Host "      Version:  v$version" -ForegroundColor White
Write-Host ""

#==============================================================================
# Frontend Build
#==============================================================================
if (-not $SkipFrontend) {
    Write-Info "[3/4] Building frontend..."
    Set-Location -Path "$PSScriptRoot\frontend"
    
    if ($packageManager -eq "bun") {
        bun run build
    } else {
        npm run build
    }
    
    if ($LASTEXITCODE -ne 0) {
        Write-Err "ERROR: Frontend build failed"
        exit 1
    }
    Write-Success "      Frontend build completed"
    Set-Location -Path $PSScriptRoot
} else {
    Write-Muted "[3/4] Frontend build skipped (--SkipFrontend)"
}

Write-Host ""

#==============================================================================
# Go Build
#==============================================================================
Write-Info "[4/4] Building executables..."

# Determine architectures
$arches = @()
if ($x64Only) {
    $arches = @("amd64")
} elseif ($x86Only) {
    $arches = @("386")
} else {
    $arches = @("amd64", "386")
}

$buildSuccess = @()
$origGOARCH = (& go env GOARCH).Trim()
$origCGO = (& go env CGO_ENABLED).Trim()

foreach ($arch in $arches) {
    switch ($arch) {
        "amd64" { $archLabel = "x64" }
        "386"   { $archLabel = "x86" }
        default { $archLabel = $arch }
    }
    
    $exeName = "$appName-$archLabel-v$version.exe"
    Write-Host ""
    Write-Info "      Building $archLabel..."
    
    # Create temporary versioninfo with correct OriginalFilename
    $verObj.StringFileInfo.OriginalFilename = $exeName
    $tempVerFile = Join-Path $PSScriptRoot ("versioninfo.$arch.temp.json")
    $verObj | ConvertTo-Json -Depth 8 | Out-File -FilePath $tempVerFile -Encoding ascii
    
    # Generate resource.syso
    if ($hasGoversioninfo) {
        Write-Muted "      Generating resource.syso..."
        
        if (Test-Path "resource.syso") {
            Remove-Item "resource.syso" -Force
        }
        
        goversioninfo -icon="assets/App.ico" -o="resource.syso" $tempVerFile
        if ($LASTEXITCODE -ne 0) {
            Write-Err "      Failed to generate resource.syso for $archLabel"
            Remove-Item $tempVerFile -ErrorAction SilentlyContinue
            continue
        }
    }
    
    # Set environment for cross-compilation
    $env:GOARCH = $arch
    $env:CGO_ENABLED = "0"
    
    # Build
    Write-Muted "      Compiling Go application..."
    go build -ldflags="-H windowsgui -s -w" -o $exeName .
    
    if ($LASTEXITCODE -ne 0) {
        Write-Err "      Build failed for $archLabel"
    } else {
        $exeSize = (Get-Item $exeName).Length / 1MB
        Write-Success "      Build successful: $exeName"
        Write-Muted ("      Size: " + [math]::Round($exeSize, 2).ToString() + " MB")
        $buildSuccess += $exeName
    }
    
    Remove-Item $tempVerFile -ErrorAction SilentlyContinue
}

# Restore environment
if ($origGOARCH) { $env:GOARCH = $origGOARCH } else { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue }
if ($origCGO) { $env:CGO_ENABLED = $origCGO } else { Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue }

#==============================================================================
# Summary
#==============================================================================
Write-Host ""
Write-Host "  ══════════════════════════════════════════════════════════════" -ForegroundColor Cyan

if ($buildSuccess.Count -eq 0) {
    Write-Err "  Build failed. No executables were created."
    exit 1
} else {
    Write-Success "  Build completed successfully!"
    Write-Host ""
    Write-Info "  Output files:"
    foreach ($exe in $buildSuccess) {
        $size = [math]::Round((Get-Item $exe).Length / 1MB, 2)
        Write-Host "    - $exe ($size MB)" -ForegroundColor White
    }
}

Write-Host ""
Write-Host "  ══════════════════════════════════════════════════════════════" -ForegroundColor Cyan
Write-Host ""
