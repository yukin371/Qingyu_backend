@echo off
echo ========================================
echo 青羽后端 - 快速测试数据初始化
echo ========================================
echo.

echo [1/3] 准备测试数据（书籍+章节）...
go run cmd/prepare_test_data/main.go
if errorlevel 1 (
    echo 错误：测试数据创建失败
    pause
    exit /b 1
)
echo.

echo [2/3] 创建 Banner 轮播图...
go run cmd/create_banners/main.go
if errorlevel 1 (
    echo 错误：Banner 创建失败
    pause
    exit /b 1
)
echo.

echo [3/3] 创建测试用户...
go run cmd/create_beta_users/main.go
if errorlevel 1 (
    echo 警告：测试用户创建失败（可能已存在）
)
echo.

echo ========================================
echo ✅ 测试环境初始化完成！
echo ========================================
echo.
echo 现在可以测试以下功能：
echo - 首页 Banner 轮播
echo - 书籍详情页
echo - 排行榜
echo - 用户登录（测试账号见下方）
echo.
echo 测试账号：
echo   管理员：admin / admin123
echo   作者：author1 / author123
echo   读者：reader1 / reader123
echo.
echo 前端访问：http://localhost:5173
echo 后端API：http://localhost:8080
echo.
pause


