#!/bin/bash
# 性能对比测试脚本

set -e

# 配置
BASE_URL=${BASE_URL:-"http://localhost:8080"}
OUTPUT_DIR=${OUTPUT_DIR:-"test_results"}
BOOK_ID=${BOOK_ID:-"507f1f77bcf86cd799439011"}

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

# 清空Redis缓存
clear_cache() {
    log_info "清空Redis缓存..."
    redis-cli FLUSHDB 2>/dev/null || log_warn "Redis未启动或FLUSHDB失败"
    sleep 1
}

# 执行基准测试
run_benchmark() {
    local cache_enabled=$1
    local output_file="$OUTPUT_DIR/result_cache_${cache_enabled}.json"

    log_info "执行测试（缓存: $cache_enabled）..."

    cd Qingyu_backend-block3-optimization

    go run cmd/benchmark/main.go \
        --base-url="$BASE_URL" \
        --requests=1000 \
        --concurrent=50 \
        --name="Cache_${cache_enabled}" \
        --output="$output_file" || log_warn "测试执行失败"

    log_info "测试完成，结果保存到: $output_file"
}

# 生成对比报告
generate_comparison_report() {
    log_info "生成性能对比报告..."

    python3 scripts/generate_comparison.py \
        --with-cache "$OUTPUT_DIR/result_cache_true.json" \
        --without-cache "$OUTPUT_DIR/result_cache_false.json" \
        --output "$OUTPUT_DIR/comparison_report.md" 2>/dev/null || log_warn "报告生成失败"

    if [ -f "$OUTPUT_DIR/comparison_report.md" ]; then
        log_info "对比报告生成完成: $OUTPUT_DIR/comparison_report.md"
    fi
}

# 主流程
main() {
    local mode=${1:-"compare"}

    case $mode in
        "with-cache")
            clear_cache
            run_benchmark true
            ;;
        "without-cache")
            clear_cache
            run_benchmark false
            ;;
        "compare")
            log_info "开始性能对比测试..."

            # 测试1: 无缓存
            clear_cache
            run_benchmark false

            echo ""

            # 测试2: 有缓存
            clear_cache
            run_benchmark true

            echo ""

            # 生成对比报告
            generate_comparison_report

            log_info "性能对比测试完成！"
            ;;
        *)
            echo "用法: $0 [with-cache|without-cache|compare]"
            exit 1
            ;;
    esac
}

main "$@"
