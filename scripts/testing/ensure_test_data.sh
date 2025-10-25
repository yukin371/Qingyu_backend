#!/bin/bash

# 测试数据准备脚本
# 确保测试环境有足够的测试数据

set -e

echo "========================================="
echo "  测试数据准备脚本"
echo "========================================="

# 配置
DB_NAME="qingyu_test"
MIN_BOOKS=10
MIN_CHAPTERS=100

# 检查MongoDB是否运行
echo ""
echo "1. 检查MongoDB连接..."
if ! mongosh --quiet --eval "db.version()" > /dev/null 2>&1; then
    echo "❌ MongoDB未运行或无法连接"
    exit 1
fi
echo "✓ MongoDB连接正常"

# 检查书籍数量
echo ""
echo "2. 检查书籍数据..."
BOOK_COUNT=$(mongosh $DB_NAME --quiet --eval "db.books.countDocuments({})")
echo "   当前书籍数量: $BOOK_COUNT"

if [ "$BOOK_COUNT" -lt "$MIN_BOOKS" ]; then
    echo "   ⚠ 书籍数量不足（需要至少 $MIN_BOOKS 本）"
    echo "   正在导入测试书籍..."
    go run cmd/migrate/main.go --seed books
    echo "   ✓ 书籍数据导入完成"
else
    echo "   ✓ 书籍数据充足"
fi

# 检查章节数量
echo ""
echo "3. 检查章节数据..."
CHAPTER_COUNT=$(mongosh $DB_NAME --quiet --eval "db.chapters.countDocuments({})")
echo "   当前章节数量: $CHAPTER_COUNT"

if [ "$CHAPTER_COUNT" -lt "$MIN_CHAPTERS" ]; then
    echo "   ⚠ 章节数量不足（需要至少 $MIN_CHAPTERS 个）"
    echo "   正在导入测试章节..."
    go run cmd/migrate/main.go --seed chapters
    echo "   ✓ 章节数据导入完成"
else
    echo "   ✓ 章节数据充足"
fi

# 检查测试用户
echo ""
echo "4. 检查测试用户..."
USER_COUNT=$(mongosh $DB_NAME --quiet --eval 'db.users.countDocuments({username: /^test_user/})')
echo "   当前测试用户数量: $USER_COUNT"

if [ "$USER_COUNT" -lt 5 ]; then
    echo "   ⚠ 测试用户不足"
    echo "   正在创建测试用户..."
    go run cmd/create_beta_users/main.go
    echo "   ✓ 测试用户创建完成"
else
    echo "   ✓ 测试用户充足"
fi

# 完成
echo ""
echo "========================================="
echo "  ✓ 测试数据准备完成"
echo "========================================="
echo ""
echo "数据统计:"
echo "  - 书籍: $BOOK_COUNT 本"
echo "  - 章节: $CHAPTER_COUNT 个"
echo "  - 测试用户: $USER_COUNT 个"
echo ""

