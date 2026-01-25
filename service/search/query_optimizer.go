package search

import (
	"fmt"
	"regexp"
	"strings"

	"Qingyu_backend/models/search"
)

// QueryOptimizer 查询优化器
type QueryOptimizer struct {
	maxResults      int
	maxPageSize     int
	minQueryLength  int
	enableFuzziness bool
}

// NewQueryOptimizer 创建查询优化器
func NewQueryOptimizer(maxResults, maxPageSize, minQueryLength int, enableFuzziness bool) *QueryOptimizer {
	return &QueryOptimizer{
		maxResults:      maxResults,
		maxPageSize:     maxPageSize,
		minQueryLength:  minQueryLength,
		enableFuzziness: enableFuzziness,
	}
}

// Optimize 优化搜索请求
func (o *QueryOptimizer) Optimize(req *search.SearchRequest) (*search.SearchRequest, error) {
	// 1. 验证分页参数
	if req.Page < 1 {
		req.Page = 1
	}

	if req.PageSize < 1 {
		req.PageSize = 20
	}

	// 2. 限制最大分页大小
	if req.PageSize > o.maxPageSize {
		req.PageSize = o.maxPageSize
	}

	// 3. 限制最大结果数
	if req.Page*req.PageSize > o.maxResults {
		return nil, fmt.Errorf("page size exceeds maximum (%d)", o.maxResults)
	}

	// 4. 验证查询长度
	if len(req.Query) < o.minQueryLength {
		return nil, fmt.Errorf("query too short (minimum %d characters)", o.minQueryLength)
	}

	// 5. 清理查询字符串
	req.Query = o.sanitizeQuery(req.Query)

	// 6. 过滤特殊字符
	if strings.ContainsAny(req.Query, "&|!(){}[]^~*?:\\/") {
		return nil, fmt.Errorf("query contains invalid characters")
	}

	return req, nil
}

// sanitizeQuery 清理查询字符串
func (o *QueryOptimizer) sanitizeQuery(query string) string {
	// 移除前后空格
	query = strings.TrimSpace(query)

	// 移除多余空格
	space := regexp.MustCompile(`\s+`)
	query = space.ReplaceAllString(query, " ")

	return query
}
