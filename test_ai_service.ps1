# AI服务测试脚本

$baseUrl = "http://localhost:8080/api/v1"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  AI服务测试" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# 等待服务器启动
Write-Host "1. 等待服务器启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# 登录获取token
Write-Host "2. 登录VIP用户..." -ForegroundColor Yellow
try {
    $loginBody = @{
        username = "vip_user01"
        password = "password123"
    } | ConvertTo-Json

    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/login" -Method Post -Body $loginBody -ContentType "application/json"

    if ($loginResponse.code -eq 0) {
        $token = $loginResponse.data.token
        Write-Host "   ✓ 登录成功" -ForegroundColor Green
        Write-Host "   Token: $($token.Substring(0, 20))..." -ForegroundColor Gray
    } else {
        Write-Host "   ✗ 登录失败: $($loginResponse.message)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "   ✗ 登录请求失败: $_" -ForegroundColor Red
    Write-Host "   提示: 请确保服务器已启动 (go run cmd/server/main.go)" -ForegroundColor Yellow
    exit 1
}

Write-Host ""

# 测试AI健康检查
Write-Host "3. 测试AI服务健康检查..." -ForegroundColor Yellow
try {
    $headers = @{
        "Authorization" = "Bearer $token"
    }

    $healthResponse = Invoke-RestMethod -Uri "$baseUrl/ai/health" -Method Get -Headers $headers

    if ($healthResponse.code -eq 0) {
        Write-Host "   ✓ AI服务健康" -ForegroundColor Green
        Write-Host "   状态: $($healthResponse.data.status)" -ForegroundColor Gray
    } else {
        Write-Host "   ✗ 健康检查失败: $($healthResponse.message)" -ForegroundColor Red
    }
} catch {
    Write-Host "   ✗ 健康检查请求失败: $_" -ForegroundColor Red
}

Write-Host ""

# 测试获取提供商列表
Write-Host "4. 获取AI提供商列表..." -ForegroundColor Yellow
try {
    $providersResponse = Invoke-RestMethod -Uri "$baseUrl/ai/providers" -Method Get -Headers $headers

    if ($providersResponse.code -eq 0) {
        Write-Host "   ✓ 获取成功" -ForegroundColor Green
        $providers = $providersResponse.data.providers
        foreach ($provider in $providers) {
            $statusIcon = if ($provider.enabled) { "✓" } else { "✗" }
            Write-Host "   $statusIcon $($provider.name) - 优先级: $($provider.priority)" -ForegroundColor Gray
        }
    }
} catch {
    Write-Host "   ✗ 获取提供商失败: $_" -ForegroundColor Red
}

Write-Host ""

# 测试AI文本生成（续写）
Write-Host "5. 测试AI文本续写..." -ForegroundColor Yellow
try {
    $generateBody = @{
        text = "在一个遥远的王国里，有一位勇敢的骑士"
        max_tokens = 100
        temperature = 0.7
    } | ConvertTo-Json

    $generateResponse = Invoke-RestMethod -Uri "$baseUrl/ai/generate" -Method Post -Body $generateBody -ContentType "application/json" -Headers $headers

    if ($generateResponse.code -eq 0) {
        Write-Host "   ✓ 续写成功" -ForegroundColor Green
        Write-Host "   生成文本: $($generateResponse.data.generated_text.Substring(0, [Math]::Min(100, $generateResponse.data.generated_text.Length)))..." -ForegroundColor Gray
        if ($generateResponse.data.usage) {
            Write-Host "   使用Token: $($generateResponse.data.usage.total_tokens)" -ForegroundColor Gray
        }
    } else {
        Write-Host "   ✗ 续写失败: $($generateResponse.message)" -ForegroundColor Red
        Write-Host "   错误: $($generateResponse.error)" -ForegroundColor Red

        # 如果是配额问题，检查配额状态
        if ($generateResponse.code -eq 429) {
            Write-Host ""
            Write-Host "   检查配额状态..." -ForegroundColor Yellow

            # 直接查询MongoDB
            Write-Host "   正在运行配额检查工具..." -ForegroundColor Gray
            go run cmd/check_quota/main.go
        }
    }
} catch {
    Write-Host "   ✗ 续写请求失败: $_" -ForegroundColor Red
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  测试完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

