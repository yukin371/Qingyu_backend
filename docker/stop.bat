@echo off
chcp 65001 >nul
cd /d "%~dp0"

echo ========================================
echo 青羽后端 - 停止服务
echo ========================================
echo.

echo 正在停止后端服务（含MongoDB和Redis）...
docker-compose -f docker-compose.dev.yml down

echo.
echo ========================================
echo 后端服务已停止！
echo ========================================

pause



