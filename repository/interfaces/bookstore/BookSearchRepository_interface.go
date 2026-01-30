package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
)

// BookSearchRepository 书籍搜索接口 - 专注于书籍的搜索和高级筛选
// 用于书城搜索页面、高级筛选场景
//
// 职责：
// - 关键词搜索
// - 高级筛选（价格、分类、标签等）
// - 价格区间查询
//
// 方法总数：3个
type BookSearchRepository interface {
	// Search 搜索书籍（简化版本）
	// 根据关键词搜索书籍的标题、作者、简介
	Search(ctx context.Context, keyword string, limit, offset int) ([]*bookstore2.Book, error)

	// SearchWithFilter 使用过滤器搜索书籍
	// 支持多条件组合筛选（分类、作者、标签、状态等）
	SearchWithFilter(ctx context.Context, filter *bookstore2.BookFilter) ([]*bookstore2.Book, error)

	// GetByPriceRange 按价格区间获取书籍
	// 用于价格筛选场景
	GetByPriceRange(ctx context.Context, minPrice, maxPrice float64, limit, offset int) ([]*bookstore2.Book, error)
}
