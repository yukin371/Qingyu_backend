#!/bin/bash

# 🧪 公开路由测试脚本
# 用于验证后端公开 API 是否正确配置

set -e

# 配置
API_BASE_URL="${API_BASE_URL:-http://localhost:8080}"
API_VERSION="v1"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}🧪 后端公开路由测试${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# 1. 测试书城首页
echo -e "${YELLOW}[Test 1]${NC} GET /api/v1/bookstore/homepage"
response=$(curl -s -X GET "${API_BASE_URL}/api/${API_VERSION}/bookstore/homepage")
code=$(echo $response | jq -r '.code' 2>/dev/null || echo "error")

if [ "$code" = "200" ]; then
    echo -e "${GREEN}✅ PASS${NC} - 首页接口可正常访问"
else
    echo -e "${RED}❌ FAIL${NC} - 响应: $response"
fi
echo ""

# 2. 测试 Banner 列表
echo -e "${YELLOW}[Test 2]${NC} GET /api/v1/bookstore/banners"
response=$(curl -s -X GET "${API_BASE_URL}/api/${API_VERSION}/bookstore/banners")
code=$(echo $response | jq -r '.code' 2>/dev/null || echo "error")

if [ "$code" = "200" ]; then
    echo -e "${GREEN}✅ PASS${NC} - Banner 列表接口可正常访问"
    # 从响应中提取一个 banner ID
    banner_id=$(echo $response | jq -r '.data[0].id' 2>/dev/null || echo "test-banner")
else
    echo -e "${RED}❌ FAIL${NC} - 响应: $response"
    banner_id="test-banner"
fi
echo ""

# 3. 测试 Banner 点击 (关键修复)
echo -e "${YELLOW}[Test 3]${NC} POST /api/v1/bookstore/banners/:id/click 🔑"
echo "          使用 Banner ID: $banner_id"
response=$(curl -s -X POST "${API_BASE_URL}/api/${API_VERSION}/bookstore/banners/${banner_id}/click" \
    -H "Content-Type: application/json")
code=$(echo $response | jq -r '.code' 2>/dev/null || echo "error")
message=$(echo $response | jq -r '.message' 2>/dev/null || echo "unknown")

if [ "$code" = "200" ] || [ "$code" = "404" ]; then
    echo -e "${GREEN}✅ PASS${NC} - Banner 点击接口 (无需认证) 可正常访问"
    echo "   响应代码: $code (成功或 Banner 不存在)"
    echo "   响应消息: $message"
elif grep -q "UNAUTHORIZED\|401\|未登录" <<< "$response"; then
    echo -e "${RED}❌ FAIL${NC} - API 仍然要求认证"
    echo "   响应: $response"
else
    echo -e "${RED}❌ FAIL${NC} - 响应: $response"
fi
echo ""

# 4. 测试推荐首页
echo -e "${YELLOW}[Test 4]${NC} GET /api/v1/recommendation/homepage"
response=$(curl -s -X GET "${API_BASE_URL}/api/${API_VERSION}/recommendation/homepage")
code=$(echo $response | jq -r '.code' 2>/dev/null || echo "error")

if [ "$code" = "200" ]; then
    echo -e "${GREEN}✅ PASS${NC} - 推荐首页接口可正常访问"
else
    echo -e "${RED}❌ FAIL${NC} - 响应: $response"
fi
echo ""

# 5. 测试推荐热门
echo -e "${YELLOW}[Test 5]${NC} GET /api/v1/recommendation/hot"
response=$(curl -s -X GET "${API_BASE_URL}/api/${API_VERSION}/recommendation/hot")
code=$(echo $response | jq -r '.code' 2>/dev/null || echo "error")

if [ "$code" = "200" ]; then
    echo -e "${GREEN}✅ PASS${NC} - 热门推荐接口可正常访问"
else
    echo -e "${RED}❌ FAIL${NC} - 响应: $response"
fi
echo ""

# 6. 测试排行榜
echo -e "${YELLOW}[Test 6]${NC} GET /api/v1/bookstore/rankings/realtime"
response=$(curl -s -X GET "${API_BASE_URL}/api/${API_VERSION}/bookstore/rankings/realtime")
code=$(echo $response | jq -r '.code' 2>/dev/null || echo "error")

if [ "$code" = "200" ]; then
    echo -e "${GREEN}✅ PASS${NC} - 排行榜接口可正常访问"
else
    echo -e "${RED}❌ FAIL${NC} - 响应: $response"
fi
echo ""

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}✨ 测试完成${NC}"
echo -e "${BLUE}========================================${NC}"

