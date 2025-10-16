package recommendation

import (
	"context"

	reco "Qingyu_backend/models/recommendation/reco"
)

// ItemFeatureRepository 物品特征仓储接口
// 用于存储和查询书籍的特征信息，支持内容推荐和相似度计算
type ItemFeatureRepository interface {
	// Create 创建物品特征
	Create(ctx context.Context, feature *reco.ItemFeature) error

	// Upsert 更新或创建物品特征
	Upsert(ctx context.Context, feature *reco.ItemFeature) error

	// GetByItemID 根据物品ID获取特征
	GetByItemID(ctx context.Context, itemID string) (*reco.ItemFeature, error)

	// BatchGetByItemIDs 批量获取物品特征
	BatchGetByItemIDs(ctx context.Context, itemIDs []string) ([]*reco.ItemFeature, error)

	// GetByCategory 根据分类获取物品特征列表
	GetByCategory(ctx context.Context, category string, limit int) ([]*reco.ItemFeature, error)

	// GetByTags 根据标签获取物品特征列表（标签权重匹配）
	GetByTags(ctx context.Context, tags map[string]float64, limit int) ([]*reco.ItemFeature, error)

	// Delete 删除物品特征
	Delete(ctx context.Context, itemID string) error

	// Health 健康检查
	Health(ctx context.Context) error
}
