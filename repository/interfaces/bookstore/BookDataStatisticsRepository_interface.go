package bookstore

import (
	bookstore2 "Qingyu_backend/models/bookstore"
	"context"
)

// BookDataStatisticsRepository 书籍数据统计接口 - 专注于书籍数据的统计和计数操作
// 用于数据统计、分析报表等场景
//
// 注意：此接口用于Book模型的统计数据，不同于BookStatisticsRepository（用于BookStatistics模型）
//
// 职责：
// - 统计书籍数量（按分类、作者、状态等）
// - 获取统计概览信息
// - 更新计数（浏览量、点赞数、评论数等）
//
// 方法总数：9个
type BookDataStatisticsRepository interface {
	// CountByCategory 统计分类下的书籍数量
	CountByCategory(ctx context.Context, categoryID string) (int64, error)

	// CountByAuthor 统计作者的书籍数量
	CountByAuthor(ctx context.Context, author string) (int64, error)

	// CountByStatus 统计指定状态的书籍数量
	CountByStatus(ctx context.Context, status bookstore2.BookStatus) (int64, error)

	// CountByFilter 根据过滤器统计书籍数量
	CountByFilter(ctx context.Context, filter *bookstore2.BookFilter) (int64, error)

	// GetStats 获取书籍统计概览信息
	// 包括总书数、已发布数、草稿数、推荐数、精选数等
	GetStats(ctx context.Context) (*bookstore2.BookStats, error)

	// IncrementViewCount 增加浏览计数
	IncrementViewCount(ctx context.Context, bookID string) error

	// IncrementLikeCount 增加点赞数
	IncrementLikeCount(ctx context.Context, bookID string) error

	// IncrementCommentCount 增加评论数
	IncrementCommentCount(ctx context.Context, bookID string) error

	// UpdateRating 更新评分
	UpdateRating(ctx context.Context, bookID string, rating float64) error
}
