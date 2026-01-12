package recommendation

import (
	recommendationModel "Qingyu_backend/models/recommendation"
	"context"
)

// RecommendationRepository 推荐相关Repository
type RecommendationRepository interface {
	// 用户行为记录
	RecordBehavior(ctx context.Context, behavior *recommendationModel.UserBehavior) error
	GetUserBehaviors(ctx context.Context, userID string, limit int) ([]*recommendationModel.UserBehavior, error)
	GetItemBehaviors(ctx context.Context, itemID string, limit int) ([]*recommendationModel.UserBehavior, error)

	// 推荐结果存储（可选，可以完全用缓存）
	// SaveRecommendations(ctx context.Context, userID string, items []*recommendationModel.RecommendedItem) error
	// GetRecommendations(ctx context.Context, userID string) ([]*recommendationModel.RecommendedItem, error)

	// Health 健康检查
	Health(ctx context.Context) error
}
