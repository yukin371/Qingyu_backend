@echo off
REM 生成gRPC Python代码

echo ========================================
echo 生成gRPC Protobuf代码
echo ========================================

cd /d %~dp0..

echo.
echo [1/3] 检查依赖...
python -c "import grpc_tools.protoc" 2>NUL
if errorlevel 1 (
    echo ❌ 未安装grpcio-tools，正在安装...
    pip install grpcio-tools
)

echo.
echo [2/3] 生成Python protobuf代码...
python -m grpc_tools.protoc ^
    -I proto ^
    --python_out=src/grpc_service ^
    --grpc_python_out=src/grpc_service ^
    proto/ai_service.proto

if errorlevel 1 (
    echo ❌ 代码生成失败
    pause
    exit /b 1
)

echo.
echo [3/3] 修复导入路径...
REM Python生成的代码可能需要修复import路径
python -c "import os; pb_file='src/grpc_service/ai_service_pb2_grpc.py'; content=open(pb_file).read(); content=content.replace('import ai_service_pb2', 'from . import ai_service_pb2'); open(pb_file, 'w').write(content)" 2>NUL

echo.
echo ✅ 代码生成完成！
echo 生成的文件:
echo   - src/grpc_service/ai_service_pb2.py
echo   - src/grpc_service/ai_service_pb2_grpc.py
echo.
pause

