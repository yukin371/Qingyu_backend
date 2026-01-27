#!/bin/bash
# 性能基线测试脚本

set -e

BASELINE_DIR="test_results/baselines"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BASELINE_FILE="$BASELINE_DIR/baseline_$TIMESTAMP.json"

# 创建基线目录
mkdir -p "$BASELINE_DIR"

echo "=== 性能基线测试 - $TIMESTAMP ==="

# 检查wrk是否安装（降级处理）
if ! command -v wrk &> /dev/null; then
    echo "⚠️  wrk未安装，跳过性能基线测试"
    echo ""
    echo "如需建立性能基线，请："
    echo "  1. 安装wrk: https://github.com/wg/wrk"
    echo "  2. 启动后端服务"
    echo "  3. 运行: ./scripts/performance_baseline.sh"
    echo ""
    echo "创建空基线文件用于后续参考..."

    # 创建空基线文件
    cat > "$BASELINE_FILE" << EOF
{
  "timestamp": "$TIMESTAMP",
  "note": "wrk未安装，性能基线未建立",
  "tests": {}
}
EOF

    echo "✓ 空基线文件已创建: $BASELINE_FILE"
    echo "   请安装wrk后重新运行测试"
    exit 0
fi

# 启动后端服务（如果未运行）
echo "检查服务状态..."
if ! curl -s http://localhost:8080/health > /dev/null; then
    echo "❌ 后端服务未运行，请先启动服务"
    exit 1
fi

echo "✓ 后端服务运行正常"

# 测试配置
HOST="http://localhost:8080"
DURATION="30s"
THREADS=4
CONNECTIONS=100

# 测试端点列表
declare -A ENDPOINTS=(
    [health]="/health"
    [user_login]="/api/v1/user/auth/login"
    [bookstore_home]="/api/v1/bookstore/homepage"
)

# JSON结果容器
echo "{" > "$BASELINE_FILE"
echo "  \"timestamp\": \"$TIMESTAMP\"," >> "$BASELINE_FILE"
echo "  \"tests\": {" >> "$BASELINE_FILE"

# 遍历测试端点
first=true
for endpoint in "${!ENDPOINTS[@]}"; do
    path="${ENDPOINTS[$endpoint]}"
    echo ""
    echo "测试端点: $endpoint ($path)"

    # 执行wrk测试
    result=$(wrk -t$THREADS -c$CONNECTIONS -D$DURATION "$HOST$path" 2>&1)

    # 提取关键指标
    rps=$(echo "$result" | grep "Requests/sec" | awk '{print $2}')
    latency_avg=$(echo "$result" | grep "Latency" | awk '{print $2}')
    latency_stdev=$(echo "$result" | grep "Latency" | awk '{print $3}')
    latency_p95=$(echo "$result" | grep "95%" | awk '{print $2}')

    echo "  RPS: $rps"
    echo "  平均延迟: $latency_avg"
    echo "  P95延迟: $latency_p95"

    # 写入JSON
    if [ "$first" = true ]; then
        first=false
    else
        echo "," >> "$BASELINE_FILE"
    fi

    cat >> "$BASELINE_FILE" << EOF
    "$endpoint": {
      "path": "$path",
      "rps": $rps,
      "latency_avg": "$latency_avg",
      "latency_stdev": "$latency_stdev",
      "latency_p95": "$latency_p95"
    }
EOF
done

echo "" >> "$BASELINE_FILE"
echo "  }" >> "$BASELINE_FILE"
echo "}" >> "$BASELINE_FILE"

echo ""
echo "=== 基线测试完成 ==="
echo "结果已保存到: $BASELINE_FILE"

# 输出摘要
echo ""
echo "=== 性能基线摘要 ==="
cat "$BASELINE_FILE" | jq '.tests'
