@echo off
REM 启动Phase3 gRPC服务器

echo ========================================
echo Phase3 gRPC服务器启动脚本
echo ========================================

cd /d %~dp0..

REM 检查API密钥
if "%GOOGLE_API_KEY%"=="" (
    echo.
    echo ❌ 错误: 未设置GOOGLE_API_KEY环境变量
    echo.
    echo 请先设置API密钥:
    echo   set GOOGLE_API_KEY=your_api_key_here
    echo.
    pause
    exit /b 1
)

echo.
echo ✅ API密钥已设置
echo.
echo [1/2] 检查依赖...
python -c "import grpc; import google.generativeai" 2>NUL
if errorlevel 1 (
    echo ❌ 依赖缺失，正在安装...
    pip install -r requirements.txt
)

echo.
echo [2/2] 启动gRPC服务器...
echo.
echo ========================================
python src/grpc_service/server.py --host 0.0.0.0 --port 50051

pause

