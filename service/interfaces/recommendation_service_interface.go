package interfaces

import (
	"context"

	recommendationModel "Qingyu_backend/models/recommendation"
)

// RecommendationService 推荐服务接口
type RecommendationService interface {
	// =========================
	// 推荐获取
	// =========================
	// GetPersonalizedRecommendations 获取个性化推荐（基于用户画像）
	GetPersonalizedRecommendations(ctx context.Context, userID string, limit int) ([]string, error)

	// GetSimilarItems 获取相似物品（协同过滤）
	GetSimilarItems(ctx context.Context, itemID string, limit int) ([]string, error)

	// GetHotRecommendations 获取热门推荐
	GetHotRecommendations(ctx context.Context, limit int, days int) ([]string, error)

	// GetHomepageRecommendations 获取首页推荐（混合推荐策略）
	GetHomepageRecommendations(ctx context.Context, userID string, limit int) ([]string, error)

	// GetCategoryRecommendations 获取分类推荐
	GetCategoryRecommendations(ctx context.Context, category string, limit int) ([]string, error)

	// =========================
	// 用户行为记录
	// =========================
	// RecordBehavior 记录用户行为
	RecordBehavior(ctx context.Context, behavior *recommendationModel.Behavior) error

	// =========================
	// 服务信息
	// =========================
	// ServiceInfo 返回服务信息（服务名和版本）
	ServiceInfo() (string, string)
}
