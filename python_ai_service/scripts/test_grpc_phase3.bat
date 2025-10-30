@echo off
REM 测试Phase3 gRPC服务

echo ========================================
echo Phase3 gRPC服务测试脚本
echo ========================================

cd /d %~dp0..

echo.
echo 📋 测试前检查:
echo   1. 确保gRPC服务器正在运行
echo   2. 确保GOOGLE_API_KEY已设置
echo.
echo 如果服务器未运行，请先执行:
echo   scripts\start_grpc_server.bat
echo.
pause

echo.
echo 🧪 开始测试...
echo.
python tests/test_grpc_phase3.py

echo.
echo ========================================
echo 测试完成
echo ========================================
pause

