@echo off
REM ServerDemo 测试模式启动脚本
REM 此脚本会启动后端服务并跳过JWT认证，方便API测试

echo ========================================
echo   青羽写作平台 - 测试模式启动
echo ========================================
echo.
echo [测试模式] JWT认证已禁用
echo   所有API都可以直接访问，无需Token
echo.
echo 测试用户信息:
echo   - user_id: test-user-id
echo   - username: test-user
echo   - roles: [reader, author, admin]
echo.
echo ========================================
echo.

REM 设置环境变量
set SKIP_AUTH=true

REM 启动服务
serverdemo.exe

pause
