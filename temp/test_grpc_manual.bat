@echo off
echo ========================================
echo Testing Python gRPC Server
echo ========================================

echo.
echo [1/3] Killing existing Python processes...
taskkill /F /IM python.exe /T 2>nul

echo.
echo [2/3] Starting gRPC server...
cd python_ai_service
start /B python start_grpc_simple.py

echo.
echo [3/3] Waiting 5 seconds for server startup...
timeout /t 5 /nobreak

echo.
echo [4/3] Checking port 50052...
netstat -ano | findstr ":50052"

echo.
echo [5/3] Testing gRPC connection...
cd ..
go run test_grpc_connection.go

pause

