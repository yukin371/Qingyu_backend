@echo off
REM 完整测试数据准备脚本 (Windows)
REM 包括书籍、章节、用户和AI配额

echo =========================================
echo   测试数据准备脚本 (完整版)
echo =========================================
echo.

REM 配置
set DB_NAME=qingyu_test
set MONGO_URI=mongodb://localhost:27017

echo 1. 检查MongoDB连接...
mongosh %MONGO_URI% --quiet --eval "db.version()" >nul 2>&1
if errorlevel 1 (
    echo X MongoDB未运行或无法连接
    echo   请启动MongoDB服务
    exit /b 1
)
echo √ MongoDB连接正常
echo.

echo 2. 准备测试用户和AI配额...
echo   正在运行create_beta_users工具...
cd /d "%~dp0..\.."
go run cmd/create_beta_users/main.go
if errorlevel 1 (
    echo X 测试用户创建失败
    exit /b 1
)
echo √ 测试用户和AI配额准备完成
echo.

echo 3. 检查书籍数据...
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.books.countDocuments({})"') do set BOOK_COUNT=%%i
echo   当前书籍数量: %BOOK_COUNT%

if %BOOK_COUNT% LSS 10 (
    echo   ! 书籍数量不足，正在导入...
    go run cmd/migrate/main.go --seed books
    echo   √ 书籍数据导入完成
) else (
    echo   √ 书籍数据充足
)
echo.

echo 4. 检查章节数据...
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.chapters.countDocuments({})"') do set CHAPTER_COUNT=%%i
echo   当前章节数量: %CHAPTER_COUNT%

if %CHAPTER_COUNT% LSS 50 (
    echo   ! 章节数量不足，正在导入...
    go run cmd/migrate/main.go --seed chapters
    echo   √ 章节数据导入完成
) else (
    echo   √ 章节数据充足
)
echo.

echo 5. 激活测试用户AI配额...
echo   正在为测试用户激活AI配额...
mongosh %DB_NAME% --quiet --eval "db.ai_quotas.updateMany({user_id: {$in: [ObjectId('670f8b9a5e6d3c001f123456'), ObjectId('670f8b9a5e6d3c001f123457')]}}, {$set: {monthly_limit: 10000, daily_limit: 1000, used_this_month: 0, used_today: 0, status: 'active', updated_at: new Date()}})" >nul 2>&1
echo √ AI配额激活完成
echo.

echo 6. 检查分类数据...
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.categories.countDocuments({})"') do set CATEGORY_COUNT=%%i
echo   当前分类数量: %CATEGORY_COUNT%

if %CATEGORY_COUNT% LSS 5 (
    echo   ! 分类数据不足，正在创建...
    mongosh %DB_NAME% --quiet --eval "db.categories.insertMany([{name: '玄幻', slug: 'xuanhuan', description: '玄幻小说', parent_id: null, level: 1, sort_order: 1, is_active: true, created_at: new Date(), updated_at: new Date()}, {name: '都市', slug: 'dushi', description: '都市小说', parent_id: null, level: 1, sort_order: 2, is_active: true, created_at: new Date(), updated_at: new Date()}, {name: '仙侠', slug: 'xianxia', description: '仙侠小说', parent_id: null, level: 1, sort_order: 3, is_active: true, created_at: new Date(), updated_at: new Date()}, {name: '科幻', slug: 'kehuan', description: '科幻小说', parent_id: null, level: 1, sort_order: 4, is_active: true, created_at: new Date(), updated_at: new Date()}, {name: '历史', slug: 'lishi', description: '历史小说', parent_id: null, level: 1, sort_order: 5, is_active: true, created_at: new Date(), updated_at: new Date()}])" >nul 2>&1
    echo   √ 分类数据创建完成
) else (
    echo   √ 分类数据充足
)
echo.

echo =========================================
echo   √ 测试数据准备完成
echo =========================================
echo.

echo 数据统计:
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.books.countDocuments({})"') do echo   - 书籍: %%i 本
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.chapters.countDocuments({})"') do echo   - 章节: %%i 个
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.users.countDocuments({username: /^test_user|^vip_user/})"') do echo   - 测试用户: %%i 个
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.categories.countDocuments({})"') do echo   - 分类: %%i 个
for /f %%i in ('mongosh %DB_NAME% --quiet --eval "db.ai_quotas.countDocuments({status: 'active'})"') do echo   - 激活的AI配额: %%i 个
echo.

echo 可以开始运行测试了！
echo   go test ./test/integration -v -count=1
echo.
pause

