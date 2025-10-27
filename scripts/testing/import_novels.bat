@echo off
REM ========================================
REM 导入小说测试数据脚本
REM 用途：从 novels_100.json 导入100本测试小说
REM ========================================

echo ========================================
echo 青羽后端 - 小说数据导入
echo ========================================
echo.

REM 检查数据文件是否存在
if not exist "data\novels_100.json" (
    echo [错误] 找不到数据文件: data\novels_100.json
    echo.
    echo 请确保：
    echo 1. 已运行 Python 脚本生成数据文件
    echo 2. 文件路径正确
    echo.
    pause
    exit /b 1
)

echo [1/3] 数据文件检查...
echo [成功] 找到数据文件: data\novels_100.json
echo.

REM 先试运行验证数据
echo [2/3] 验证数据格式（试运行）...
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json -dry-run=true
if errorlevel 1 (
    echo.
    echo [错误] 数据验证失败，请检查数据格式
    pause
    exit /b 1
)
echo.

REM 正式导入
echo [3/3] 正式导入数据...
echo.
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json

if errorlevel 1 (
    echo.
    echo [错误] 数据导入失败
    pause
    exit /b 1
)

echo.
echo ========================================
echo ✓ 小说数据导入成功
echo ========================================
echo.
echo 提示：
echo - 可以通过 API 访问书城查看导入的书籍
echo - 运行集成测试验证导入结果
echo.
pause


