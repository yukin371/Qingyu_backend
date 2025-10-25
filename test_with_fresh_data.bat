@echo off
echo =========================================
echo   准备测试数据并运行AI测试
echo =========================================
echo.

echo 1. 准备测试数据...
go run cmd/prepare_test_data/main.go
if errorlevel 1 (
    echo 数据准备失败！
    pause
    exit /b 1
)

echo.
echo 2. 运行AI测试...
go test ./test/integration -run "^TestAIGenerationScenario$" -v -count=1

echo.
pause

