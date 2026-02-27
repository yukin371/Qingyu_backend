#!/bin/bash

# 青羽写作平台 API - 管理员API示例
# 使用curl进行API调用

BASE_URL="http://localhost:9090/api/v1"
ADMIN_TOKEN="your-admin-jwt-token-here"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# 检查token
if [ "$ADMIN_TOKEN" = "your-admin-jwt-token-here" ]; then
  echo -e "${RED}错误: 请先设置有效的管理员Token${NC}"
  echo "编辑此文件，将 ADMIN_TOKEN 变量设置为有效的Token"
  exit 1
fi

echo -e "${GREEN}=== 青羽写作平台 管理员API示例 ===${NC}\n"

# 1. 获取用户列表
echo -e "${YELLOW}1. 获取用户列表${NC}"
echo "GET /admin/users"
curl -X GET "${BASE_URL}/admin/users?page=1&limit=10" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 2. 创建用户
echo -e "${YELLOW}2. 创建用户${NC}"
echo "POST /admin/users"
curl -X POST "${BASE_URL}/admin/users" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "SecurePass123!",
    "nickname": "新用户",
    "role": "user"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 3. 获取权限列表
echo -e "${YELLOW}3. 获取权限列表${NC}"
echo "GET /admin/permissions"
curl -X GET "${BASE_URL}/admin/permissions" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 4. 创建权限
echo -e "${YELLOW}4. 创建权限${NC}"
echo "POST /admin/permissions"
curl -X POST "${BASE_URL}/admin/permissions" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "code": "test.permission",
    "name": "测试权限",
    "description": "这是一个测试权限",
    "category": "test"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 5. 获取权限模板列表
echo -e "${YELLOW}5. 获取权限模板列表${NC}"
echo "GET /admin/permission-templates"
curl -X GET "${BASE_URL}/admin/permission-templates" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 6. 创建权限模板
echo -e "${YELLOW}6. 创建权限模板${NC}"
echo "POST /admin/permission-templates"
curl -X POST "${BASE_URL}/admin/permission-templates" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "测试模板",
    "code": "test_template",
    "description": "这是一个测试权限模板",
    "permissions": ["user.read", "content.read"],
    "category": "test"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 7. 获取审计日志
echo -e "${YELLOW}7. 获取审计日志${NC}"
echo "GET /admin/audit/trail"
curl -X GET "${BASE_URL}/admin/audit/trail?page=1&size=20" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 8. 获取用户增长趋势
echo -e "${YELLOW}8. 获取用户增长趋势${NC}"
echo "GET /admin/analytics/user-growth"
curl -X GET "${BASE_URL}/admin/analytics/user-growth?start_date=2026-01-01&end_date=2026-01-31&interval=daily" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 9. 获取内容统计
echo -e "${YELLOW}9. 获取内容统计${NC}"
echo "GET /admin/analytics/content-statistics"
curl -X GET "${BASE_URL}/admin/analytics/content-statistics" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 10. 创建公告
echo -e "${YELLOW}10. 创建公告${NC}"
echo "POST /admin/announcements"
curl -X POST "${BASE_URL}/admin/announcements" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "系统维护通知",
    "content": "系统将于今晚22:00-24:00进行维护",
    "type": "maintenance",
    "priority": "high",
    "is_pinned": true
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 11. 导出书籍数据
echo -e "${YELLOW}11. 导出书籍数据${NC}"
echo "GET /admin/content/books/export"
curl -X GET "${BASE_URL}/admin/content/books/export?format=csv" \
  -H "Authorization: Bearer ${ADMIN_TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s

echo -e "\n${GREEN}=== 示例完成 ===${NC}"
