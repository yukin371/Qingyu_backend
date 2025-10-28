# Windows PowerShell script to generate Python protobuf code

Write-Host "Generating Python protobuf code..." -ForegroundColor Green

# Check if grpc_tools is installed
python -c "import grpc_tools" 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "Error: grpc_tools not installed." -ForegroundColor Red
    Write-Host "Installing grpcio-tools..." -ForegroundColor Yellow
    pip install grpcio-tools
}

# Change to python_ai_service directory
Push-Location python_ai_service

# Create output directory if not exists
if (-not (Test-Path "src\grpc_server")) {
    New-Item -ItemType Directory -Path "src\grpc_server" -Force | Out-Null
}

# Generate Python code
Write-Host "Generating protobuf code..." -ForegroundColor Cyan

python -m grpc_tools.protoc -I proto `
    --python_out=src\grpc_server `
    --grpc_python_out=src\grpc_server `
    proto\ai_service.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Python protobuf code generated successfully in src\grpc_server\" -ForegroundColor Green

    # Fix import paths
    Write-Host "Fixing import paths..." -ForegroundColor Cyan
    $grpcFile = "src\grpc_server\ai_service_pb2_grpc.py"
    if (Test-Path $grpcFile) {
        $content = Get-Content $grpcFile -Raw
        $content = $content -replace 'import ai_service_pb2', 'from . import ai_service_pb2'
        Set-Content -Path $grpcFile -Value $content
        Write-Host "✓ Import paths fixed" -ForegroundColor Green
    }
} else {
    Write-Host "✗ Failed to generate Python protobuf code" -ForegroundColor Red
    Pop-Location
    exit 1
}

Pop-Location

