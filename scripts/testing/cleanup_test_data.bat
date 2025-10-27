@echo off
REM ========================================
REM 清理测试数据脚本
REM 用途：清空数据库中的测试数据，为集成测试做准备
REM ========================================

echo ========================================
echo 青羽后端 - 测试数据清理
echo ========================================
echo.

REM 检查 MongoDB 是否运行
echo [1/3] 检查 MongoDB 服务...
mongosh --eval "db.version()" >nul 2>&1
if errorlevel 1 (
    echo [错误] MongoDB 服务未运行，请先启动 MongoDB
    echo 提示：可以运行 docker-compose up -d mongodb 或启动本地 MongoDB 服务
    pause
    exit /b 1
)
echo [成功] MongoDB 服务正常运行
echo.

REM 询问确认
echo [警告] 此操作将清空以下数据：
echo   - books（书籍）
echo   - chapters（章节）
echo   - users（用户，保留系统用户）
echo   - ranking_items（榜单）
echo   - reading_progress（阅读进度）
echo   - annotations（书签和笔记）
echo   - user_collections（收藏）
echo.
set /p confirm="确认清理？(yes/no): "
if not "%confirm%"=="yes" (
    echo 操作已取消
    pause
    exit /b 0
)

echo.
echo [2/3] 连接到数据库并清理数据...

REM 使用 mongosh 执行清理命令
mongosh --quiet --eval "
use qingyu_test;

print('[清理] 删除书籍数据...');
db.books.deleteMany({});
print('  ✓ books 集合已清空');

print('[清理] 删除章节数据...');
db.chapters.deleteMany({});
print('  ✓ chapters 集合已清空');

print('[清理] 删除用户数据（保留系统用户）...');
db.users.deleteMany({ role: { $ne: 'system' } });
print('  ✓ users 集合已清空（保留系统用户）');

print('[清理] 删除榜单数据...');
db.ranking_items.deleteMany({});
print('  ✓ ranking_items 集合已清空');

print('[清理] 删除阅读进度...');
db.reading_progress.deleteMany({});
print('  ✓ reading_progress 集合已清空');

print('[清理] 删除书签和笔记...');
db.annotations.deleteMany({});
print('  ✓ annotations 集合已清空');

print('[清理] 删除收藏记录...');
db.user_collections.deleteMany({});
print('  ✓ user_collections 集合已清空');

print('[清理] 删除评论数据...');
db.comments.deleteMany({});
print('  ✓ comments 集合已清空');

print('[清理] 删除写作项目（测试项目）...');
db.projects.deleteMany({ title: /测试/ });
print('  ✓ projects 测试数据已清空');

print('');
print('[统计] 清理后数据统计：');
print('  - 书籍数量: ' + db.books.countDocuments());
print('  - 章节数量: ' + db.chapters.countDocuments());
print('  - 用户数量: ' + db.users.countDocuments());
print('  - 榜单数量: ' + db.ranking_items.countDocuments());
print('  - 阅读进度: ' + db.reading_progress.countDocuments());
print('  - 书签笔记: ' + db.annotations.countDocuments());
"

if errorlevel 1 (
    echo.
    echo [错误] 数据清理失败
    pause
    exit /b 1
)

echo.
echo [3/3] 清理完成
echo.
echo ========================================
echo ✓ 测试数据清理成功
echo ========================================
echo.
echo 提示：现在可以运行 import_test_users.go 导入测试用户
echo.
pause


