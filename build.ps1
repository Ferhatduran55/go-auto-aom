# Build Script
param(
    [switch]$SkipFrontend,
    [switch]$x64Only,
    [switch]$x86Only
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
Write-Host "AutoManagement Builder - Single EXE" -ForegroundColor Cyan
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

# Check for Bun or Node.js (prefer Bun if available)
$hasBun = Get-Command bun -ErrorAction SilentlyContinue
$hasNpm = Get-Command npm -ErrorAction SilentlyContinue
if (-not $hasBun -and -not $hasNpm) {
    Write-Err "ERROR: Neither Bun nor Node.js (npm) is installed or in PATH"
    Write-Muted "       Install Bun: https://bun.sh or Node.js: https://nodejs.org"
    exit 1
}
if ($hasBun) {
    $bunVersion = bun --version
    Write-Success "      bun: v$bunVersion"
} else {
    $npmVersion = npm --version
    Write-Success "      npm: v$npmVersion"
}

# Check goversioninfo
$hasGoversioninfo = Get-Command goversioninfo -ErrorAction SilentlyContinue
if ($hasGoversioninfo) {
    Write-Success "      goversioninfo: installed"
} else {
    Write-Warn "      goversioninfo: not found"
    Write-Muted "      Installing goversioninfo..."
    go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
    if ($LASTEXITCODE -eq 0) {
        $hasGoversioninfo = $true
        Write-Success "      goversioninfo: installed"
    } else {
        Write-Err "      Failed to install goversioninfo"
        exit 1
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
    
    if ($hasBun) { bun run build } else { npm run build }
    
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
Write-Info "[4/4] Building Go application..."

# Determine architectures
$archs = @()

# Mutually exclusive flags check
if ($x64Only -and $x86Only) {
    Write-Err "ERROR: Cannot specify both -x64Only and -x86Only"
    exit 1
}

if ($x64Only) {
    $archs = @("amd64")
} elseif ($x86Only) {
    $archs = @("386")
} else {
    $archs = @("amd64", "386")  # Default: build both x64 and x86
}

$buildSuccess = @()
$origGOARCH = $env:GOARCH
$origCGO = $env:CGO_ENABLED

foreach ($arch in $archs) {
    $archLabel = if ($arch -eq "amd64") { "x64" } else { "x86" }
    $exeName = "$appName-$archLabel-v$version.exe"
    
    Write-Muted "      Building $archLabel..."
    
    # Generate version info
    $verObj.StringFileInfo.OriginalFilename = $exeName
    $tempVerFile = Join-Path $PSScriptRoot ("versioninfo.$arch.temp.json")
    $verObj | ConvertTo-Json -Depth 8 | Out-File -FilePath $tempVerFile -Encoding ascii
    
    if ($hasGoversioninfo) {
        Write-Muted "      Generating resource.syso..."
        if (Test-Path "resource.syso") { Remove-Item "resource.syso" -Force }
        $goviArgs = @("-icon=assets/App.ico", "-o=resource.syso")
        if ($arch -eq "amd64") { $goviArgs += "-64" }
        goversioninfo @goviArgs $tempVerFile
        if ($LASTEXITCODE -ne 0) {
            Write-Err "      Failed to generate resource.syso"
            Remove-Item $tempVerFile -ErrorAction SilentlyContinue
            continue
        }
    }
    
    # Environment
    $env:GOARCH = $arch
    $env:CGO_ENABLED = "0"  # Pure Go - no CGO needed
    
    # Build
    Write-Muted "      Compiling..."
    # Determine Bleve/WebView versions (env override -> parse go.mod -> fallback)
    if ($env:BLEVE_VERSION) {
        $bleveVersion = $env:BLEVE_VERSION
    } else {
        try {
            # Use 'go list -m' to reliably fetch the module version
            $bleveVersion = (& go list -m -f '{{.Version}}' github.com/blevesearch/bleve/v2).Trim()
            if (-not $bleveVersion) {
                Write-Err "ERROR: Could not determine Bleve version via 'go list'. Set BLEVE_VERSION in CI environment."
                exit 1
            }
            if (-not $bleveVersion.StartsWith('v')) { $bleveVersion = 'v' + $bleveVersion }
        } catch {
            Write-Err "ERROR: Failed to determine Bleve version: $_"
            exit 1
        }
    }

    if ($env:WEBVIEW_VERSION) {
        $webviewVersion = $env:WEBVIEW_VERSION
    } else {
        try {
            # Use 'go list -m' to get the go-webview2 module version
            $webviewVersion = (& go list -m -f '{{.Version}}' github.com/jchv/go-webview2).Trim()
            if (-not $webviewVersion) {
                Write-Err "ERROR: Could not determine WebView module version via 'go list'. Set WEBVIEW_VERSION in CI environment."
                exit 1
            }
        } catch {
            Write-Err "ERROR: Failed to determine WebView module version: $_"
            exit 1
        }
    }


    $ld = "-H windowsgui -s -w -X 'main.AppVersion=$version' -X main.BleveVersion=$bleveVersion -X 'main.WebViewVersion=$webviewVersion'"
    & go build -ldflags $ld -o $exeName .
    
    if ($LASTEXITCODE -ne 0) {
        Write-Err "      Build failed for $archLabel"
    } else {
        Write-Success "      Build successful for $archLabel"
        $buildSuccess += $exeName
    }
    
    Remove-Item $tempVerFile -ErrorAction SilentlyContinue
}

# Restore environment
if ($origGOARCH) { $env:GOARCH = $origGOARCH } else { Remove-Item Env:GOARCH -ErrorAction SilentlyContinue }
if ($origCGO) { $env:CGO_ENABLED = $origCGO } else { Remove-Item Env:CGO_ENABLED -ErrorAction SilentlyContinue }

Write-Host ""

#==============================================================================
# Summary
#==============================================================================
Write-Host ""

if ($buildSuccess.Count -eq 0) {
    Write-Err "  Build failed. No executables were created."
    exit 1
}

Write-Host "BUILD SUCCESSFUL" -ForegroundColor Green
Write-Host ""

foreach ($exe in $buildSuccess) {
    $size = [math]::Round((Get-Item $exe).Length / 1MB, 2)
    Write-Host "  📦 $exe ($size MB) (All-in-One)" -ForegroundColor White
    Write-Host ""
}