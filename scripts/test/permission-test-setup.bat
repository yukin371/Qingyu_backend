@echo off
REM ============================================
REM Qingyu Backend - 权限系统测试环境准备脚本 (Windows)
REM ============================================
REM
REM 功能:
REM   1. 检查MongoDB连接状态
REM   2. 检查Redis连接状态
REM   3. 创建/重置测试数据库
REM   4. 初始化测试数据
REM
REM 使用方法:
REM   scripts\test\permission-test-setup.bat
REM   scripts\test\permission-test-setup.bat --skip-data
REM   scripts\test\permission-test-setup.bat --db-only
REM
REM ============================================

setlocal enabledelayedexpansion

REM 默认配置
if "%QINGYU_TEST_DB_NAME%"=="" set QINGYU_TEST_DB_NAME=qingyu_permission_test
if "%QINGYU_MONGO_HOST%"=="" set QINGYU_MONGO_HOST=localhost
if "%QINGYU_MONGO_PORT%"=="" set QINGYU_MONGO_PORT=27017
if "%QINGYU_REDIS_HOST%"=="" set QINGYU_REDIS_HOST=localhost
if "%QINGYU_REDIS_PORT%"=="" set QINGYU_REDIS_PORT=6379

set SKIP_DATA=0
set DB_ONLY=0
set FORCE_RECREATE=0

REM 解析参数
:parse_args
if "%~1"=="--skip-data" (
    set SKIP_DATA=1
    shift
    goto parse_args
)
if "%~1"=="--db-only" (
    set DB_ONLY=1
    shift
    goto parse_args
)
if "%~1"=="--force" (
    set FORCE_RECREATE=1
    shift
    goto parse_args
)
if "%~1"=="--help" (
    goto show_help
)

REM ==================== 主流程 ====================

echo.
echo ========================================
echo   Qingyu Backend - 权限系统测试环境准备
echo ========================================
echo.

REM 步骤1: 检查MongoDB
echo [1/4] 检查MongoDB连接...
echo.

where mongosh >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    set MONGO_CMD=mongosh
) else (
    where mongo >nul 2>&1
    if %ERRORLEVEL% EQU 0 (
        set MONGO_CMD=mongo
    ) else (
        echo [错误] 未找到MongoDB客户端 (mongosh 或 mongo)
        echo 请安装MongoDB客户端: https://www.mongodb.com/try/download
        goto error_exit
    )
)

REM 尝试连接MongoDB
%MONGO_CMD% --quiet --host %QINGYU_MONGO_HOST% --port %QINGYU_MONGO_PORT% --eval "db.version()" >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    echo [成功] MongoDB连接成功
    echo   主机: %QINGYU_MONGO_HOST%:%QINGYU_MONGO_PORT%
    echo.
) else (
    echo [错误] 无法连接到MongoDB
    echo 请确保MongoDB服务已启动:
    echo   Windows: net start MongoDB
    goto error_exit
)

REM 步骤2: 检查Redis
if %DB_ONLY% EQU 1 (
    echo [信息] 跳过Redis检查 (--db-only 模式)
    echo.
    goto skip_redis
)

echo [2/4] 检查Redis连接...
echo.

where redis-cli >nul 2>&1
if %ERRORLEVEL% EQU 0 (
    redis-cli -h %QINGYU_REDIS_HOST% -p %QINGYU_REDIS_PORT% ping >nul 2>&1
    if !ERRORLEVEL! EQU 0 (
        echo [成功] Redis连接成功
        echo   主机: %QINGYU_REDIS_HOST%:%QINGYU_REDIS_PORT%
        echo.
    ) else (
        echo [警告] 无法连接到Redis
        echo 请确保Redis服务已启动或使用Docker:
        echo   docker run -d -p 6379:6379 redis:alpine
        echo.
    )
) else (
    echo [警告] 未找到redis-cli命令
    echo 请安装Redis客户端或使用Docker启动Redis
    echo.
)

:skip_redis

REM 步骤3: 准备数据库
echo [3/4] 准备测试数据库...
echo.
echo   数据库名称: %QINGYU_TEST_DB_NAME%

if %FORCE_RECREATE% EQU 1 (
    echo   强制重建数据库...
    %MONGO_CMD% --quiet --host %QINGYU_MONGO_HOST% --port %QINGYU_MONGO_PORT% --eval "db.getSiblingDB('%QINGYU_TEST_DB_NAME%').dropDatabase()"
    echo   数据库已删除
)

echo   创建集合和索引...

REM 创建roles集合
%MONGO_CMD% --quiet --host %QINGYU_MONGO_HOST% --port %QINGYU_MONGO_PORT% %QINGYU_TEST_DB_NAME% --eval "db.createCollection('roles'); db.roles.createIndex({ name: 1 }, { unique: true });" >nul 2>&1

REM 创建users集合
%MONGO_CMD% --quiet --host %QINGYU_MONGO_HOST% --port %QINGYU_MONGO_PORT% %QINGYU_TEST_DB_NAME% --eval "if (!db.getCollectionNames().includes('users')) { db.createCollection('users'); db.users.createIndex({ username: 1 }, { unique: true }); }" >nul 2>&1

echo [成功] 数据库准备完成
echo.

REM 步骤4: 初始化测试数据
if %SKIP_DATA% EQU 1 (
    echo [信息] 跳过测试数据填充 (--skip-data)
    echo.
    goto skip_data
)

echo [4/4] 初始化测试数据...
echo.

where go >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [错误] 未找到Go命令
    echo 请安装Go: https://golang.org/dl/
    goto error_exit
)

if not exist "go.mod" (
    echo [错误] 请在项目根目录运行此脚本
    goto error_exit
)

echo   运行测试数据填充脚本...
go run scripts/test/permission-test-data.go --db=%QINGYU_TEST_DB_NAME%
if %ERRORLEVEL% EQU 0 (
    echo [成功] 测试数据填充完成
    echo.
) else (
    echo [错误] 测试数据填充失败
    goto error_exit
)

:skip_data

REM 打印摘要
echo.
echo ========================================
echo   测试环境准备完成！
echo ========================================
echo.
echo 测试环境信息:
echo   数据库: %QINGYU_TEST_DB_NAME%
echo   MongoDB: %QINGYU_MONGO_HOST%:%QINGYU_MONGO_PORT%
echo   Redis: %QINGYU_REDIS_HOST%:%QINGYU_REDIS_PORT%
echo.

if %SKIP_DATA% EQU 0 (
    echo 测试账号:
    echo   管理员: admin@test.com / Admin@123
    echo   作者:   author@test.com / Author@123
    echo   读者:   reader@test.com / Reader@123
    echo   编辑:   editor@test.com / Editor@123
    echo.
)

echo 下一步:
echo   1. 设置环境变量:
echo      set QINGYU_DATABASE_NAME=%QINGYU_TEST_DB_NAME%
echo.
echo   2. 启动测试服务器:
echo      go run cmd/server/main.go
echo.
echo   3. 运行权限测试:
echo      go test ./internal/middleware/auth/... -v
echo.
echo ========================================
echo.

goto end

:show_help
echo 使用方法: %~nx0 [选项]
echo.
echo 选项:
echo   --skip-data    跳过测试数据填充
echo   --db-only      仅准备数据库（不检查Redis）
echo   --force        强制重建数据库
echo   --help         显示此帮助信息
echo.
echo 环境变量:
echo   QINGYU_TEST_DB_NAME      测试数据库名称 (默认: qingyu_permission_test)
echo   QINGYU_MONGO_HOST        MongoDB主机 (默认: localhost)
echo   QINGYU_MONGO_PORT        MongoDB端口 (默认: 27017)
echo   QINGYU_REDIS_HOST        Redis主机 (默认: localhost)
echo   QINGYU_REDIS_PORT        Redis端口 (默认: 6379)
echo.
echo 示例:
echo   # 完整设置（包括数据）
echo   %~nx0
echo.
echo   # 仅检查数据库
echo   %~nx0 --db-only
echo.
echo   # 仅设置数据库，不填充数据
echo   %~nx0 --skip-data
echo.
exit /B 0

:error_exit
echo.
echo [错误] 设置失败，请检查错误信息
exit /B 1

:end
echo [成功] 所有步骤完成！
exit /B 0
