#!/bin/bash
# 青羽后端 - 快速验证脚本
# 用于验证项目编译和测试是否通过

echo "======================================"
echo "青羽后端 - 快速验证脚本"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 计数器
SUCCESS_COUNT=0
FAIL_COUNT=0

echo "步骤 1: 编译项目..."
if go build -o Qingyu_backend.exe; then
    echo -e "${GREEN}✅ 编译成功${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ 编译失败${NC}"
    ((FAIL_COUNT++))
    exit 1
fi

echo ""
echo "步骤 2: 运行Repository层测试..."
if go test ./test/repository/ -v; then
    echo -e "${GREEN}✅ Repository测试通过${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Repository测试失败${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "步骤 3: 运行Service层测试..."
if go test ./test/service/ -v; then
    echo -e "${GREEN}✅ Service测试通过${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${RED}❌ Service测试失败${NC}"
    ((FAIL_COUNT++))
fi

echo ""
echo "步骤 4: 运行书城系统测试..."
if go test ./test/ -run "Bookstore" -v; then
    echo -e "${GREEN}✅ 书城测试通过${NC}"
    ((SUCCESS_COUNT++))
else
    echo -e "${YELLOW}⚠️  书城测试警告（可能需要数据库连接）${NC}"
fi

echo ""
echo "======================================"
echo "验证完成！"
echo "======================================"
echo -e "成功: ${GREEN}${SUCCESS_COUNT}${NC} 项"
echo -e "失败: ${RED}${FAIL_COUNT}${NC} 项"
echo ""

if [ $FAIL_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 所有检查通过！项目状态健康！${NC}"
    exit 0
else
    echo -e "${RED}⚠️  有 ${FAIL_COUNT} 项检查失败，请修复后再提交代码${NC}"
    exit 1
fi

