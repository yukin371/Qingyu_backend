#!/bin/bash

# 青羽后端测试脚本
# 功能：运行所有测试并生成报告

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
COVERAGE_DIR="coverage"
COVERAGE_FILE="${COVERAGE_DIR}/coverage.txt"
COVERAGE_HTML="${COVERAGE_DIR}/coverage.html"
TEST_TIMEOUT="10m"

echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${BLUE}   青羽后端 - 自动化测试套件${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""

# 创建覆盖率目录
mkdir -p ${COVERAGE_DIR}

# 1. 检查依赖
echo -e "${YELLOW}[1/6]${NC} 检查依赖..."
if ! command -v go &> /dev/null; then
    echo -e "${RED}✗ Go未安装${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Go版本: $(go version)${NC}"

# 2. 检查服务
echo -e "${YELLOW}[2/6]${NC} 检查服务状态..."

# 检查MongoDB
if ! mongosh --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    echo -e "${RED}✗ MongoDB未运行${NC}"
    echo -e "${YELLOW}提示: 运行 docker-compose up -d mongodb${NC}"
    exit 1
fi
echo -e "${GREEN}✓ MongoDB运行中${NC}"

# 检查Redis
if ! redis-cli ping > /dev/null 2>&1; then
    echo -e "${RED}✗ Redis未运行${NC}"
    echo -e "${YELLOW}提示: 运行 docker-compose up -d redis${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Redis运行中${NC}"

# 3. 代码格式检查
echo -e "${YELLOW}[3/6]${NC} 代码格式检查..."
UNFORMATTED=$(gofmt -l .)
if [ -n "$UNFORMATTED" ]; then
    echo -e "${RED}✗ 以下文件需要格式化:${NC}"
    echo "$UNFORMATTED"
    echo -e "${YELLOW}运行: gofmt -s -w .${NC}"
    exit 1
fi
echo -e "${GREEN}✓ 代码格式正确${NC}"

# 4. Lint检查
echo -e "${YELLOW}[4/6]${NC} 运行Lint检查..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run --timeout=${TEST_TIMEOUT} || {
        echo -e "${RED}✗ Lint检查失败${NC}"
        exit 1
    }
    echo -e "${GREEN}✓ Lint检查通过${NC}"
else
    echo -e "${YELLOW}⚠ golangci-lint未安装，跳过${NC}"
fi

# 5. 运行单元测试
echo -e "${YELLOW}[5/6]${NC} 运行单元测试..."
echo ""

export CONFIG_PATH="config/config.test.yaml"

go test \
    -v \
    -race \
    -timeout=${TEST_TIMEOUT} \
    -coverprofile=${COVERAGE_FILE} \
    -covermode=atomic \
    ./... || {
        echo -e "${RED}✗ 测试失败${NC}"
        exit 1
    }

echo ""
echo -e "${GREEN}✓ 单元测试通过${NC}"

# 6. 生成覆盖率报告
echo -e "${YELLOW}[6/6]${NC} 生成覆盖率报告..."

# 计算覆盖率
COVERAGE=$(go tool cover -func=${COVERAGE_FILE} | grep total | awk '{print $3}')
echo -e "${BLUE}总体覆盖率: ${COVERAGE}${NC}"

# 生成HTML报告
go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_HTML}
echo -e "${GREEN}✓ HTML报告已生成: ${COVERAGE_HTML}${NC}"

# 按包统计覆盖率
echo ""
echo -e "${BLUE}各包覆盖率统计:${NC}"
go tool cover -func=${COVERAGE_FILE} | grep -v "total" | sort -k3 -n -r | head -20

# 检查覆盖率阈值
COVERAGE_NUM=$(echo $COVERAGE | sed 's/%//')
THRESHOLD=80

echo ""
if (( $(echo "$COVERAGE_NUM < $THRESHOLD" | bc -l) )); then
    echo -e "${YELLOW}⚠ 覆盖率 ${COVERAGE} 低于阈值 ${THRESHOLD}%${NC}"
else
    echo -e "${GREEN}✓ 覆盖率 ${COVERAGE} 达标 (>=${THRESHOLD}%)${NC}"
fi

# 运行集成测试（可选）
if [ "$RUN_INTEGRATION" = "true" ]; then
    echo ""
    echo -e "${YELLOW}[额外]${NC} 运行集成测试..."
    go test -v -tags=integration ./test/integration/... || {
        echo -e "${RED}✗ 集成测试失败${NC}"
        exit 1
    }
    echo -e "${GREEN}✓ 集成测试通过${NC}"
fi

# 运行性能测试（可选）
if [ "$RUN_BENCHMARK" = "true" ]; then
    echo ""
    echo -e "${YELLOW}[额外]${NC} 运行性能测试..."
    go test -bench=. -benchmem -run=^$ ./... > ${COVERAGE_DIR}/benchmark.txt
    echo -e "${GREEN}✓ 性能测试完成: ${COVERAGE_DIR}/benchmark.txt${NC}"
fi

# 总结
echo ""
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo -e "${GREEN}✓ 所有测试通过！${NC}"
echo -e "${BLUE}════════════════════════════════════════${NC}"
echo ""
echo -e "${BLUE}报告位置:${NC}"
echo -e "  - 覆盖率报告: ${COVERAGE_HTML}"
echo -e "  - 覆盖率数据: ${COVERAGE_FILE}"
if [ "$RUN_BENCHMARK" = "true" ]; then
    echo -e "  - 性能测试: ${COVERAGE_DIR}/benchmark.txt"
fi
echo ""
echo -e "${YELLOW}提示: 在浏览器中打开 ${COVERAGE_HTML} 查看详细报告${NC}"

