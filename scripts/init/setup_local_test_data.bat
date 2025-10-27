@echo off
REM ================================================
REM 青羽后端 - 本地测试数据初始化脚本
REM ================================================
echo.
echo ====================================
echo 青羽写作系统 - 本地测试数据初始化
echo ====================================
echo.

REM 检查数据文件是否存在
if not exist "data\novels_100.json" (
    echo [错误] 数据文件 data\novels_100.json 不存在
    echo 请先运行 Python 脚本生成数据：
    echo python scripts\import_novels.py --max-novels 100 --output data\novels_100.json
    pause
    exit /b 1
)

echo [1/3] 检查 MongoDB 连接...
echo.

REM 先测试数据库连接（使用 status 命令）
go run cmd/migrate/main.go -command=status -config=.
if errorlevel 1 (
    echo.
    echo [错误] MongoDB 连接失败
    echo 请确保：
    echo 1. MongoDB 服务已启动
    echo 2. config/config.local.yaml 中的数据库配置正确
    pause
    exit /b 1
)

echo.
echo [2/3] 导入小说数据（100本）...
echo.

REM 导入小说数据
go run cmd/migrate/main.go -command=import-novels -file=data/novels_100.json -config=.
if errorlevel 1 (
    echo.
    echo [错误] 小说数据导入失败
    pause
    exit /b 1
)

echo.
echo [3/3] 创建内测用户...
echo.

REM 创建内测用户
go run cmd/create_beta_users/main.go
if errorlevel 1 (
    echo.
    echo [错误] 内测用户创建失败
    pause
    exit /b 1
)

echo.
echo ====================================
echo 测试数据初始化完成！
echo ====================================
echo.
echo 你现在可以：
echo 1. 启动服务器：go run cmd/server/main.go
echo 2. 使用以下测试账号登录：
echo.
echo    管理员账号：
echo      用户名: admin
echo      密码: Admin@123456
echo.
echo    VIP作家账号：
echo      用户名: vip_writer01
echo      密码: Vip@123456
echo.
echo    普通作家：
echo      用户名: writer_xuanhuan
echo      密码: Writer@123456
echo.
echo    普通读者：
echo      用户名: reader01
echo      密码: Reader@123456
echo.
echo 详细账号列表请查看上方输出
echo ====================================
echo.
pause

