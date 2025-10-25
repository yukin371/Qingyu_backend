@echo off
REM 青羽后端 - 运行单个场景测试
REM Windows 批处理脚本

setlocal enabledelayedexpansion

echo.
echo ============================================================
echo 青羽后端 - 场景测试选择
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

echo 可用的测试场景:
echo   1. 书城流程测试 (TestBookstoreScenario)
echo   2. 搜索功能测试 (TestSearchScenario)
echo   3. 阅读流程测试 (TestReadingScenario)
echo   4. AI生成测试 (TestAIGenerationScenario)
echo   5. 认证流程测试 (TestAuthScenario)
echo   6. 写作流程测试 (TestWritingScenario)
echo   7. 互动功能测试 (TestInteractionScenario)
echo   8. 全部测试
echo   0. 退出
echo.

set /p choice="请选择要执行的测试 (0-8): "

if "%choice%"=="0" goto :end
if "%choice%"=="1" set testname=TestBookstoreScenario
if "%choice%"=="2" set testname=TestSearchScenario
if "%choice%"=="3" set testname=TestReadingScenario
if "%choice%"=="4" set testname=TestAIGenerationScenario
if "%choice%"=="5" set testname=TestAuthScenario
if "%choice%"=="6" set testname=TestWritingScenario
if "%choice%"=="7" set testname=TestInteractionScenario
if "%choice%"=="8" set testname=Scenario

if not defined testname (
    echo.
    echo [错误] 无效的选择
    pause
    exit /b 1
)

echo.
echo 运行测试: %testname%
echo.

go test -v ./test/integration/ -run %testname%

echo.
echo ============================================================
echo 测试完成
echo ============================================================
echo.

:end
pause



