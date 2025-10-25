@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo =========================================
echo   测试数据准备脚本
echo =========================================

set DB_NAME=qingyu_test
set MIN_BOOKS=10
set MIN_CHAPTERS=100

:: 检查MongoDB
echo.
echo 1. 检查MongoDB连接...
mongosh --quiet --eval "db.version()" >nul 2>&1
if errorlevel 1 (
    echo ❌ MongoDB未运行或无法连接
    exit /b 1
)
echo ✓ MongoDB连接正常

:: 检查书籍数量
echo.
echo 2. 检查书籍数据...
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.books.countDocuments({})"') do set BOOK_COUNT=%%i
echo    当前书籍数量: %BOOK_COUNT%

if %BOOK_COUNT% LSS %MIN_BOOKS% (
    echo    ⚠ 书籍数量不足（需要至少 %MIN_BOOKS% 本）
    echo    正在导入测试书籍...
    go run cmd/migrate/main.go --seed books
    echo    ✓ 书籍数据导入完成
) else (
    echo    ✓ 书籍数据充足
)

:: 检查章节数量
echo.
echo 3. 检查章节数据...
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.chapters.countDocuments({})"') do set CHAPTER_COUNT=%%i
echo    当前章节数量: %CHAPTER_COUNT%

if %CHAPTER_COUNT% LSS %MIN_CHAPTERS% (
    echo    ⚠ 章节数量不足（需要至少 %MIN_CHAPTERS% 个）
    echo    正在导入测试章节...
    go run cmd/migrate/main.go --seed chapters
    echo    ✓ 章节数据导入完成
) else (
    echo    ✓ 章节数据充足
)

:: 检查测试用户
echo.
echo 4. 检查测试用户...
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.users.countDocuments({username: /^test_user/})"') do set USER_COUNT=%%i
echo    当前测试用户数量: %USER_COUNT%

if %USER_COUNT% LSS 5 (
    echo    ⚠ 测试用户不足
    echo    正在创建测试用户...
    go run cmd/create_beta_users/main.go
    echo    ✓ 测试用户创建完成
) else (
    echo    ✓ 测试用户充足
)

:: 完成
echo.
echo =========================================
echo   ✓ 测试数据准备完成
echo =========================================
echo.
echo 数据统计:
echo   - 书籍: %BOOK_COUNT% 本
echo   - 章节: %CHAPTER_COUNT% 个
echo   - 测试用户: %USER_COUNT% 个
echo.

endlocal

