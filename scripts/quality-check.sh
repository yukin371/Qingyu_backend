#!/bin/bash
# Qingyu Backend 代码质量检查
# 集成所有质量检查脚本

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "======================================"
echo "  Qingyu Backend 代码质量检查"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

TOTAL_EXIT_CODE=0

# 1. 类型一致性检查
echo -e "${GREEN}[1/5]${NC} 类型一致性检查..."
if bash "$SCRIPT_DIR/check-type-consistency.sh"; then
    echo -e "${GREEN}✓ 类型一致性检查通过${NC}"
else
    echo -e "${RED}✗ 类型一致性检查失败${NC}"
    TOTAL_EXIT_CODE=1
fi
echo ""

# 2. Repository职责检查
echo -e "${GREEN}[2/5]${NC} Repository职责检查..."
if bash "$SCRIPT_DIR/check-repository-responsibility.sh"; then
    echo -e "${GREEN}✓ Repository职责检查通过${NC}"
else
    echo -e "${RED}✗ Repository职责检查失败${NC}"
    TOTAL_EXIT_CODE=1
fi
echo ""

# 3. 测试覆盖率检查
echo -e "${GREEN}[3/5]${NC} 测试覆盖率检查..."
if command -v go &> /dev/null; then
    echo "运行测试并收集覆盖率..."
    if go test -coverprofile=coverage.out ./... 2>/dev/null; then
        echo ""
        echo "覆盖率统计："
        go tool cover -func=coverage.out | grep total || true

        # 检查总覆盖率是否达到要求
        TOTAL_COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//' || echo "0")
        echo ""
        if (( $(echo "$TOTAL_COVERAGE >= 60" | bc -l 2>/dev/null || echo "0") )); then
            echo -e "${GREEN}✓ 测试覆盖率: ${TOTAL_COVERAGE}%${NC}"
        else
            echo -e "${YELLOW}⚠ 测试覆盖率: ${TOTAL_COVERAGE}% (建议 >= 60%)${NC}"
        fi
    else
        echo -e "${YELLOW}⚠ 部分测试失败，请检查${NC}"
        TOTAL_EXIT_CODE=1
    fi
else
    echo -e "${YELLOW}⚠ 未找到go命令，跳过测试检查${NC}"
fi
echo ""

# 4. 代码格式检查
echo -e "${GREEN}[4/5]${NC} 代码格式检查..."
if command -v gofmt &> /dev/null; then
    UNFORMATTED=$(gofmt -l . 2>/dev/null | grep -v vendor | head -10 || true)
    if [ -n "$UNFORMATTED" ]; then
        echo -e "${YELLOW}⚠ 以下文件需要格式化：${NC}"
        echo "$UNFORMATTED"
        echo "提示: 运行 'gofmt -w .' 进行格式化"
        TOTAL_EXIT_CODE=1
    else
        echo -e "${GREEN}✓ 代码格式检查通过${NC}"
    fi
else
    echo -e "${YELLOW}⚠ 未找到gofmt命令，跳过格式检查${NC}"
fi
echo ""

# 5. 静态分析
echo -e "${GREEN}[5/5]${NC} 静态分析..."
if command -v go &> /dev/null; then
    if go vet ./... 2>/dev/null; then
        echo -e "${GREEN}✓ 静态分析通过${NC}"
    else
        echo -e "${RED}✗ 静态分析发现问题${NC}"
        TOTAL_EXIT_CODE=1
    fi
else
    echo -e "${YELLOW}⚠ 未找到go命令，跳过静态分析${NC}"
fi
echo ""

# 总结
echo "======================================"
if [ $TOTAL_EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}✓ 所有检查通过！${NC}"
else
    echo -e "${RED}✗ 部分检查未通过，请修复后重试${NC}"
fi
echo "======================================"

exit $TOTAL_EXIT_CODE
