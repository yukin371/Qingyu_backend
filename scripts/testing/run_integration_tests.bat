@echo off
REM ========================================
REM 青羽后端 - 一键集成测试脚本
REM 用途：完整执行所有集成测试场景
REM ========================================

echo.
echo ========================================
echo 青羽后端 - 集成测试套件
echo ========================================
echo.

REM 检查Go环境
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到Go环境，请先安装Go
    pause
    exit /b 1
)

echo [环境检查] Go环境正常
echo.

REM 询问是否需要清理和导入数据
set /p prepare_data="是否需要准备测试数据？(yes/no，首次运行选yes): "

if "%prepare_data%"=="yes" (
    echo.
    echo ========================================
    echo 阶段1: 准备测试数据
    echo ========================================
    echo.

    echo [1/3] 清理旧测试数据...
    call scripts\testing\cleanup_test_data.bat
    if errorlevel 1 (
        echo [错误] 数据清理失败
        pause
        exit /b 1
    )

    echo.
    echo [2/3] 导入测试用户...
    go run scripts\testing\import_test_users.go
    if errorlevel 1 (
        echo [错误] 用户导入失败
        pause
        exit /b 1
    )

    echo.
    echo [3/3] 导入小说数据...
    call scripts\testing\import_novels.bat
    if errorlevel 1 (
        echo [错误] 小说数据导入失败
        pause
        exit /b 1
    )

    echo.
    echo [成功] 测试数据准备完成
    echo.
)

REM 检查服务器是否在运行
echo ========================================
echo 阶段2: 检查服务器状态
echo ========================================
echo.

curl -s http://localhost:8080/api/v1/system/health >nul 2>&1
if errorlevel 1 (
    echo [警告] 服务器未运行，正在启动...
    echo.

    REM 后台启动服务器
    start /B cmd /c "go run cmd/server/main.go > server_test.log 2>&1"

    echo [等待] 服务器启动中（等待10秒）...
    timeout /t 10 /nobreak >nul

    REM 再次检查
    curl -s http://localhost:8080/api/v1/system/health >nul 2>&1
    if errorlevel 1 (
        echo [错误] 服务器启动失败，请手动启动
        echo 命令: go run cmd/server/main.go
        pause
        exit /b 1
    )

    echo [成功] 服务器已启动
) else (
    echo [成功] 服务器正在运行
)

echo.
echo ========================================
echo 阶段3: 执行集成测试
echo ========================================
echo.

REM 设置测试配置
set GO_ENV=test

echo 可用的测试场景:
echo   1. 书城流程测试 (scenario_bookstore_test.go)
echo   2. 搜索功能测试 (scenario_search_test.go)
echo   3. 阅读流程测试 (scenario_reading_test.go)
echo   4. AI生成测试 (scenario_ai_generation_test.go)
echo   5. 认证流程测试 (scenario_auth_test.go)
echo   6. 写作流程测试 (scenario_writing_test.go)
echo   7. 互动功能测试 (scenario_interaction_test.go)
echo   8. 全部测试
echo.

set /p test_choice="请选择要执行的测试 (1-8): "

if "%test_choice%"=="1" (
    echo.
    echo [执行] 书城流程测试...
    go test -v ./test/integration/scenario_bookstore_test.go
) else if "%test_choice%"=="2" (
    echo.
    echo [执行] 搜索功能测试...
    go test -v ./test/integration/scenario_search_test.go
) else if "%test_choice%"=="3" (
    echo.
    echo [执行] 阅读流程测试...
    go test -v ./test/integration/scenario_reading_test.go
) else if "%test_choice%"=="4" (
    echo.
    echo [执行] AI生成测试...
    go test -v ./test/integration/scenario_ai_generation_test.go
) else if "%test_choice%"=="5" (
    echo.
    echo [执行] 认证流程测试...
    go test -v ./test/integration/scenario_auth_test.go
) else if "%test_choice%"=="6" (
    echo.
    echo [执行] 写作流程测试...
    go test -v ./test/integration/scenario_writing_test.go
) else if "%test_choice%"=="7" (
    echo.
    echo [执行] 互动功能测试...
    go test -v ./test/integration/scenario_interaction_test.go
) else if "%test_choice%"=="8" (
    echo.
    echo [执行] 全部集成测试...
    go test -v ./test/integration/scenario_*.go
) else (
    echo [错误] 无效的选择
    pause
    exit /b 1
)

echo.
echo ========================================
echo 测试完成
echo ========================================
echo.

REM 询问是否查看详细日志
set /p view_log="是否查看服务器日志？(yes/no): "
if "%view_log%"=="yes" (
    if exist server_test.log (
        type server_test.log
    ) else (
        echo 日志文件不存在
    )
)

echo.
echo 提示：
echo - 测试结果已显示在上方
echo - 服务器日志: server_test.log
echo - 如需停止服务器，请手动关闭或使用 Ctrl+C
echo.

pause



