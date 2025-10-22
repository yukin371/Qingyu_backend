@echo off
chcp 65001 >nul
echo ========================================
echo 项目修复验证脚本
echo ========================================
echo.

cd /d "%~dp0"

echo [1] 检查 Go 环境...
go version
if %errorlevel% neq 0 (
    echo 错误: Go 环境未安装或未配置
    pause
    exit /b 1
)
echo.

echo [2] 编译项目...
cd ../
go build -o Qingyu_backend.exe .
if %errorlevel% neq 0 (
    echo 编译失败! 请检查错误信息
    pause
    exit /b 1
)
echo 编译成功!
echo.

echo [3] 启动项目...
echo 按 Ctrl+C 停止服务器
echo.
.\Qingyu_backend.exe

pause

