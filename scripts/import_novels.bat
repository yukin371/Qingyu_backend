@echo off
chcp 65001 >nul
echo ╔════════════════════════════════════════╗
echo ║    青羽写作平台 - 小说批量导入工具     ║
echo ╚════════════════════════════════════════╝
echo.

cd /d "%~dpdp0"
echo 当前目录: %CD%
echo.

echo 正在编译并运行小说导入工具...
echo.

echo 1 > import_input.txt
echo 5 >> import_input.txt

go run cmd/import_novels/main.go < import_input.txt

if %ERRORLEVEL% NEQ 0 (
    echo.
    echo ❌ 运行失败，错误代码: %ERRORLEVEL%
    echo.
    pause
    exit /b %ERRORLEVEL%
)

del import_input.txt

echo.
echo ✓ 运行完成
echo.
pause
