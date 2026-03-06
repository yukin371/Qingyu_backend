Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $PSScriptRoot

git -C $root config core.hooksPath .githooks | Out-Null

Write-Host "Git hooks installed."
Write-Host ("core.hooksPath=" + (git -C $root config --get core.hooksPath))
