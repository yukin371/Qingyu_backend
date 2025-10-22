# PowerShell 脚本：使用 Docker 运行用户 Repository 集成测试

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  用户管理模块 Repository 集成测试" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 1. 检查 Docker 是否运行
Write-Host "[1/5] 检查 Docker 服务..." -ForegroundColor Yellow
$dockerRunning = docker ps 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Docker 未运行或未安装" -ForegroundColor Red
    Write-Host "   请先启动 Docker Desktop" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Docker 服务正常" -ForegroundColor Green
Write-Host ""

# 2. 启动数据库服务
Write-Host "[2/5] 启动 Docker 数据库服务..." -ForegroundColor Yellow
Push-Location ../../../docker
docker-compose -f docker-compose.db-only.yml up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ 启动 Docker 服务失败" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location
Write-Host "✅ Docker 服务已启动" -ForegroundColor Green
Write-Host ""

# 3. 等待 MongoDB 就绪
Write-Host "[3/5] 等待 MongoDB 就绪..." -ForegroundColor Yellow
$retries = 0
$maxRetries = 30
$ready = $false

while ($retries -lt $maxRetries) {
    $mongoStatus = docker exec qingyu-mongodb mongosh --eval "db.runCommand('ping').ok" --quiet 2>$null
    if ($LASTEXITCODE -eq 0 -and $mongoStatus -eq "1") {
        $ready = $true
        break
    }
    $retries++
    Write-Host "   等待中... ($retries/$maxRetries)" -ForegroundColor Gray
    Start-Sleep -Seconds 1
}

if (-not $ready) {
    Write-Host "❌ MongoDB 启动超时" -ForegroundColor Red
    Write-Host "   请检查: docker logs qingyu-mongodb" -ForegroundColor Red
    exit 1
}
Write-Host "✅ MongoDB 已就绪" -ForegroundColor Green
Write-Host ""

# 4. 运行测试
Write-Host "[4/5] 运行集成测试..." -ForegroundColor Yellow
Write-Host ""
Push-Location ../../..
go test -v ./test/repository/user/...
$testResult = $LASTEXITCODE
Pop-Location
Write-Host ""

if ($testResult -eq 0) {
    Write-Host "✅ 所有测试通过" -ForegroundColor Green
} else {
    Write-Host "❌ 部分测试失败" -ForegroundColor Red
}
Write-Host ""

# 5. 询问是否停止服务
Write-Host "[5/5] 清理..." -ForegroundColor Yellow
$cleanup = Read-Host "是否停止 Docker 服务? (y/N)"
if ($cleanup -eq "y" -or $cleanup -eq "Y") {
    Push-Location ../../../docker
    docker-compose -f docker-compose.db-only.yml down
    Pop-Location
    Write-Host "✅ Docker 服务已停止" -ForegroundColor Green
} else {
    Write-Host "ℹ️  Docker 服务继续运行" -ForegroundColor Cyan
    Write-Host "   手动停止: cd docker && docker-compose -f docker-compose.db-only.yml down" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  测试完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

exit $testResult

