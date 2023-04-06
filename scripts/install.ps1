#Requires -RunAsAdministrator

Param([Parameter(Mandatory=$true)][string]$version, [Parameter(Mandatory=$true)][string]$arch)
$ErrorActionPreference = "Stop"

$appName = "Hoist"
$exeName = "hoist.exe"
$appPath = "${env:ProgramFiles}\${appName}"
$path = ([System.Environment]::GetEnvironmentVariable("PATH", "User"))

$path = $path.TrimEnd(";")
$splitPath = $path.Split(";")

Write-Output "Creating application directory..."
New-Item -Path $env:ProgramFiles -Name $appName -ItemType "directory" -Force

Write-Output "Downloading application executable..."
Invoke-WebRequest "https://github.com/aiden-deloryn/Hoist/releases/download/v${version}/hoist_${version}_${arch}.exe" -OutFile "${appPath}\${exeName}"

If (-not ($splitPath -contains $appPath)) {
    Write-Output "Adding application directory to PATH..."
    $splitPath += $appPath
    $newPath = ($splitPath -join ";")
    [System.Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
}
