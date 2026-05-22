#!/bin/bash

# 测试运行脚本（使用Docker）
set -e

echo "🚀 启动测试环境..."

# 清理旧的测试环境
docker-compose -f docker/docker-compose.test.yml down -v 2>/dev/null || true

# 启动测试基础设施
echo "📦 启动 MongoDB 和 Redis..."
docker-compose -f docker/docker-compose.test.yml up -d

# 等待服务就绪
echo "⏳ 等待服务启动..."
for i in {1..30}; do
    if docker exec qingyu-mongodb-test mongo --eval "db.adminCommand('ping')" --quiet > /dev/null 2>&1; then
        echo "✅ MongoDB 已就绪"
        break
    fi
    echo "   等待 MongoDB... ($i/30)"
    sleep 2
done

for i in {1..15}; do
    if docker exec qingyu-redis-test redis-cli ping > /dev/null 2>&1; then
        echo "✅ Redis 已就绪"
        break
    fi
    echo "   等待 Redis... ($i/15)"
    sleep 1
done

# 设置环境变量
export MONGODB_URI="${MONGODB_URI:-mongodb://localhost:27017}"
export MONGODB_DATABASE="qingyu_test"
export REDIS_ADDR="localhost:6379"
export ENVIRONMENT="test"

# 运行测试
echo ""
echo "🧪 运行测试..."
echo "================================"

TEST_FAILED=0

# 运行单元测试
echo ""
echo "📝 运行单元测试..."
if go test -v -race -short -coverprofile=coverage_unit.txt -covermode=atomic \
    $(go list ./... | grep -v /test/); then
    echo "✅ 单元测试通过"
else
    echo "❌ 单元测试失败"
    TEST_FAILED=1
fi

# 运行集成测试
echo ""
echo "🔗 运行集成测试..."
if go test -v -race -timeout 10m ./test/integration/...; then
    echo "✅ 集成测试通过"
else
    echo "❌ 集成测试失败"
    TEST_FAILED=1
fi

# 运行API测试
echo ""
echo "🌐 运行API测试..."
if go test -v -race -timeout 10m ./test/api/...; then
    echo "✅ API测试通过"
else
    echo "❌ API测试失败"
    TEST_FAILED=1
fi

# 清理
echo ""
echo "🧹 清理测试环境..."
docker-compose -f docker/docker-compose.test.yml down -v

# 返回结果
if [ $TEST_FAILED -eq 0 ]; then
    echo ""
    echo "🎉 所有测试通过！"
    exit 0
else
    echo ""
    echo "💥 部分测试失败"
    exit 1
fi

