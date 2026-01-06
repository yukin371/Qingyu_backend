@echo off
chcp 65001 >nul
echo ╔════════════════════════════════════════╗
echo ║   青羽写作平台 - 测试数据更新工具    ║
echo ╚════════════════════════════════════════╝
echo.

cd /d "%~dp0.."
echo 当前目录: %CD%
echo.

echo 正在编译并运行数据更新工具...
echo.

go run cmd/seed_data/main.go

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ 运行失败，错误代码: %ERRORLEVEL%
    echo.
    pause
    exit /b %ERRORLEVEL%
)

echo.
echo ✓ 运行完成
echo.
pause
