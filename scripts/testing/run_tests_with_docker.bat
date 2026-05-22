@echo off
setlocal enabledelayedexpansion

echo 🚀 启动测试环境...

:: 清理旧的测试环境
docker-compose -f docker/docker-compose.test.yml down -v 2>nul

:: 启动测试基础设施
echo 📦 启动 MongoDB 和 Redis...
docker-compose -f docker/docker-compose.test.yml up -d

:: 等待MongoDB就绪
echo ⏳ 等待 MongoDB 启动...
for /L %%i in (1,1,30) do (
    docker exec qingyu-mongodb-test mongo --eval "db.adminCommand('ping')" --quiet >nul 2>&1
    if !errorlevel! equ 0 (
        echo ✅ MongoDB 已就绪
        goto mongodb_ready
    )
    echo    等待 MongoDB... (%%i/30^)
    timeout /t 2 /nobreak >nul
)
echo ❌ MongoDB 启动失败
docker-compose -f docker/docker-compose.test.yml down -v
exit /b 1

:mongodb_ready

:: 等待Redis就绪
echo ⏳ 等待 Redis 启动...
for /L %%i in (1,1,15) do (
    docker exec qingyu-redis-test redis-cli ping >nul 2>&1
    if !errorlevel! equ 0 (
        echo ✅ Redis 已就绪
        goto redis_ready
    )
    echo    等待 Redis... (%%i/15^)
    timeout /t 1 /nobreak >nul
)
echo ❌ Redis 启动失败
docker-compose -f docker/docker-compose.test.yml down -v
exit /b 1

:redis_ready

:: 设置环境变量
if not defined MONGODB_URI set MONGODB_URI=mongodb://localhost:27017
set MONGODB_DATABASE=qingyu_test
set REDIS_ADDR=localhost:6379
set ENVIRONMENT=test

:: 运行测试
echo.
echo 🧪 运行测试...
echo ================================

set TEST_FAILED=0

:: 运行单元测试
echo.
echo 📝 运行单元测试...
go test -v -race -short -coverprofile=coverage_unit.txt -covermode=atomic ./service/... ./api/... ./middleware/...
if !errorlevel! equ 0 (
    echo ✅ 单元测试通过
) else (
    echo ❌ 单元测试失败
    set TEST_FAILED=1
)

:: 运行集成测试
echo.
echo 🔗 运行集成测试...
go test -v -race -timeout 10m ./test/integration/...
if !errorlevel! equ 0 (
    echo ✅ 集成测试通过
) else (
    echo ❌ 集成测试失败
    set TEST_FAILED=1
)

:: 运行API测试
echo.
echo 🌐 运行API测试...
go test -v -race -timeout 10m ./test/api/...
if !errorlevel! equ 0 (
    echo ✅ API测试通过
) else (
    echo ❌ API测试失败
    set TEST_FAILED=1
)

:: 清理
echo.
echo 🧹 清理测试环境...
docker-compose -f docker/docker-compose.test.yml down -v

:: 返回结果
if !TEST_FAILED! equ 0 (
    echo.
    echo 🎉 所有测试通过！
    exit /b 0
) else (
    echo.
    echo 💥 部分测试失败
    exit /b 1
)

