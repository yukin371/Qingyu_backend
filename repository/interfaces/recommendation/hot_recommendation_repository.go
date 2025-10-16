package recommendation

import (
	"context"
)

// HotRecommendationRepository 热门推荐仓储接口
// 用于获取热门书籍数据，基于浏览量、收藏量、评分等指标
type HotRecommendationRepository interface {
	// GetHotBooks 获取热门书籍列表
	// 参数:
	//   ctx: 上下文
	//   limit: 返回数量
	//   days: 统计天数（例如：最近7天的热门）
	// 返回: 书籍ID列表
	GetHotBooks(ctx context.Context, limit int, days int) ([]string, error)

	// GetHotBooksByCategory 获取分类下的热门书籍
	// 参数:
	//   ctx: 上下文
	//   category: 分类名称
	//   limit: 返回数量
	//   days: 统计天数
	// 返回: 书籍ID列表
	GetHotBooksByCategory(ctx context.Context, category string, limit int, days int) ([]string, error)

	// GetTrendingBooks 获取正在飙升的书籍（增长趋势）
	// 参数:
	//   ctx: 上下文
	//   limit: 返回数量
	// 返回: 书籍ID列表
	GetTrendingBooks(ctx context.Context, limit int) ([]string, error)

	// GetNewPopularBooks 获取新书中的热门书籍
	// 参数:
	//   ctx: 上下文
	//   limit: 返回数量
	//   daysThreshold: 新书天数阈值（例如：30天内上架的书）
	// 返回: 书籍ID列表
	GetNewPopularBooks(ctx context.Context, limit int, daysThreshold int) ([]string, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
