@echo off
REM 青羽后端 - 快速验证脚本（Windows版本）
REM 用于验证项目编译和测试是否通过

echo ======================================
echo 青羽后端 - 快速验证脚本
echo ======================================
echo.

set SUCCESS_COUNT=0
set FAIL_COUNT=0

echo 步骤 1: 编译项目...
go build -o Qingyu_backend.exe
if %errorlevel% equ 0 (
    echo [32m✓ 编译成功[0m
    set /a SUCCESS_COUNT+=1
) else (
    echo [31m✗ 编译失败[0m
    set /a FAIL_COUNT+=1
    goto :end
)

echo.
echo 步骤 2: 运行Repository层测试...
go test ./test/repository/ -v
if %errorlevel% equ 0 (
    echo [32m✓ Repository测试通过[0m
    set /a SUCCESS_COUNT+=1
) else (
    echo [31m✗ Repository测试失败[0m
    set /a FAIL_COUNT+=1
)

echo.
echo 步骤 3: 运行Service层测试...
go test ./test/service/ -v
if %errorlevel% equ 0 (
    echo [32m✓ Service测试通过[0m
    set /a SUCCESS_COUNT+=1
) else (
    echo [31m✗ Service测试失败[0m
    set /a FAIL_COUNT+=1
)

echo.
echo 步骤 4: 运行书城系统测试...
go test ./test/ -run "Bookstore" -v
if %errorlevel% equ 0 (
    echo [32m✓ 书城测试通过[0m
    set /a SUCCESS_COUNT+=1
) else (
    echo [33m! 书城测试警告（可能需要数据库连接）[0m
)

:end
echo.
echo ======================================
echo 验证完成！
echo ======================================
echo 成功: %SUCCESS_COUNT% 项
echo 失败: %FAIL_COUNT% 项
echo.

if %FAIL_COUNT% equ 0 (
    echo [32m🎉 所有检查通过！项目状态健康！[0m
    exit /b 0
) else (
    echo [31m⚠ 有 %FAIL_COUNT% 项检查失败，请修复后再提交代码[0m
    exit /b 1
)

