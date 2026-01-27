#!/bin/bash
# Block 3 阶段3验收检查清单

set -e

echo "=================================================="
echo "        Block 3 阶段3验收检查清单"
echo "=================================================="
echo ""

PASS_COUNT=0
FAIL_COUNT=0

# 检查项1: 缓存装饰器代码
echo "检查项1: 验证缓存装饰器代码"
if [ -f "repository/cache/cached_repository.go" ]; then
    LINES=$(wc -l < repository/cache/cached_repository.go)
    echo "✅ cached_repository.go 存在 ($LINES 行)"
    ((PASS_COUNT++))
else
    echo "❌ cached_repository.go 不存在"
    ((FAIL_COUNT++))
fi

# 检查项2: 缓存装饰器测试
echo ""
echo "检查项2: 运行缓存装饰器测试"
if go test ./repository/cache/... -v > /tmp/cache_test.log 2>&1; then
    TEST_COUNT=$(grep -c "PASS:" /tmp/cache_test.log || echo "0")
    echo "✅ 缓存装饰器测试通过 ($TEST_COUNT 个测试)"
    ((PASS_COUNT++))
else
    echo "❌ 缓存装饰器测试失败"
    cat /tmp/cache_test.log
    ((FAIL_COUNT++))
fi

# 检查项3: 缓存包装器
echo ""
echo "检查项3: 验证缓存包装器代码"
if [ -f "repository/mongodb/bookstore/cached_book_repository.go" ] && \
   [ -f "repository/mongodb/user/cached_user_repository.go" ]; then
    echo "✅ BookRepository 和 UserRepository 缓存包装器存在"
    ((PASS_COUNT++))
else
    echo "❌ 缓存包装器代码不完整"
    ((FAIL_COUNT++))
fi

# 检查项4: 缓存配置
echo ""
echo "检查项4: 验证缓存配置"
if grep -q "CacheConfig" config/cache.go 2>/dev/null; then
    echo "✅ 缓存配置已定义"
    ((PASS_COUNT++))
else
    echo "❌ 缓存配置不存在"
    ((FAIL_COUNT++))
fi

# 检查项5: 缓存预热
echo ""
echo "检查项5: 验证缓存预热机制"
if [ -f "pkg/cache/warmer.go" ]; then
    LINES=$(wc -l < pkg/cache/warmer.go)
    echo "✅ 缓存预热器存在 ($LINES 行)"
    ((PASS_COUNT++))
else
    echo "❌ 缓存预热器不存在"
    ((FAIL_COUNT++))
fi

# 检查项6: 运行缓存预热测试
echo ""
echo "检查项6: 运行缓存预热测试"
if go test ./pkg/cache/... -v > /tmp/warmer_test.log 2>&1; then
    TEST_COUNT=$(grep -c "PASS:" /tmp/warmer_test.log || echo "0")
    echo "✅ 缓存预热测试通过 ($TEST_COUNT 个测试)"
    ((PASS_COUNT++))
else
    echo "❌ 缓存预热测试失败"
    cat /tmp/warmer_test.log
    ((FAIL_COUNT++))
fi

# 总结
echo ""
echo "=================================================="
echo "        验收结果"
echo "=================================================="
echo "✅ 通过: $PASS_COUNT 项"
echo "❌ 失败: $FAIL_COUNT 项"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo "🎉 阶段3验收通过！"
    exit 0
else
    echo "⚠️  阶段3验收失败，请检查上述错误"
    exit 1
fi
