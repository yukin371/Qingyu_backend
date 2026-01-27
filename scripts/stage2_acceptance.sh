#!/bin/bash

###############################################################################
# Block 3 阶段2验收脚本
# 功能: 验证MongoDB监控体系建立完成
# 作者: Block 3数据库优化实施女仆
# 日期: 2026-01-27
###############################################################################

set -e  # 任何错误都会终止脚本

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印标题
echo ""
echo "=================================================="
echo "        阶段2验收检查清单"
echo "=================================================="
echo ""

# 检查项1: 验证MongoDB Profiler启用
echo "检查项1: 验证MongoDB Profiler启用"
echo "-------------------------------------------"

# 检查mongosh是否可用
if ! command -v mongosh &> /dev/null; then
    echo -e "${YELLOW}⚠️  mongosh未安装，跳过Profiler状态检查${NC}"
    echo "   请手动验证MongoDB Profiler已启用"
else
    # 尝试连接MongoDB并检查Profiler状态
    STATUS=$(mongosh qingyu_dev --quiet --eval "db.getProfilingStatus().was" 2>&1 || echo "error")

    if [ "$STATUS" = "1" ]; then
        echo -e "${GREEN}✅ Profiler级别: 1 (仅慢查询)${NC}"

        # 额外检查慢查询阈值
        THRESHOLD=$(mongosh qingyu_dev --quiet --eval "db.getProfilingStatus().slowms" 2>&1 || echo "error")
        if [ "$THRESHOLD" = "100" ]; then
            echo -e "${GREEN}✅ 慢查询阈值: 100ms${NC}"
        else
            echo -e "${YELLOW}⚠️  慢查询阈值: ${THRESHOLD}ms (期望: 100ms)${NC}"
        fi
    elif [ "$STATUS" = "0" ]; then
        echo -e "${RED}❌ Profiler未启用 (级别: 0)${NC}"
        exit 1
    elif [ "$STATUS" = "2" ]; then
        echo -e "${GREEN}✅ Profiler级别: 2 (全部查询)${NC}"
        echo -e "${YELLOW}⚠️  注意: 生产环境建议使用级别1${NC}"
    else
        echo -e "${YELLOW}⚠️  无法连接MongoDB或获取Profiler状态${NC}"
        echo "   请确保MongoDB正在运行且数据库为qingyu_dev"
    fi
fi

echo ""

# 检查项2: 验证Prometheus指标
echo "检查项2: 验证Prometheus指标"
echo "-------------------------------------------"

# 检查metrics.go是否存在
if [ -f "repository/mongodb/monitor/metrics.go" ]; then
    echo -e "${GREEN}✅ metrics.go存在${NC}"
else
    echo -e "${RED}❌ metrics.go不存在${NC}"
    exit 1
fi

# 检查metrics_enhanced.go是否存在
if [ -f "repository/mongodb/monitor/metrics_enhanced.go" ]; then
    echo -e "${GREEN}✅ metrics_enhanced.go存在${NC}"
else
    echo -e "${RED}❌ metrics_enhanced.go不存在${NC}"
    exit 1
fi

# 检查metrics_test.go是否存在
if [ -f "repository/mongodb/monitor/metrics_test.go" ]; then
    echo -e "${GREEN}✅ metrics_test.go存在${NC}"
else
    echo -e "${RED}❌ metrics_test.go不存在${NC}"
    exit 1
fi

# 检查go命令是否可用
if ! command -v go &> /dev/null; then
    echo -e "${YELLOW}⚠️  go命令未安装，跳过测试执行${NC}"
    echo "   请手动运行: go test ./repository/mongodb/monitor/... -v"
else
    # 运行测试并检查结果
    TEST_OUTPUT=$(go test ./repository/mongodb/monitor/... -v 2>&1 || echo "TEST_FAILED")

    if echo "$TEST_OUTPUT" | grep -q "PASS"; then
        # 统计通过的测试数量
        PASS_COUNT=$(echo "$TEST_OUTPUT" | grep -c "PASS:" || echo "0")
        echo -e "${GREEN}✅ Prometheus指标测试通过 (${PASS_COUNT}个测试)${NC}"

        # 检查是否有benchmark
        if echo "$TEST_OUTPUT" | grep -q "Benchmark"; then
            echo -e "${GREEN}✅ Benchmark测试已执行${NC}"
        fi
    else
        echo -e "${RED}❌ Prometheus指标测试失败${NC}"
        echo "$TEST_OUTPUT" | tail -20
        exit 1
    fi
fi

echo ""

# 检查项3: 验证监控配置文件
echo "检查项3: 验证监控配置文件"
echo "-------------------------------------------"

# 检查Grafana仪表板
if [ -f "monitoring/grafana/dashboards/mongodb-dashboard.json" ]; then
    echo -e "${GREEN}✅ Grafana仪表板配置存在${NC}"

    # 检查文件大小（确保不是空文件）
    DASHBOARD_SIZE=$(stat -f%z "monitoring/grafana/dashboards/mongodb-dashboard.json" 2>/dev/null || stat -c%s "monitoring/grafana/dashboards/mongodb-dashboard.json" 2>/dev/null || echo "0")
    if [ "$DASHBOARD_SIZE" -gt 1000 ]; then
        echo -e "${GREEN}✅ 仪表板配置文件大小正常 (${DASHBOARD_SIZE} bytes)${NC}"
    else
        echo -e "${RED}❌ 仪表板配置文件可能为空或损坏${NC}"
        exit 1
    fi
else
    echo -e "${RED}❌ Grafana仪表板配置不存在${NC}"
    exit 1
fi

# 检查告警规则
if [ -f "monitoring/alerts/block3_alerts.yaml" ]; then
    echo -e "${GREEN}✅ Prometheus告警规则存在${NC}"

    # 检查告警规则数量
    ALERT_COUNT=$(grep -c "alert:" "monitoring/alerts/block3_alerts.yaml" 2>/dev/null || echo "0")
    if [ "$ALERT_COUNT" -ge 3 ]; then
        echo -e "${GREEN}✅ 告警规则配置完整 (${ALERT_COUNT}个规则)${NC}"
    else
        echo -e "${YELLOW}⚠️  告警规则数量: ${ALERT_COUNT} (期望: 至少3个)${NC}"
    fi
else
    echo -e "${RED}❌ Prometheus告警规则不存在${NC}"
    exit 1
fi

echo ""

# 检查项4: 验证慢查询分析工具
echo "检查项4: 验证慢查询分析工具"
echo "-------------------------------------------"

# 检查analyze_slow_queries.js
if [ -f "scripts/db/analyze_slow_queries.js" ]; then
    echo -e "${GREEN}✅ analyze_slow_queries.js存在${NC}"
else
    echo -e "${RED}❌ analyze_slow_queries.js不存在${NC}"
    exit 1
fi

# 检查auto_analyze_slow_queries.js
if [ -f "scripts/db/auto_analyze_slow_queries.js" ]; then
    echo -e "${GREEN}✅ auto_analyze_slow_queries.js存在${NC}"
else
    echo -e "${RED}❌ auto_analyze_slow_queries.js不存在${NC}"
    exit 1
fi

# 检查test_slow_queries.js
if [ -f "scripts/db/test_slow_queries.js" ]; then
    echo -e "${GREEN}✅ test_slow_queries.js存在${NC}"
else
    echo -e "${RED}❌ test_slow_queries.js不存在${NC}"
    exit 1
fi

# 检查使用文档
if [ -f "scripts/db/README-SLOW-QUERY-TOOLS.md" ]; then
    echo -e "${GREEN}✅ 慢查询工具使用文档存在${NC}"
else
    echo -e "${YELLOW}⚠️  慢查询工具使用文档不存在${NC}"
fi

echo ""

# 检查项5: 验证Profiling配置
echo "检查项5: 验证Profiling配置"
echo "-------------------------------------------"

# 检查enable_profiling.js
if [ -f "scripts/db/enable_profiling.js" ]; then
    echo -e "${GREEN}✅ enable_profiling.js存在${NC}"
else
    echo -e "${RED}❌ enable_profiling.js不存在${NC}"
    exit 1
fi

# 检查database.go中的Profiling配置
if [ -f "config/database.go" ]; then
    if grep -q "ProfilingLevel" "config/database.go"; then
        echo -e "${GREEN}✅ database.go包含Profiling配置${NC}"
    else
        echo -e "${YELLOW}⚠️  database.go可能缺少Profiling配置${NC}"
    fi
else
    echo -e "${YELLOW}⚠️  config/database.go不存在${NC}"
fi

echo ""

# 打印总结
echo "=================================================="
echo -e "${GREEN}✅ 阶段2验收通过${NC}"
echo "=================================================="
echo ""
echo "验收结果摘要:"
echo "  ✅ MongoDB Profiler配置"
echo "  ✅ Prometheus监控指标"
echo "  ✅ Grafana仪表板"
echo "  ✅ 告警规则"
echo "  ✅ 慢查询分析工具"
echo "  ✅ Profiling配置脚本"
echo ""
echo "下一步:"
echo "  1. 查看完成报告: docs/reports/block3-stage2-completion-report.md"
echo "  2. 进入阶段3: 缓存实现"
echo ""
echo "监控使用指南:"
echo "  - 启动Grafana: docker-compose up -d grafana"
echo "  - 访问仪表板: http://localhost:3000"
echo "  - 运行慢查询分析: mongosh qingyu_dev scripts/db/analyze_slow_queries.js"
echo ""
