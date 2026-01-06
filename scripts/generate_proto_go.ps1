# Windows PowerShell script to generate Go protobuf code

Write-Host "Generating Go protobuf code..." -ForegroundColor Green

# Check if protoc is installed
if (-not (Get-Command protoc -ErrorAction SilentlyContinue)) {
    Write-Host "Error: protoc not found. Please install Protocol Buffers compiler." -ForegroundColor Red
    Write-Host "Download from: https://github.com/protocolbuffers/protobuf/releases" -ForegroundColor Yellow
    exit 1
}

# Check if Go plugins are installed
$goPath = $env:GOPATH
if (-not $goPath) {
    $goPath = "$env:USERPROFILE\go"
}

$protoc_gen_go = Join-Path $goPath "bin\protoc-gen-go.exe"
$protoc_gen_go_grpc = Join-Path $goPath "bin\protoc-gen-go-grpc.exe"

if (-not (Test-Path $protoc_gen_go)) {
    Write-Host "Installing protoc-gen-go..." -ForegroundColor Yellow
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
}

if (-not (Test-Path $protoc_gen_go_grpc)) {
    Write-Host "Installing protoc-gen-go-grpc..." -ForegroundColor Yellow
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
}

# Create output directory if not exists
$outputDir = "pkg\grpc\pb"
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

# Generate Go code
Write-Host "Generating protobuf code..." -ForegroundColor Cyan

protoc --go_out=pkg\grpc\pb --go-grpc_out=pkg\grpc\pb `
    --go_opt=paths=source_relative `
    --go-grpc_opt=paths=source_relative `
    --proto_path=python_ai_service\proto `
    python_ai_service\proto\ai_service.proto

if ($LASTEXITCODE -eq 0) {
    Write-Host "✓ Go protobuf code generated successfully in pkg\grpc\pb\" -ForegroundColor Green
} else {
    Write-Host "✗ Failed to generate Go protobuf code" -ForegroundColor Red
    exit 1
}

