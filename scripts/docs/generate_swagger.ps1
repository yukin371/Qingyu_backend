$ErrorActionPreference = "Stop"

$RepoRoot = Split-Path -Parent (Split-Path -Parent $PSScriptRoot)
Set-Location $RepoRoot

$swaggerDirs = @(
    "api/v1",
    "pkg/response",
    "models",
    "models/dto",
    "service/interfaces",
    "service/ai/dto",
    "service/shared/storage",
    "service/shared/stats"
) -join ","

Write-Host "Generating Swagger artifacts..."
swag init `
    -g swagger.go `
    -d $swaggerDirs `
    --parseDependency=false `
    -o docs

Write-Host ""
Write-Host "Swagger artifacts updated:"
Write-Host "  docs/docs.go"
Write-Host "  docs/swagger.json"
Write-Host "  docs/swagger.yaml"
