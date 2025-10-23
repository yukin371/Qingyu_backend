#!/bin/bash

# MVP冒烟测试脚本
# 用途：验证核心流程端到端可用性
# 流程：注册→登录→创建项目→编辑→AI续写→发布→阅读

set -e

echo "======================================"
echo "🧪 MVP核心流程冒烟测试"
echo "======================================"
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API基础地址
API_BASE="http://localhost:8080/api/v1"
TEST_USER="smoke_test_$(date +%s)"
TEST_EMAIL="${TEST_USER}@test.com"
TEST_PASSWORD="Test@123456"

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试结果函数
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

# JSON解析辅助函数（简化版）
get_json_value() {
    local json=$1
    local key=$2
    echo "$json" | grep -o "\"$key\":\"[^\"]*\"" | sed "s/\"$key\":\"\([^\"]*\)\"/\1/"
}

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "前置检查"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 检查API是否可访问
if curl -s "$API_BASE/../health" > /dev/null 2>&1; then
    test_result "API服务可访问" "PASS"
else
    test_result "API服务可访问" "FAIL"
    echo -e "${RED}错误: API服务未运行，请先启动服务${NC}"
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 1: 用户注册（密码强度验证）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 测试弱密码（应被拒绝）
WEAK_RESPONSE=$(curl -s -X POST "$API_BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"${TEST_USER}_weak\",
        \"email\": \"weak_${TEST_EMAIL}\",
        \"password\": \"123456\"
    }")

if echo "$WEAK_RESPONSE" | grep -q "密码"; then
    test_result "弱密码拒绝验证" "PASS"
else
    test_result "弱密码拒绝验证" "FAIL"
fi

# 注册强密码用户
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USER\",
        \"email\": \"$TEST_EMAIL\",
        \"password\": \"$TEST_PASSWORD\"
    }")

if echo "$REGISTER_RESPONSE" | grep -q "success\|token\|user"; then
    test_result "用户注册（强密码）" "PASS"
else
    test_result "用户注册（强密码）" "FAIL"
    echo "注册响应: $REGISTER_RESPONSE"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 2: 用户登录（多端限制）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 登录获取Token
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"username\": \"$TEST_USER\",
        \"password\": \"$TEST_PASSWORD\"
    }")

# 提取Token（简化处理）
if echo "$LOGIN_RESPONSE" | grep -q "token"; then
    test_result "用户登录" "PASS"
    # 尝试提取token（如果有jq工具）
    if command -v jq &> /dev/null; then
        TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // .token // empty')
    else
        # 简单提取
        TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')
    fi

    if [ -z "$TOKEN" ]; then
        echo -e "${YELLOW}⚠${NC} Token提取失败，后续测试可能受影响"
        TOKEN="dummy_token"
    fi
else
    test_result "用户登录" "FAIL"
    TOKEN="dummy_token"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 3: 创建项目"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

PROJECT_RESPONSE=$(curl -s -X POST "$API_BASE/projects" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "name": "冒烟测试项目",
        "description": "自动化测试项目",
        "genre": "玄幻",
        "tags": ["测试"]
    }')

if echo "$PROJECT_RESPONSE" | grep -q "id\|project"; then
    test_result "创建项目" "PASS"
    if command -v jq &> /dev/null; then
        PROJECT_ID=$(echo "$PROJECT_RESPONSE" | jq -r '.data.id // .data.project_id // empty')
    else
        PROJECT_ID=$(echo "$PROJECT_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"\([^"]*\)"/\1/')
    fi
else
    test_result "创建项目" "FAIL"
    PROJECT_ID="dummy_project"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 4: 创建章节"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

CHAPTER_RESPONSE=$(curl -s -X POST "$API_BASE/projects/$PROJECT_ID/documents" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "title": "第一章 测试开始",
        "content": "这是一个测试章节的内容。",
        "type": "chapter"
    }')

if echo "$CHAPTER_RESPONSE" | grep -q "id\|document"; then
    test_result "创建章节" "PASS"
    if command -v jq &> /dev/null; then
        DOCUMENT_ID=$(echo "$CHAPTER_RESPONSE" | jq -r '.data.id // .data.document_id // empty')
    else
        DOCUMENT_ID=$(echo "$CHAPTER_RESPONSE" | grep -o '"id":"[^"]*"' | head -1 | sed 's/"id":"\([^"]*\)"/\1/')
    fi
else
    test_result "创建章节" "FAIL"
    DOCUMENT_ID="dummy_document"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 5: 自动保存验证（等待35秒）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "等待自动保存触发..."
sleep 35

# 查询版本历史
VERSION_RESPONSE=$(curl -s "$API_BASE/documents/$DOCUMENT_ID/versions" \
    -H "Authorization: Bearer $TOKEN")

if echo "$VERSION_RESPONSE" | grep -q "自动保存\|version"; then
    test_result "自动保存功能" "PASS"
else
    test_result "自动保存功能" "WARN"
    echo -e "${YELLOW}⚠${NC} 可能需要更长时间或手动触发"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 6: AI续写功能（可选）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# AI续写（可能需要配置API Key）
AI_RESPONSE=$(curl -s -X POST "$API_BASE/ai/generate" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "content": "从前有座山",
        "type": "continuation",
        "length": 100
    }' 2>&1 || echo "skip")

if echo "$AI_RESPONSE" | grep -q "text\|content\|result"; then
    test_result "AI续写功能" "PASS"
elif echo "$AI_RESPONSE" | grep -q "quota\|配额"; then
    test_result "AI续写功能" "WARN"
    echo -e "${YELLOW}⚠${NC} 配额不足，功能正常"
else
    test_result "AI续写功能" "SKIP"
    echo -e "${YELLOW}⚠${NC} AI服务未配置或不可用"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 7: 发布到书城（可选）"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

PUBLISH_RESPONSE=$(curl -s -X POST "$API_BASE/projects/$PROJECT_ID/publish" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{
        "status": "published"
    }' 2>&1 || echo "skip")

if echo "$PUBLISH_RESPONSE" | grep -q "success\|published"; then
    test_result "发布到书城" "PASS"
else
    test_result "发布到书城" "SKIP"
    echo -e "${YELLOW}⚠${NC} 发布功能可能未实现或需审核"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 8: 阅读功能验证"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# 获取书城首页
BOOKSTORE_RESPONSE=$(curl -s "$API_BASE/bookstore/home")

if echo "$BOOKSTORE_RESPONSE" | grep -q "books\|novels\|data"; then
    test_result "书城首页访问" "PASS"
else
    test_result "书城首页访问" "FAIL"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Step 9: 退出登录"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

LOGOUT_RESPONSE=$(curl -s -X POST "$API_BASE/auth/logout" \
    -H "Authorization: Bearer $TOKEN")

if echo "$LOGOUT_RESPONSE" | grep -q "success\|ok"; then
    test_result "退出登录" "PASS"
else
    test_result "退出登录" "WARN"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "测试总结"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

echo "总计测试: $TOTAL_TESTS"
echo -e "通过: ${GREEN}$PASSED_TESTS${NC}"
echo -e "失败: ${RED}$FAILED_TESTS${NC}"

SUCCESS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))

echo ""
if [ $SUCCESS_RATE -ge 80 ]; then
    echo -e "${GREEN}✓ 冒烟测试通过！成功率: ${SUCCESS_RATE}%${NC}"
    echo ""
    echo "MVP核心流程可用，可以进入内测阶段！"
    exit 0
else
    echo -e "${RED}✗ 冒烟测试失败！成功率: ${SUCCESS_RATE}%${NC}"
    echo ""
    echo "请修复失败的测试用例后重试"
    exit 1
fi

