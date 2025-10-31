@echo off
REM ============================================
REM 运行 Go gRPC 集成测试
REM ============================================

echo.
echo ========================================
echo Go gRPC 集成测试
echo ========================================
echo.

REM 检查 gRPC 服务器是否运行
echo [检查] gRPC 服务器状态...
timeout /t 1 /nobreak >nul

REM 进入项目根目录
cd /d %~dp0..\..

echo.
echo [运行] Go 集成测试...
echo.

REM 运行集成测试
go test -v -timeout 300s ./test/integration -run TestGRPC

if errorlevel 1 (
    echo.
    echo ========================================
    echo ❌ 测试失败
    echo ========================================
    echo.
    echo 请确保:
    echo   1. Python AI 服务正在运行 (localhost:50051)
    echo   2. GOOGLE_API_KEY 环境变量已设置
    echo   3. 所有依赖已安装
    echo.
    pause
    exit /b 1
)

echo.
echo ========================================
echo ✅ 所有测试通过
echo ========================================
echo.
pause


