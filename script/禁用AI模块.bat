@echo off
chcp 65001 >nul
echo ========================================
echo 禁用 AI 模块的 *_new.go 文件
echo ========================================

cd /d "%~dp0.."

echo.
echo [步骤1] 禁用 adapter_manager_new.go...
if exist "service\ai\adapter_manager_new.go" (
    move /Y "service\ai\adapter_manager_new.go" "service\ai\adapter_manager_new.go.disabled" >nul 2>&1
    echo [成功] adapter_manager_new.go 已禁用
) else (
    echo [跳过] adapter_manager_new.go 文件不存在或已禁用
)

echo.
echo [步骤2] 禁用 context_service_new.go...
if exist "service\ai\context_service_new.go" (
    move /Y "service\ai\context_service_new.go" "service\ai\context_service_new.go.disabled" >nul 2>&1
    echo [成功] context_service_new.go 已禁用
) else (
    echo [跳过] context_service_new.go 文件不存在或已禁用
)

echo.
echo [步骤3] 禁用 external_api_service_new.go...
if exist "service\ai\external_api_service_new.go" (
    move /Y "service\ai\external_api_service_new.go" "service\ai\external_api_service_new.go.disabled" >nul 2>&1
    echo [成功] external_api_service_new.go 已禁用
) else (
    echo [跳过] external_api_service_new.go 文件不存在或已禁用
)

echo.
echo [步骤4] 编译检查...
go build -o nul . 2>compile_errors.txt
if %ERRORLEVEL% EQU 0 (
    echo [成功] 编译通过！
    del compile_errors.txt
    echo.
    echo ========================================
    echo AI 模块已成功禁用，项目编译通过
    echo ========================================
) else (
    echo [失败] 编译仍有错误，请查看 compile_errors.txt
    type compile_errors.txt
    echo.
    echo ========================================
    echo 编译失败，请检查错误
    echo ========================================
)

pause

