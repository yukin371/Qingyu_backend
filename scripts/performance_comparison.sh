#!/bin/bash
# 性能对比测试脚本

# 获取脚本所在目录的绝对路径
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 验证项目目录存在
if [ ! -d "$PROJECT_DIR" ]; then
    echo "[ERROR] 项目目录不存在: $PROJECT_DIR" >&2
    exit 1
fi

# 切换到项目根目录
cd "$PROJECT_DIR" || exit 1

# 配置
BASE_URL=${BASE_URL:-"http://localhost:8080"}
OUTPUT_DIR=${OUTPUT_DIR:-"test_results"}
BOOK_ID=${BOOK_ID:-"507f1f77bcf86cd799439011"}
PROMETHEUS_URL=${PROMETHEUS_URL:-"http://localhost:9090"}

# 错误日志函数
log_error() {
    echo -e "\033[0;31m[ERROR]\033[0m $1" >&2
}

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
    if ! redis-cli FLUSHDB 2>/dev/null; then
        log_warn "Redis未启动或FLUSHDB失败，继续测试..."
    fi
    sleep 1
}

# 执行基准测试
run_benchmark() {
    local cache_enabled=$1
    local output_file="$OUTPUT_DIR/result_cache_${cache_enabled}.json"

    log_info "执行测试（缓存: $cache_enabled）..."

    # 验证benchmark程序存在
    if [ ! -f "cmd/benchmark/main.go" ]; then
        log_error "benchmark程序不存在: cmd/benchmark/main.go"
        return 1
    fi

    if ! go run cmd/benchmark/main.go \
        --base-url="$BASE_URL" \
        --requests=1000 \
        --concurrent=50 \
        --name="Cache_${cache_enabled}" \
        --output="$output_file"; then
        log_error "测试执行失败（缓存: $cache_enabled）"
        return 1
    fi

    # 验证输出文件已生成
    if [ ! -f "$output_file" ]; then
        log_error "测试结果文件未生成: $output_file"
        return 1
    fi

    log_info "测试完成，结果保存到: $output_file"
    return 0
}

# 生成对比报告
generate_comparison_report() {
    local cache_true_file="$OUTPUT_DIR/result_cache_true.json"
    local cache_false_file="$OUTPUT_DIR/result_cache_false.json"
    local report_file="$OUTPUT_DIR/comparison_report.md"

    log_info "生成性能对比报告..."

    # 验证输入文件存在
    if [ ! -f "$cache_true_file" ]; then
        log_error "有缓存测试结果不存在: $cache_true_file"
        return 1
    fi

    if [ ! -f "$cache_false_file" ]; then
        log_error "无缓存测试结果不存在: $cache_false_file"
        return 1
    fi

    # 验证Python脚本存在
    if [ ! -f "scripts/generate_comparison.py" ]; then
        log_error "报告生成脚本不存在: scripts/generate_comparison.py"
        return 1
    fi

    if ! python3 scripts/generate_comparison.py \
        --with-cache "$cache_true_file" \
        --without-cache "$cache_false_file" \
        --output "$report_file"; then
        log_error "报告生成失败"
        return 1
    fi

    if [ -f "$report_file" ]; then
        log_info "对比报告生成完成: $report_file"
    else
        log_error "报告文件未生成: $report_file"
        return 1
    fi

    return 0
}

# 主流程
main() {
    local mode=${1:-"compare"}
    local exit_code=0

    case $mode in
        "with-cache")
            clear_cache
            if ! run_benchmark true; then
                log_error "有缓存测试失败"
                exit 1
            fi
            ;;
        "without-cache")
            clear_cache
            if ! run_benchmark false; then
                log_error "无缓存测试失败"
                exit 1
            fi
            ;;
        "compare")
            log_info "开始性能对比测试..."
            log_info "收集Prometheus指标从: $PROMETHEUS_URL"

            # 测试1: 无缓存
            clear_cache
            if ! run_benchmark false; then
                log_error "无缓存测试失败"
                exit 1
            fi

            echo ""

            # 测试2: 有缓存
            clear_cache
            if ! run_benchmark true; then
                log_error "有缓存测试失败"
                exit 1
            fi

            echo ""

            # 生成对比报告
            if ! generate_comparison_report; then
                log_error "报告生成失败"
                exit 1
            fi

            log_info "性能对比测试完成！"
            ;;
        *)
            log_error "用法: $0 [with-cache|without-cache|compare]"
            exit 1
            ;;
    esac
}

main "$@"
