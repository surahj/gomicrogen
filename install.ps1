# gomicrogen Installation Script for Windows
# This script downloads and installs gomicrogen CLI tool on Windows

param(
    [string]$Version = "",
    [string]$InstallDir = "$env:LOCALAPPDATA\gomicrogen",
    [switch]$Help
)

# Configuration
$Repo = "surahj/gomicrogen"
$BinaryName = "gomicrogen.exe"
$TempDir = "$env:TEMP\gomicrogen-install"

# Function to write colored output
function Write-Status {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Blue
}

function Write-Success {
    param([string]$Message)
    Write-Host "[SUCCESS] $Message" -ForegroundColor Green
}

function Write-Warning {
    param([string]$Message)
    Write-Host "[WARNING] $Message" -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

# Function to detect architecture
function Get-Architecture {
    if ([Environment]::Is64BitOperatingSystem) {
        return "amd64"
    } else {
        return "386"
    }
}

# Function to get latest version
function Get-LatestVersion {
    Write-Status "Fetching latest version..."
    
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Method Get
        $latestVersion = $response.tag_name
        Write-Status "Latest version: $latestVersion"
        return $latestVersion
    }
    catch {
        Write-Warning "Could not fetch latest version, using 'latest'"
        return "latest"
    }
}

# Function to download binary
function Download-Binary {
    param([string]$Version)
    
    $arch = Get-Architecture
    $downloadUrl = "https://github.com/$Repo/releases/download/$Version/${BinaryName.Replace('.exe', '')}-windows-$arch.exe"
    $localFilename = "${BinaryName.Replace('.exe', '')}-windows-$arch.exe"
    
    Write-Status "Downloading from: $downloadUrl"
    
    # Create temp directory
    if (!(Test-Path $TempDir)) {
        New-Item -ItemType Directory -Path $TempDir -Force | Out-Null
    }
    
    $downloadPath = Join-Path $TempDir $localFilename
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath
        Write-Success "Download completed"
        return $downloadPath
    }
    catch {
        Write-Error "Failed to download binary: $($_.Exception.Message)"
        exit 1
    }
}

# Function to install binary
function Install-Binary {
    param([string]$BinaryPath)
    
    Write-Status "Installing to $InstallDir..."
    
    # Create install directory if it doesn't exist
    if (!(Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    
    $installPath = Join-Path $InstallDir $BinaryName
    
    try {
        Copy-Item -Path $BinaryPath -Destination $installPath -Force
        Write-Success "Installation completed"
    }
    catch {
        Write-Error "Failed to install binary: $($_.Exception.Message)"
        exit 1
    }
}

# Function to add to PATH
function Add-ToPath {
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    
    if ($currentPath -notlike "*$InstallDir*") {
        Write-Status "Adding $InstallDir to PATH..."
        $newPath = "$currentPath;$InstallDir"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Success "Added to PATH"
    } else {
        Write-Status "Already in PATH"
    }
}

# Function to verify installation
function Test-Installation {
    Write-Status "Verifying installation..."
    
    $installPath = Join-Path $InstallDir $BinaryName
    
    if (Test-Path $installPath) {
        Write-Success "$BinaryName is now available"
        try {
            $version = & $installPath version 2>$null
            if ($version) {
                Write-Status "Version: $version"
            }
        }
        catch {
            Write-Status "Version info not available"
        }
    } else {
        Write-Error "Installation verification failed. $BinaryName is not found"
        exit 1
    }
}

# Function to cleanup
function Remove-TempFiles {
    Write-Status "Cleaning up temporary files..."
    if (Test-Path $TempDir) {
        Remove-Item -Path $TempDir -Recurse -Force
    }
}

# Function to show usage
function Show-Usage {
    Write-Host "Usage: .\install.ps1 [OPTIONS]"
    Write-Host ""
    Write-Host "Options:"
    Write-Host "  -Version VERSION     Install specific version (default: latest)"
    Write-Host "  -InstallDir DIR      Install to specific directory (default: $env:LOCALAPPDATA\gomicrogen)"
    Write-Host "  -Help               Show this help message"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\install.ps1                      # Install latest version"
    Write-Host "  .\install.ps1 -Version v1.0.0      # Install specific version"
    Write-Host "  .\install.ps1 -InstallDir C:\tools # Install to custom directory"
}

# Main execution
if ($Help) {
    Show-Usage
    exit 0
}

try {
    Write-Status "Starting gomicrogen installation..."
    
    # Get version if not specified
    if ([string]::IsNullOrEmpty($Version)) {
        $Version = Get-LatestVersion
    }
    
    # Download binary
    $binaryPath = Download-Binary -Version $Version
    
    # Install binary
    Install-Binary -BinaryPath $binaryPath
    
    # Add to PATH
    Add-ToPath
    
    # Verify installation
    Test-Installation
    
    # Cleanup
    Remove-TempFiles
    
    Write-Success "Installation completed successfully!"
    Write-Host ""
    Write-Host "You can now use gomicrogen:"
    Write-Host "  gomicrogen --help"
    Write-Host ""
    Write-Host "Example usage:"
    Write-Host "  gomicrogen new my-service --module github.com/your-org/my-service"
    Write-Host ""
    Write-Host "Note: You may need to restart your terminal or run 'refreshenv' to use gomicrogen"
}
catch {
    Write-Error "Installation failed: $($_.Exception.Message)"
    exit 1
}
finally {
    Remove-TempFiles
} 