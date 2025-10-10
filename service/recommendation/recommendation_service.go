package recommendation

import (
	"context"
	"fmt"
	"math"
	"sort"

	reco "Qingyu_backend/models/recommendation/reco"
	recoRepo "Qingyu_backend/repository/interfaces/recommendation"
)

// RecommendationService 推荐服务接口
type RecommendationService interface {
	// GetPersonalizedRecommendations 获取个性化推荐（基于用户画像）
	GetPersonalizedRecommendations(ctx context.Context, userID string, limit int) ([]string, error)

	// GetSimilarItems 获取相似物品（协同过滤）
	GetSimilarItems(ctx context.Context, itemID string, limit int) ([]string, error)

	// RecordBehavior 记录用户行为
	RecordBehavior(ctx context.Context, behavior *reco.Behavior) error

	// ServiceInfo 服务信息
	ServiceInfo() (string, string)
}

// RecommendationServiceImpl 推荐服务实现
type RecommendationServiceImpl struct {
	behaviorRepo recoRepo.BehaviorRepository
	profileRepo  recoRepo.ProfileRepository
	serviceName  string
	version      string
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(
	behaviorRepo recoRepo.BehaviorRepository,
	profileRepo recoRepo.ProfileRepository,
) RecommendationService {
	return &RecommendationServiceImpl{
		behaviorRepo: behaviorRepo,
		profileRepo:  profileRepo,
		serviceName:  "RecommendationService",
		version:      "1.0.0",
	}
}

// GetPersonalizedRecommendations 获取个性化推荐
// 基于用户画像的内容推荐（标签匹配）
func (s *RecommendationServiceImpl) GetPersonalizedRecommendations(ctx context.Context, userID string, limit int) ([]string, error) {
	// 1. 获取用户画像
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	// 2. 如果画像不存在，返回空或热门推荐
	if profile == nil {
		// TODO: 返回热门书籍作为冷启动
		return []string{}, nil
	}

	// 3. 基于标签/分类/作者偏好计算推荐（简化版本）
	// TODO: 实现基于画像的标签匹配算法
	// 当前返回空列表作为占位
	recommendations := []string{}

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

// GetSimilarItems 获取相似物品（协同过滤）
// 基于物品的协同过滤算法
func (s *RecommendationServiceImpl) GetSimilarItems(ctx context.Context, itemID string, limit int) ([]string, error) {
	// TODO: 实现基于物品的协同过滤
	// 1. 获取物品特征
	// 2. 计算相似度（余弦相似度、皮尔逊相关系数等）
	// 3. 返回TopN相似物品

	// 当前返回空列表作为占位
	similarItems := []string{}

	if len(similarItems) > limit {
		similarItems = similarItems[:limit]
	}

	return similarItems, nil
}

// RecordBehavior 记录用户行为
func (s *RecommendationServiceImpl) RecordBehavior(ctx context.Context, behavior *reco.Behavior) error {
	if behavior.UserID == "" || behavior.ItemID == "" {
		return fmt.Errorf("user_id and item_id are required")
	}

	err := s.behaviorRepo.Create(ctx, behavior)
	if err != nil {
		return fmt.Errorf("failed to record behavior: %w", err)
	}

	// TODO: 异步触发画像更新

	return nil
}

// ServiceInfo 返回服务信息
func (s *RecommendationServiceImpl) ServiceInfo() (string, string) {
	return s.serviceName, s.version
}

// ===== 辅助函数 =====

// cosineSimilarity 计算余弦相似度（占位实现）
func cosineSimilarity(vec1, vec2 map[string]float64) float64 {
	if len(vec1) == 0 || len(vec2) == 0 {
		return 0
	}

	dotProduct := 0.0
	norm1 := 0.0
	norm2 := 0.0

	// 计算所有键的并集
	allKeys := make(map[string]bool)
	for k := range vec1 {
		allKeys[k] = true
	}
	for k := range vec2 {
		allKeys[k] = true
	}

	for k := range allKeys {
		v1 := vec1[k]
		v2 := vec2[k]
		dotProduct += v1 * v2
		norm1 += v1 * v1
		norm2 += v2 * v2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

// ScoredItem 带分数的物品
type ScoredItem struct {
	ItemID string
	Score  float64
}

// sortByScore 按分数排序
func sortByScore(items []ScoredItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Score > items[j].Score
	})
}
