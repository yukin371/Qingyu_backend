@echo off
REM ============================================
REM Go后端 与 Python AI微服务 集成测试脚本
REM ============================================

echo.
echo ========================================
echo Go Backend ^<-^> Python AI 微服务集成测试
echo ========================================
echo.

REM 检查环境变量
if "%GOOGLE_API_KEY%"=="" (
    echo [错误] GOOGLE_API_KEY 环境变量未设置
    echo.
    echo 请先设置 API Key:
    echo   set GOOGLE_API_KEY=your_api_key_here
    echo.
    pause
    exit /b 1
)

echo [✓] GOOGLE_API_KEY 已设置
echo.

REM 进入项目根目录
cd /d %~dp0..\..

echo ========================================
echo 步骤 1/4: 检查依赖
echo ========================================
echo.

REM 检查 Go 依赖
echo [1/3] 检查 Go 依赖...
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] Go 未安装或不在 PATH 中
    pause
    exit /b 1
)
echo [✓] Go 已安装
echo.

REM 检查 Python 依赖
echo [2/3] 检查 Python 依赖...
python --version >nul 2>&1
if errorlevel 1 (
    echo [错误] Python 未安装或不在 PATH 中
    pause
    exit /b 1
)
echo [✓] Python 已安装
echo.

REM 检查 Python AI 服务依赖
echo [3/3] 检查 Python AI 服务依赖...
cd python_ai_service
python -c "import grpc; import google.generativeai" >nul 2>&1
if errorlevel 1 (
    echo [警告] Python 依赖不完整，正在安装...
    pip install -r requirements.txt
    if errorlevel 1 (
        echo [错误] 依赖安装失败
        cd ..
        pause
        exit /b 1
    )
)
echo [✓] Python 依赖已安装
cd ..
echo.

echo ========================================
echo 步骤 2/4: 启动 Python AI 微服务
echo ========================================
echo.

echo [启动] gRPC 服务器 (端口: 50051)...
echo.

REM 在新窗口启动 Python gRPC 服务器
start "Python AI gRPC Server" cmd /k "cd /d %CD%\python_ai_service && python run_grpc_server.py"

echo [等待] 等待服务器启动 (5秒)...
timeout /t 5 /nobreak >nul
echo.

echo ========================================
echo 步骤 3/4: 运行 Python 端测试
echo ========================================
echo.

echo [测试] Python gRPC 客户端测试...
cd python_ai_service
python tests\test_grpc_phase3.py
set PYTHON_TEST_RESULT=%ERRORLEVEL%
cd ..
echo.

if %PYTHON_TEST_RESULT% neq 0 (
    echo [错误] Python 测试失败
    echo 请检查 gRPC 服务器日志
    echo.
    pause
    exit /b 1
)

echo [✓] Python 测试通过
echo.

echo ========================================
echo 步骤 4/4: 运行 Go 端集成测试
echo ========================================
echo.

echo [测试] Go gRPC 客户端测试...
echo.

REM 运行 Go 测试
go run cmd\test_phase3_grpc\main.go --addr localhost:50051
set GO_TEST_RESULT=%ERRORLEVEL%

if %GO_TEST_RESULT% neq 0 (
    echo.
    echo [错误] Go 测试失败
    echo 请检查 gRPC 服务器日志
    echo.
    pause
    exit /b 1
)

echo.
echo [✓] Go 测试通过
echo.

echo ========================================
echo 测试完成摘要
echo ========================================
echo.
echo [✓] Python AI 微服务启动成功
echo [✓] Python 客户端测试通过
echo [✓] Go 客户端测试通过
echo.
echo ========================================
echo ✅ 所有测试通过！
echo ========================================
echo.
echo 注意: Python AI 服务仍在运行
echo 请手动关闭窗口: "Python AI gRPC Server"
echo.
pause

