#!/bin/bash
# 运行后端集成测试脚本
#
# 使用方法:
#   ./scripts/run_integration_tests.sh [test_pattern]
#
# 示例:
#   ./scripts/run_integration_tests.sh                    # 运行所有集成测试
#   ./scripts/run_integration_tests.sh TestCollectionScenario  # 运行特定测试
#
# 前提条件:
#   1. Docker服务已启动: docker-compose -f docker/docker-compose.integration.yml up -d
#   2. 等待MongoDB副本集初始化完成（约10秒）

set -e

# 设置环境变量
export MONGODB_URI="mongodb://localhost:27017/?replicaSet=rs0&directConnection=true"
export REDIS_HOST="localhost"
export REDIS_PORT="6379"

# 进入后端目录
cd "$(dirname "$0")/.."

# 检查Docker服务是否运行
echo "检查Docker服务状态..."
if ! docker ps | grep -q "qingyu-mongodb-primary"; then
    echo "❌ MongoDB容器未运行，请先启动:"
    echo "   docker-compose -f docker/docker-compose.integration.yml up -d"
    exit 1
fi

# 等待MongoDB就绪
echo "等待MongoDB副本集就绪..."
sleep 3
docker exec qingyu-mongodb-primary mongosh --eval "rs.status().ok" --quiet > /dev/null 2>&1 || {
    echo "❌ MongoDB副本集未就绪，请稍后重试"
    exit 1
}
echo "✓ MongoDB副本集已就绪"

# 运行测试
TEST_PATTERN="${1:-.}"
echo ""
echo "=========================================="
echo "运行集成测试: $TEST_PATTERN"
echo "=========================================="
echo ""

go test -v -count=1 -timeout 15m -run "$TEST_PATTERN" ./test/integration/...

echo ""
echo "=========================================="
echo "集成测试完成"
echo "=========================================="
