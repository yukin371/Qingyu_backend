@echo off
chcp 65001 >nul
cd /d "%~dp0"

echo ========================================
echo 青羽后端 - 开发环境启动
echo ========================================
echo.

echo [1] 启动后端服务（含MongoDB和Redis）...
docker-compose -f docker-compose.dev.yml up -d

echo.
echo [2] 等待服务启动...
timeout /t 5 >nul

echo.
echo [3] 显示服务状态
docker-compose -f docker-compose.dev.yml ps

echo.
echo ========================================
echo 后端服务已启动！
echo ========================================
echo 后端API: http://localhost:8080
echo MongoDB: localhost:27017
echo Redis:   localhost:6379
echo ========================================
echo.
echo 查看日志: docker-compose -f docker-compose.dev.yml logs -f
echo 停止服务: docker-compose -f docker-compose.dev.yml down
echo ========================================

pause



