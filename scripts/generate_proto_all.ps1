# Windows PowerShell script to generate all protobuf code

Write-Host "=== Generating All Protobuf Code ===" -ForegroundColor Cyan
Write-Host ""

# Generate Go code
Write-Host "[1/2] Generating Go protobuf code..." -ForegroundColor Yellow
& "$PSScriptRoot\generate_proto_go.ps1"

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to generate Go protobuf code" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Generate Python code
Write-Host "[2/2] Generating Python protobuf code..." -ForegroundColor Yellow
& "$PSScriptRoot\generate_proto_python.ps1"

if ($LASTEXITCODE -ne 0) {
    Write-Host "Failed to generate Python protobuf code" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "=== All protobuf code generated successfully! ===" -ForegroundColor Green

