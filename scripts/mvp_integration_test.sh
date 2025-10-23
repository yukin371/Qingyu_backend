#!/bin/bash

# MVP集成测试脚本
# 用途：验证MVP新增功能在实际环境中的运行情况
# 日期：2025-10-23

set -e

echo "======================================"
echo "MVP集成测试脚本"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果记录函数
test_result() {
    local test_name=$1
    local result=$2

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    if [ "$result" = "PASS" ]; then
        echo -e "${GREEN}✓${NC} $test_name"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗${NC} $test_name"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# 1. 环境检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1. 环境检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查MongoDB
if docker ps | grep -q mongodb; then
    test_result "MongoDB运行状态" "PASS"
else
    test_result "MongoDB运行状态" "FAIL"
    echo -e "${RED}错误: MongoDB未运行，请执行 docker-compose up -d${NC}"
fi

# 检查Redis
if docker ps | grep -q redis; then
    test_result "Redis运行状态" "PASS"
else
    test_result "Redis运行状态" "FAIL"
    echo -e "${RED}错误: Redis未运行，请执行 docker-compose up -d${NC}"
fi

# 检查Go环境
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    test_result "Go环境 ($GO_VERSION)" "PASS"
else
    test_result "Go环境" "FAIL"
fi

echo ""

# 2. 编译检查
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2. 编译检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 编译自动保存服务
if go build ./service/project/... 2>/dev/null; then
    test_result "自动保存服务编译" "PASS"
else
    test_result "自动保存服务编译" "FAIL"
fi

# 编译认证服务
if go build ./service/shared/auth/... 2>/dev/null; then
    test_result "认证服务编译" "PASS"
else
    test_result "认证服务编译" "FAIL"
fi

# 编译整个项目
if go build ./... 2>/dev/null; then
    test_result "项目完整编译" "PASS"
else
    test_result "项目完整编译" "FAIL"
fi

echo ""

# 3. 单元测试（快速验证）
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3. 单元测试（MVP新增功能）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 测试密码验证器
if go test ./service/shared/auth/ -run TestMVP_PasswordValidation -v 2>/dev/null; then
    test_result "密码验证器单元测试" "PASS"
else
    test_result "密码验证器单元测试" "FAIL"
fi

echo ""

# 4. 集成测试（跳过需要真实环境的测试）
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4. 集成测试准备"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 运行集成测试（跳过标记为Skip的测试）
echo -e "${YELLOW}注意: 完整集成测试需要真实环境，当前仅运行Mock测试${NC}"

if go test ./test/integration/ -run TestMVP_PasswordValidation -v 2>/dev/null; then
    test_result "密码验证集成测试" "PASS"
else
    test_result "密码验证集成测试" "FAIL"
fi

echo ""

# 5. API健康检查（如果服务在运行）
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5. API健康检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查服务是否运行
if curl -s http://localhost:8080/health &> /dev/null; then
    test_result "API健康检查" "PASS"

    # 检查各个模块健康状态
    HEALTH_RESPONSE=$(curl -s http://localhost:8080/health)
    echo "健康检查响应: $HEALTH_RESPONSE"
else
    echo -e "${YELLOW}⚠${NC} API服务未运行（可选）"
fi

echo ""

# 6. 功能验证（手动测试清单）
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "6. 手动功能验证清单"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "请手动验证以下功能："
echo ""
echo "[ ] 1. 密码强度验证"
echo "    - 注册时输入弱密码（如'123456'）应被拒绝"
echo "    - 注册时输入强密码（如'Abc12345'）应成功"
echo ""
echo "[ ] 2. 多端登录限制"
echo "    - 同一用户在5台设备登录应成功"
echo "    - 第6台设备登录应被拒绝"
echo "    - 退出一台设备后，应允许新设备登录"
echo ""
echo "[ ] 3. 自动保存功能"
echo "    - 打开文档编辑"
echo "    - 等待30秒后检查版本历史"
echo "    - 应自动创建新版本"
echo ""
echo "[ ] 4. 完整核心流程"
echo "    - 注册新用户"
echo "    - 登录"
echo "    - 创建项目"
echo "    - 创建章节"
echo "    - 编辑内容（等待自动保存）"
echo "    - AI续写"
echo "    - 发布到书城"
echo "    - 退出登录"
echo ""

# 7. 测试总结
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "测试总结"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "总计测试: $TOTAL_TESTS"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ 所有自动化测试通过！${NC}"
    echo ""
    echo "下一步："
    echo "1. 完成上述手动功能验证"
    echo "2. 执行部署准备检查（运行 ./scripts/deployment_check.sh）"
    echo "3. 部署到内测环境"
    exit 0
else
    echo ""
    echo -e "${RED}✗ 部分测试失败，请修复后重试${NC}"
    exit 1
fi

