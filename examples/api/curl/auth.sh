#!/bin/bash

# 青羽写作平台 API - 认证示例
# 使用curl进行API调用

BASE_URL="http://localhost:9090/api/v1"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== 青羽写作平台 API 认证示例 ===${NC}\n"

# 1. 用户注册
echo -e "${YELLOW}1. 用户注册${NC}"
echo "POST /auth/register"
curl -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "nickname": "测试用户"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s

echo -e "\n"

# 2. 用户登录
echo -e "${YELLOW}2. 用户登录${NC}"
echo "POST /auth/login"
LOGIN_RESPONSE=$(curl -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "SecurePass123!"
  }' \
  -s)

echo "$LOGIN_RESPONSE" | jq '.'

# 提取token
TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.data.token // .data.access_token // empty')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo -e "\n${RED}登录失败，无法获取token${RED}"
  exit 1
fi

echo -e "\n${GREEN}获取到Token: ${TOKEN:0:20}...${NC}\n"

# 3. 使用Token访问受保护的API
echo -e "${YELLOW}3. 获取用户信息（使用Token）${NC}"
echo "GET /user/profile"
curl -X GET "${BASE_URL}/user/profile" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 4. 刷新Token
echo -e "${YELLOW}4. 刷新Token${NC}"
echo "POST /auth/refresh"
curl -X POST "${BASE_URL}/auth/refresh" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 5. 登出
echo -e "${YELLOW}5. 登出${NC}"
echo "POST /auth/logout"
curl -X POST "${BASE_URL}/auth/logout" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n${GREEN}=== 示例完成 ===${NC}"
