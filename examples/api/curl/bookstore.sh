#!/bin/bash

# 青羽写作平台 API - 书城API示例
# 使用curl进行API调用

BASE_URL="http://localhost:9090/api/v1"
TOKEN="your-jwt-token-here"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}=== 青羽写作平台 书城API示例 ===${NC}\n"

# 1. 获取首页数据
echo -e "${YELLOW}1. 获取首页数据${NC}"
echo "GET /bookstore/home"
curl -X GET "${BASE_URL}/bookstore/home" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 2. 获取书籍列表
echo -e "${YELLOW}2. 获取书籍列表${NC}"
echo "GET /bookstore/books"
curl -X GET "${BASE_URL}/bookstore/books?page=1&limit=10" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 3. 搜索书籍
echo -e "${YELLOW}3. 搜索书籍（按标题）${NC}"
echo "GET /bookstore/books/search/title"
curl -X GET "${BASE_URL}/bookstore/books/search/title?title=玄幻&page=1&limit=5" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 4. 按作者搜索
echo -e "${YELLOW}4. 按作者搜索${NC}"
echo "GET /bookstore/books/search/author"
curl -X GET "${BASE_URL}/bookstore/books/search/author?author=唐家三少&page=1&limit=5" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 5. 按标签筛选
echo -e "${YELLOW}5. 按标签筛选${NC}"
echo "GET /bookstore/books/tags"
curl -X GET "${BASE_URL}/bookstore/books/tags?tags=玄幻,修真&page=1&limit=5" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 6. 获取书籍详情
echo -e "${YELLOW}6. 获取书籍详情${NC}"
echo "GET /bookstore/books/{id}"
BOOK_ID="book123"
curl -X GET "${BASE_URL}/bookstore/books/${BOOK_ID}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 7. 获取书籍章节列表
echo -e "${YELLOW}7. 获取书籍章节列表${NC}"
echo "GET /bookstore/books/{id}/chapters"
curl -X GET "${BASE_URL}/bookstore/books/${BOOK_ID}/chapters?page=1&limit=20" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 8. 获取章节内容
echo -e "${YELLOW}8. 获取章节内容（需要认证）${NC}"
echo "GET /bookstore/chapters/{id}"
CHAPTER_ID="chapter123"
curl -X GET "${BASE_URL}/bookstore/chapters/${CHAPTER_ID}" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 9. 获取相似书籍推荐
echo -e "${YELLOW}9. 获取相似书籍推荐${NC}"
echo "GET /bookstore/books/{id}/similar"
curl -X GET "${BASE_URL}/bookstore/books/${BOOK_ID}/similar?limit=10" \
  -H "Content-Type: application/json" \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n"

# 10. 提交书籍评分
echo -e "${YELLOW}10. 提交书籍评分（需要认证）${NC}"
echo "POST /bookstore/books/{id}/rating"
curl -X POST "${BASE_URL}/bookstore/books/${BOOK_ID}/rating" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "rating": 5,
    "comment": "非常好看！"
  }' \
  -w "\nHTTP Status: %{http_code}\n" \
  -s | jq '.'

echo -e "\n${GREEN}=== 示例完成 ===${NC}"
