#!/bin/bash
# 阅读端功能自动化测试脚本

set -e  # 遇到错误立即退出

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 配置
BASE_URL="${BASE_URL:-http://localhost:8080}"
API_PREFIX="/api/v1"
OUTPUT_DIR="test_results"

# 创建输出目录
mkdir -p "$OUTPUT_DIR"

# 测试结果统计
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 打印标题
print_header() {
    echo ""
    echo "========================================"
    echo "$1"
    echo "========================================"
}

# 打印测试结果
print_result() {
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    if [ $2 -eq 0 ]; then
        echo -e "${GREEN}[OK]${NC} $1"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}[FAIL]${NC} $1"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
}

# 测试 API
test_api() {
    local name=$1
    local url=$2
    local method=${3:-GET}
    local data=$4

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$BASE_URL$API_PREFIX$url" || echo "000")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$BASE_URL$API_PREFIX$url" || echo "000")
    fi

    http_code=$(echo "$response" | tail -n 1)
    body=$(echo "$response" | sed '$d')

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        print_result "$name" 0
        echo "$body" > "$OUTPUT_DIR/$(echo $name | tr ' ' '_').json"
        return 0
    else
        print_result "$name (HTTP $http_code)" 1
        return 1
    fi
}

# 从 JSON 响应中提取第一个书籍 ID
get_first_book_id() {
    cat "$OUTPUT_DIR/获取书籍列表.json" | grep -o '"id":"[^"]*"' | head -n 1 | cut -d'"' -f4
}

# 从 JSON 响应中提取第一个章节 ID
get_first_chapter_id() {
    cat "$OUTPUT_DIR/获取章节列表.json" | grep -o '"id":"[^"]*"' | head -n 1 | cut -d'"' -f4
}

echo "========================================"
echo "青羽阅读端功能自动化测试"
echo "========================================"
echo "测试服务器: $BASE_URL"
echo "测试时间: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# ========================================
# 一、书城浏览功能测试
# ========================================
print_header "1. 书城浏览功能测试"

test_api "获取书籍列表" "/bookstore/books?page=1&pageSize=10"
test_api "分页测试-第2页" "/bookstore/books?page=2&pageSize=5"
test_api "按字数排序" "/bookstore/books?sortBy=word_count&sortOrder=desc&limit=10"
test_api "按章节数排序" "/bookstore/books?sortBy=chapter_count&sortOrder=desc&limit=10"
test_api "玄幻分类筛选" "/bookstore/books?category=玄幻&limit=10"
test_api "仙侠分类筛选" "/bookstore/books?category=仙侠&limit=10"
test_api "都市分类筛选" "/bookstore/books?category=都市&limit=10"

# 搜索功能（如果有书籍数据）
if [ -f "$OUTPUT_DIR/获取书籍列表.json" ]; then
    # 从结果中提取一个书名进行搜索测试
    test_api "搜索功能测试" "/bookstore/books/search?keyword=书"
fi

# ========================================
# 二、榜单功能测试
# ========================================
print_header "2. 榜单功能测试"

test_api "热门榜-按字数" "/bookstore/books?sortBy=word_count&sortOrder=desc&limit=20"
test_api "热门榜-按章节" "/bookstore/books?sortBy=chapter_count&sortOrder=desc&limit=20"
test_api "热门书籍标记" "/bookstore/books?is_hot=true&limit=20"
test_api "推荐书籍" "/bookstore/books?is_recommended=true&limit=20"
test_api "精选书籍" "/bookstore/books?is_featured=true&limit=20"
test_api "新书榜" "/bookstore/books?sortBy=created_at&sortOrder=desc&limit=20"
test_api "最近更新榜" "/bookstore/books?sortBy=updated_at&sortOrder=desc&limit=20"

# 分类榜单
test_api "玄幻热门榜" "/bookstore/books?category=玄幻&sortBy=word_count&sortOrder=desc&limit=20"
test_api "仙侠热门榜" "/bookstore/books?category=仙侠&sortBy=word_count&sortOrder=desc&limit=20"
test_api "都市热门榜" "/bookstore/books?category=都市&sortBy=word_count&sortOrder=desc&limit=20"

# ========================================
# 三、书籍详情测试
# ========================================
print_header "3. 书籍详情测试"

if [ -f "$OUTPUT_DIR/获取书籍列表.json" ]; then
    BOOK_ID=$(get_first_book_id)
    if [ -n "$BOOK_ID" ]; then
        echo "使用书籍 ID: $BOOK_ID"
        test_api "获取书籍详情" "/bookstore/books/$BOOK_ID"
        test_api "获取章节列表" "/bookstore/books/$BOOK_ID/chapters?page=1&pageSize=50"

        # 如果有章节，测试章节内容
        if [ -f "$OUTPUT_DIR/获取章节列表.json" ]; then
            CHAPTER_ID=$(get_first_chapter_id)
            if [ -n "$CHAPTER_ID" ]; then
                echo "使用章节 ID: $CHAPTER_ID"
                test_api "获取章节内容" "/bookstore/chapters/$CHAPTER_ID"
            fi
        fi
    fi
fi

# ========================================
# 四、评分系统测试（只读操作）
# ========================================
print_header "4. 评分系统测试"

if [ -n "$BOOK_ID" ]; then
    test_api "获取书籍评分列表" "/reading/books/$BOOK_ID/ratings?page=1&limit=10"
    test_api "获取平均评分" "/reading/books/$BOOK_ID/average-rating"
    test_api "获取评分分布" "/reading/books/$BOOK_ID/rating-distribution"
fi

# ========================================
# 五、推荐系统测试
# ========================================
print_header "5. 推荐系统测试"

test_api "通用推荐" "/recommendation/books?limit=10"
if [ -n "$BOOK_ID" ]; then
    test_api "相似书籍推荐" "/recommendation/similar/$BOOK_ID?limit=10"
fi

# ========================================
# 六、统计信息测试
# ========================================
print_header "6. 统计信息测试"

if [ -n "$BOOK_ID" ]; then
    test_api "书籍统计信息" "/reading/books/$BOOK_ID/statistics"
fi
test_api "分类统计" "/reading/statistics/categories" || echo "接口可能未实现"

# ========================================
# 测试总结
# ========================================
print_header "测试总结"

echo ""
echo "总测试项: $TOTAL_TESTS"
echo -e "${GREEN}通过: $PASSED_TESTS${NC}"
echo -e "${RED}失败: $FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}✓ 所有测试通过！${NC}"
    PASS_RATE=100
else
    PASS_RATE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "${YELLOW}通过率: ${PASS_RATE}%${NC}"
fi

# 生成测试报告
REPORT_FILE="$OUTPUT_DIR/test_report_$(date '+%Y%m%d_%H%M%S').txt"
{
    echo "========================================"
    echo "青羽阅读端功能测试报告"
    echo "========================================"
    echo "测试时间: $(date '+%Y-%m-%d %H:%M:%S')"
    echo "测试服务器: $BASE_URL"
    echo ""
    echo "测试结果统计:"
    echo "- 总测试项: $TOTAL_TESTS"
    echo "- 通过项: $PASSED_TESTS"
    echo "- 失败项: $FAILED_TESTS"
    echo "- 通过率: ${PASS_RATE}%"
    echo ""
    echo "详细结果请查看: $OUTPUT_DIR 目录"
} > "$REPORT_FILE"

echo ""
echo "测试报告已保存: $REPORT_FILE"
echo "详细响应数据保存在: $OUTPUT_DIR 目录"
echo ""

# 退出码
if [ $FAILED_TESTS -eq 0 ]; then
    exit 0
else
    exit 1
fi

