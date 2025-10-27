@echo off
REM 青羽后端 - 运行所有场景测试
REM Windows 批处理脚本

echo.
echo ============================================================
echo 青羽后端 - 场景测试执行
echo ============================================================
echo.

REM 检查服务器是否运行
curl -s http://localhost:8080/api/v1/system/health >nul 2>&1
if %errorlevel% neq 0 (
    echo [错误] 服务器未运行，请先启动服务器
    echo 命令: go run cmd/server/main.go
    echo.
    pause
    exit /b 1
)

echo [成功] 服务器正在运行
echo.

REM 运行所有场景测试
echo 开始运行场景测试...
echo.

go test -v ./test/integration/ -run Scenario

echo.
echo ============================================================
echo 测试完成
echo ============================================================
echo.

pause



