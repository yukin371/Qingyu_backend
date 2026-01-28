#!/bin/bash
# Prometheus指标采集脚本

set -e

PROMETHEUS_URL=${PROMETHEUS_URL:-"http://localhost:9090"}
OUTPUT_FILE=${OUTPUT_FILE:-"metrics.log"}
INTERVAL=${INTERVAL:-10} # 采集间隔（秒）

log_info() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

# 查询Prometheus指标
query_metric() {
    local metric_name=$1
    local query=$2

    curl -s "${PROMETHEUS_URL}/api/v1/query?query=${query}" \
        | jq -r '.data.result[0].value[1]' \
        >> "$OUTPUT_FILE"
}

# 采集所有指标
collect_all_metrics() {
    log_info "开始采集Prometheus指标..."

    while true; do
        echo "=== $(date '+%Y-%m-%d %H:%M:%S') ===" >> "$OUTPUT_FILE"

        # 缓存命中次数
        query_metric "cache_hits_total" "sum(cache_hits_total)" >> "$OUTPUT_FILE"
        echo "cache_hits_total" >> "$OUTPUT_FILE"

        # 缓存未命中次数
        query_metric "cache_misses_total" "sum(cache_misses_total)" >> "$OUTPUT_FILE"
        echo "cache_misses_total" >> "$OUTPUT_FILE"

        sleep "$INTERVAL"
    done
}

# 主流程
main() {
    log_info "Prometheus指标采集器启动"
    log_info "采集间隔: ${INTERVAL}秒"
    log_info "输出文件: $OUTPUT_FILE"

    collect_all_metrics
}

main
