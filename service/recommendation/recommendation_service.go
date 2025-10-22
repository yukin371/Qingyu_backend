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

	// GetHotRecommendations 获取热门推荐
	GetHotRecommendations(ctx context.Context, limit int, days int) ([]string, error)

	// GetHomepageRecommendations 获取首页推荐（混合推荐策略）
	GetHomepageRecommendations(ctx context.Context, userID string, limit int) ([]string, error)

	// GetCategoryRecommendations 获取分类推荐
	GetCategoryRecommendations(ctx context.Context, category string, limit int) ([]string, error)

	// RecordBehavior 记录用户行为
	RecordBehavior(ctx context.Context, behavior *reco.Behavior) error

	// ServiceInfo 服务信息
	ServiceInfo() (string, string)
}

// RecommendationServiceImpl 推荐服务实现
type RecommendationServiceImpl struct {
	behaviorRepo    recoRepo.BehaviorRepository
	profileRepo     recoRepo.ProfileRepository
	itemFeatureRepo recoRepo.ItemFeatureRepository
	hotRepo         recoRepo.HotRecommendationRepository
	serviceName     string
	version         string
}

// NewRecommendationService 创建推荐服务实例
func NewRecommendationService(
	behaviorRepo recoRepo.BehaviorRepository,
	profileRepo recoRepo.ProfileRepository,
	itemFeatureRepo recoRepo.ItemFeatureRepository,
	hotRepo recoRepo.HotRecommendationRepository,
) RecommendationService {
	return &RecommendationServiceImpl{
		behaviorRepo:    behaviorRepo,
		profileRepo:     profileRepo,
		itemFeatureRepo: itemFeatureRepo,
		hotRepo:         hotRepo,
		serviceName:     "RecommendationService",
		version:         "1.0.0",
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

	// 2. 如果画像不存在，返回热门推荐作为冷启动
	if profile == nil || (len(profile.Tags) == 0 && len(profile.Categories) == 0) {
		return s.GetHotRecommendations(ctx, limit, 7) // 返回最近7天的热门书籍
	}

	// 3. 基于用户画像的标签和分类进行推荐
	recommendations := make(map[string]float64)

	// 基于标签推荐
	if len(profile.Tags) > 0 {
		tagItems, err := s.itemFeatureRepo.GetByTags(ctx, profile.Tags, limit*2)
		if err == nil && len(tagItems) > 0 {
			for _, item := range tagItems {
				// 计算标签匹配分数
				score := s.calculateTagMatchScore(profile.Tags, item.Tags)
				recommendations[item.ItemID] = score
			}
		}
	}

	// 基于分类推荐
	if len(profile.Categories) > 0 {
		for category := range profile.Categories {
			categoryItems, err := s.itemFeatureRepo.GetByCategory(ctx, category, limit)
			if err == nil && len(categoryItems) > 0 {
				for _, item := range categoryItems {
					// 累加分类匹配分数
					if existingScore, exists := recommendations[item.ItemID]; exists {
						recommendations[item.ItemID] = existingScore + profile.Categories[category]
					} else {
						recommendations[item.ItemID] = profile.Categories[category]
					}
				}
			}
		}
	}

	// 4. 按分数排序并返回TopN
	sortedItems := sortMapByValue(recommendations)
	result := make([]string, 0, limit)
	for i := 0; i < len(sortedItems) && i < limit; i++ {
		result = append(result, sortedItems[i].ItemID)
	}

	// 5. 如果推荐结果不足，用热门书籍补充
	if len(result) < limit {
		hotBooks, err := s.GetHotRecommendations(ctx, limit-len(result), 7)
		if err == nil {
			result = append(result, hotBooks...)
		}
	}

	return result, nil
}

// GetSimilarItems 获取相似物品（协同过滤）
// 基于物品特征的相似度计算
func (s *RecommendationServiceImpl) GetSimilarItems(ctx context.Context, itemID string, limit int) ([]string, error) {
	// 1. 获取目标物品特征
	targetFeature, err := s.itemFeatureRepo.GetByItemID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get item feature: %w", err)
	}

	if targetFeature == nil {
		// 物品特征不存在，返回热门推荐
		return s.GetHotRecommendations(ctx, limit, 7)
	}

	// 2. 获取相同分类的物品
	similarItems := make(map[string]float64)

	// 基于分类查找相似物品
	if len(targetFeature.Categories) > 0 {
		for _, category := range targetFeature.Categories {
			categoryItems, err := s.itemFeatureRepo.GetByCategory(ctx, category, limit*2)
			if err == nil {
				for _, item := range categoryItems {
					// 跳过自己
					if item.ItemID == itemID {
						continue
					}

					// 计算相似度分数
					score := s.calculateItemSimilarity(targetFeature, item)
					if existingScore, exists := similarItems[item.ItemID]; exists {
						similarItems[item.ItemID] = existingScore + score
					} else {
						similarItems[item.ItemID] = score
					}
				}
			}
		}
	}

	// 基于标签查找相似物品
	if len(targetFeature.Tags) > 0 {
		tagItems, err := s.itemFeatureRepo.GetByTags(ctx, targetFeature.Tags, limit*2)
		if err == nil {
			for _, item := range tagItems {
				if item.ItemID == itemID {
					continue
				}

				tagScore := s.calculateTagMatchScore(targetFeature.Tags, item.Tags)
				if existingScore, exists := similarItems[item.ItemID]; exists {
					similarItems[item.ItemID] = existingScore + tagScore*0.5
				} else {
					similarItems[item.ItemID] = tagScore * 0.5
				}
			}
		}
	}

	// 3. 按相似度排序并返回TopN
	sortedItems := sortMapByValue(similarItems)
	result := make([]string, 0, limit)
	for i := 0; i < len(sortedItems) && i < limit; i++ {
		result = append(result, sortedItems[i].ItemID)
	}

	// 4. 如果相似物品不足，用热门书籍补充
	if len(result) < limit {
		hotBooks, err := s.GetHotRecommendations(ctx, limit-len(result), 7)
		if err == nil {
			for _, bookID := range hotBooks {
				// 避免重复
				isDuplicate := false
				for _, existingID := range result {
					if existingID == bookID {
						isDuplicate = true
						break
					}
				}
				if !isDuplicate {
					result = append(result, bookID)
				}
			}
		}
	}

	return result, nil
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

// GetHotRecommendations 获取热门推荐
func (s *RecommendationServiceImpl) GetHotRecommendations(ctx context.Context, limit int, days int) ([]string, error) {
	hotBooks, err := s.hotRepo.GetHotBooks(ctx, limit, days)
	if err != nil {
		return nil, fmt.Errorf("failed to get hot recommendations: %w", err)
	}
	return hotBooks, nil
}

// GetHomepageRecommendations 获取首页推荐（混合推荐策略）
// 结合热门推荐和个性化推荐
func (s *RecommendationServiceImpl) GetHomepageRecommendations(ctx context.Context, userID string, limit int) ([]string, error) {
	recommendations := make([]string, 0, limit)
	recommendationSet := make(map[string]bool) // 去重

	// 1. 如果有用户ID，获取个性化推荐（占50%）
	if userID != "" {
		personalizedLimit := limit / 2
		personalizedBooks, err := s.GetPersonalizedRecommendations(ctx, userID, personalizedLimit)
		if err == nil {
			for _, bookID := range personalizedBooks {
				if !recommendationSet[bookID] {
					recommendations = append(recommendations, bookID)
					recommendationSet[bookID] = true
				}
			}
		}
	}

	// 2. 添加热门推荐（占30%）
	hotLimit := limit * 3 / 10
	if hotLimit == 0 {
		hotLimit = limit - len(recommendations)
	}
	hotBooks, err := s.GetHotRecommendations(ctx, hotLimit*2, 7) // 多获取一些，用于去重
	if err == nil {
		addedCount := 0
		for _, bookID := range hotBooks {
			if !recommendationSet[bookID] && addedCount < hotLimit {
				recommendations = append(recommendations, bookID)
				recommendationSet[bookID] = true
				addedCount++
			}
		}
	}

	// 3. 添加新书热门（占20%）
	if len(recommendations) < limit {
		newBooksLimit := limit - len(recommendations)
		newBooks, err := s.hotRepo.GetNewPopularBooks(ctx, newBooksLimit*2, 30)
		if err == nil {
			addedCount := 0
			for _, bookID := range newBooks {
				if !recommendationSet[bookID] && addedCount < newBooksLimit {
					recommendations = append(recommendations, bookID)
					recommendationSet[bookID] = true
					addedCount++
				}
			}
		}
	}

	// 4. 如果还不足，添加飙升榜书籍
	if len(recommendations) < limit {
		trendingLimit := limit - len(recommendations)
		trendingBooks, err := s.hotRepo.GetTrendingBooks(ctx, trendingLimit*2)
		if err == nil {
			for _, bookID := range trendingBooks {
				if !recommendationSet[bookID] && len(recommendations) < limit {
					recommendations = append(recommendations, bookID)
					recommendationSet[bookID] = true
				}
			}
		}
	}

	return recommendations, nil
}

// GetCategoryRecommendations 获取分类推荐
func (s *RecommendationServiceImpl) GetCategoryRecommendations(ctx context.Context, category string, limit int) ([]string, error) {
	categoryBooks, err := s.hotRepo.GetHotBooksByCategory(ctx, category, limit, 7)
	if err != nil {
		return nil, fmt.Errorf("failed to get category recommendations: %w", err)
	}
	return categoryBooks, nil
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

// calculateTagMatchScore 计算标签匹配分数
func (s *RecommendationServiceImpl) calculateTagMatchScore(userTags, itemTags map[string]float64) float64 {
	if len(userTags) == 0 || len(itemTags) == 0 {
		return 0
	}

	// 计算加权交集
	score := 0.0
	for tag, userWeight := range userTags {
		if itemWeight, exists := itemTags[tag]; exists {
			score += userWeight * itemWeight
		}
	}

	return score
}

// calculateItemSimilarity 计算物品相似度
func (s *RecommendationServiceImpl) calculateItemSimilarity(item1, item2 *reco.ItemFeature) float64 {
	similarity := 0.0

	// 分类相似度（权重0.4）
	categorySimilarity := calculateCategorySimilarity(item1.Categories, item2.Categories)
	similarity += categorySimilarity * 0.4

	// 标签相似度（权重0.4）
	tagSimilarity := cosineSimilarity(item1.Tags, item2.Tags)
	similarity += tagSimilarity * 0.4

	// 作者相似度（权重0.2）
	authorSimilarity := calculateAuthorSimilarity(item1.Authors, item2.Authors)
	similarity += authorSimilarity * 0.2

	return similarity
}

// calculateCategorySimilarity 计算分类相似度
func calculateCategorySimilarity(categories1, categories2 []string) float64 {
	if len(categories1) == 0 || len(categories2) == 0 {
		return 0
	}

	// 计算Jaccard相似度
	set1 := make(map[string]bool)
	for _, c := range categories1 {
		set1[c] = true
	}

	intersection := 0
	for _, c := range categories2 {
		if set1[c] {
			intersection++
		}
	}

	union := len(categories1) + len(categories2) - intersection
	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// calculateAuthorSimilarity 计算作者相似度
func calculateAuthorSimilarity(authors1, authors2 []string) float64 {
	if len(authors1) == 0 || len(authors2) == 0 {
		return 0
	}

	// 计算Jaccard相似度
	set1 := make(map[string]bool)
	for _, a := range authors1 {
		set1[a] = true
	}

	intersection := 0
	for _, a := range authors2 {
		if set1[a] {
			intersection++
		}
	}

	if intersection > 0 {
		return 1.0 // 有共同作者
	}
	return 0
}

// sortMapByValue 对map按值排序
func sortMapByValue(m map[string]float64) []ScoredItem {
	items := make([]ScoredItem, 0, len(m))
	for k, v := range m {
		items = append(items, ScoredItem{
			ItemID: k,
			Score:  v,
		})
	}
	sortByScore(items)
	return items
}
