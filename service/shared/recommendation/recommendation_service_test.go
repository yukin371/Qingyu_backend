package recommendation

import (
	recModel "Qingyu_backend/models/recommendation"
	"context"
	"fmt"
	"testing"
	"time"
)

// MockRecommendationRepository Mock推荐Repository
type MockRecommendationRepository struct {
	behaviors map[string][]*recModel.UserBehavior
}

func NewMockRecommendationRepository() *MockRecommendationRepository {
	return &MockRecommendationRepository{
		behaviors: make(map[string][]*recModel.UserBehavior),
	}
}

func (m *MockRecommendationRepository) RecordBehavior(ctx context.Context, behavior *recModel.UserBehavior) error {
	behavior.ID = fmt.Sprintf("beh_%d", len(m.behaviors)+1)
	behavior.CreatedAt = time.Now()

	key := "user_" + behavior.UserID
	m.behaviors[key] = append(m.behaviors[key], behavior)
	return nil
}

func (m *MockRecommendationRepository) GetUserBehaviors(ctx context.Context, userID string, limit int) ([]*recModel.UserBehavior, error) {
	key := "user_" + userID
	behaviors := m.behaviors[key]

	if len(behaviors) > limit {
		behaviors = behaviors[:limit]
	}

	return behaviors, nil
}

func (m *MockRecommendationRepository) GetItemBehaviors(ctx context.Context, itemID string, limit int) ([]*recModel.UserBehavior, error) {
	var result []*recModel.UserBehavior

	for _, behaviorList := range m.behaviors {
		for _, b := range behaviorList {
			if b.ItemID == itemID {
				result = append(result, b)
			}
		}
	}

	if len(result) > limit {
		result = result[:limit]
	}

	return result, nil
}

func (m *MockRecommendationRepository) Health(ctx context.Context) error {
	return nil
}

// MockCacheClient Mock缓存客户端
type MockCacheClient struct {
	data map[string]string
}

func NewMockCacheClient() *MockCacheClient {
	return &MockCacheClient{
		data: make(map[string]string),
	}
}

func (m *MockCacheClient) Get(ctx context.Context, key string) (string, error) {
	val, exists := m.data[key]
	if !exists {
		return "", fmt.Errorf("key not found")
	}
	return val, nil
}

func (m *MockCacheClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	m.data[key] = value
	return nil
}

func (m *MockCacheClient) Delete(ctx context.Context, key string) error {
	delete(m.data, key)
	return nil
}

// ============ 测试用例 ============

// TestRecordUserBehavior 测试记录用户行为
func TestRecordUserBehavior(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()
	req := &RecordBehaviorRequest{
		UserID:     "user123",
		ItemID:     "book_001",
		ItemType:   "book",
		ActionType: "view",
		Duration:   300,
	}

	err := service.RecordUserBehavior(ctx, req)
	if err != nil {
		t.Fatalf("记录用户行为失败: %v", err)
	}

	// 验证记录已保存
	behaviors, _ := service.GetUserBehaviors(ctx, "user123", 10)
	if len(behaviors) != 1 {
		t.Errorf("行为记录数量错误: %d", len(behaviors))
	}

	if behaviors[0].ItemID != "book_001" {
		t.Errorf("物品ID错误: %s", behaviors[0].ItemID)
	}

	t.Logf("记录用户行为测试通过: %+v", behaviors[0])
}

// TestGetUserBehaviors 测试获取用户行为记录
func TestGetUserBehaviors(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 记录多条行为
	for i := 1; i <= 5; i++ {
		req := &RecordBehaviorRequest{
			UserID:     "user123",
			ItemID:     fmt.Sprintf("book_%03d", i),
			ItemType:   "book",
			ActionType: "view",
			Duration:   int64(i * 100),
		}
		service.RecordUserBehavior(ctx, req)
	}

	// 获取行为记录
	behaviors, err := service.GetUserBehaviors(ctx, "user123", 10)
	if err != nil {
		t.Fatalf("获取用户行为失败: %v", err)
	}

	if len(behaviors) != 5 {
		t.Errorf("行为记录数量错误，期望5，实际%d", len(behaviors))
	}

	t.Logf("获取用户行为测试通过，共%d条记录", len(behaviors))
}

// TestGetPersonalizedRecommendations 测试个性化推荐
func TestGetPersonalizedRecommendations(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 记录用户行为
	for i := 1; i <= 3; i++ {
		req := &RecordBehaviorRequest{
			UserID:     "user123",
			ItemID:     fmt.Sprintf("book_%03d", i),
			ItemType:   "book",
			ActionType: "view",
			Duration:   300,
		}
		service.RecordUserBehavior(ctx, req)
	}

	// 获取个性化推荐
	recommendations, err := service.GetPersonalizedRecommendations(ctx, "user123", 5)
	if err != nil {
		t.Fatalf("获取个性化推荐失败: %v", err)
	}

	if len(recommendations) == 0 {
		t.Error("推荐列表为空")
	}

	t.Logf("个性化推荐测试通过，共%d条推荐", len(recommendations))
	for i, rec := range recommendations {
		t.Logf("  [%d] %s (score: %.2f) - %s", i+1, rec.ItemID, rec.Score, rec.Reason)
	}
}

// TestGetPersonalizedRecommendations_NoHistory 测试无历史记录的推荐
func TestGetPersonalizedRecommendations_NoHistory(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 获取推荐（无历史记录）
	recommendations, err := service.GetPersonalizedRecommendations(ctx, "new_user", 5)
	if err != nil {
		t.Fatalf("获取个性化推荐失败: %v", err)
	}

	// 应该返回默认的热门推荐
	if len(recommendations) == 0 {
		t.Error("推荐列表为空，应返回默认推荐")
	}

	t.Logf("无历史记录推荐测试通过，共%d条默认推荐", len(recommendations))
}

// TestGetSimilarItems 测试相似内容推荐
func TestGetSimilarItems(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 模拟多个用户浏览同一本书
	for i := 1; i <= 3; i++ {
		req1 := &RecordBehaviorRequest{
			UserID:     fmt.Sprintf("user%d", i),
			ItemID:     "book_001",
			ItemType:   "book",
			ActionType: "view",
		}
		service.RecordUserBehavior(ctx, req1)

		// 这些用户还浏览了其他书
		req2 := &RecordBehaviorRequest{
			UserID:     fmt.Sprintf("user%d", i),
			ItemID:     fmt.Sprintf("book_%03d", i+10),
			ItemType:   "book",
			ActionType: "view",
		}
		service.RecordUserBehavior(ctx, req2)
	}

	// 获取相似推荐
	recommendations, err := service.GetSimilarItems(ctx, "book_001", 5)
	if err != nil {
		t.Fatalf("获取相似推荐失败: %v", err)
	}

	if len(recommendations) == 0 {
		t.Error("相似推荐列表为空")
	}

	t.Logf("相似推荐测试通过，共%d条推荐", len(recommendations))
	for i, rec := range recommendations {
		t.Logf("  [%d] %s (score: %.2f) - %s", i+1, rec.ItemID, rec.Score, rec.Reason)
	}
}

// TestGetHotItems 测试热门推荐
func TestGetHotItems(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 获取热门推荐
	recommendations, err := service.GetHotItems(ctx, "book", 10)
	if err != nil {
		t.Fatalf("获取热门推荐失败: %v", err)
	}

	if len(recommendations) != 10 {
		t.Errorf("热门推荐数量错误，期望10，实际%d", len(recommendations))
	}

	// 验证排名
	for i, rec := range recommendations {
		if rec.Rank != i+1 {
			t.Errorf("排名错误，期望%d，实际%d", i+1, rec.Rank)
		}
	}

	t.Logf("热门推荐测试通过，共%d条推荐", len(recommendations))
}

// TestRefreshRecommendations 测试刷新推荐
func TestRefreshRecommendations(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 刷新用户推荐
	err := service.RefreshRecommendations(ctx, "user123")
	if err != nil {
		t.Fatalf("刷新推荐失败: %v", err)
	}

	t.Logf("刷新推荐测试通过")
}

// TestRefreshHotItems 测试刷新热门内容
func TestRefreshHotItems(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 刷新热门内容
	err := service.RefreshHotItems(ctx, "book")
	if err != nil {
		t.Fatalf("刷新热门内容失败: %v", err)
	}

	t.Logf("刷新热门内容测试通过")
}

// TestMultipleBehaviorTypes 测试多种行为类型
func TestMultipleBehaviorTypes(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 记录不同类型的行为
	behaviorTypes := []string{"view", "click", "favorite", "read"}
	for _, actionType := range behaviorTypes {
		req := &RecordBehaviorRequest{
			UserID:     "user123",
			ItemID:     "book_001",
			ItemType:   "book",
			ActionType: actionType,
		}
		service.RecordUserBehavior(ctx, req)
	}

	// 获取用户行为
	behaviors, err := service.GetUserBehaviors(ctx, "user123", 10)
	if err != nil {
		t.Fatalf("获取用户行为失败: %v", err)
	}

	if len(behaviors) != 4 {
		t.Errorf("行为记录数量错误，期望4，实际%d", len(behaviors))
	}

	// 验证所有行为类型都被记录
	actionTypeMap := make(map[string]bool)
	for _, b := range behaviors {
		actionTypeMap[b.ActionType] = true
	}

	for _, expectedType := range behaviorTypes {
		if !actionTypeMap[expectedType] {
			t.Errorf("缺少行为类型: %s", expectedType)
		}
	}

	t.Logf("多种行为类型测试通过，记录了%d种行为类型", len(actionTypeMap))
}

// TestRecommendationScoring 测试推荐分数计算
func TestRecommendationScoring(t *testing.T) {
	mockRepo := NewMockRecommendationRepository()
	mockCache := NewMockCacheClient()
	service := NewRecommendationService(mockRepo, mockCache)

	ctx := context.Background()

	// 记录不同时长的行为
	durations := []int64{100, 500, 1000}
	for i, duration := range durations {
		req := &RecordBehaviorRequest{
			UserID:     "user123",
			ItemID:     fmt.Sprintf("book_%03d", i+1),
			ItemType:   "book",
			ActionType: "read",
			Duration:   duration,
		}
		service.RecordUserBehavior(ctx, req)
	}

	// 获取个性化推荐
	recommendations, err := service.GetPersonalizedRecommendations(ctx, "user123", 5)
	if err != nil {
		t.Fatalf("获取推荐失败: %v", err)
	}

	// 验证推荐存在
	if len(recommendations) == 0 {
		t.Error("推荐列表为空")
	}

	// 验证分数递减
	for i := 0; i < len(recommendations)-1; i++ {
		if recommendations[i].Score < recommendations[i+1].Score {
			t.Errorf("分数应该递减: [%d]=%.2f < [%d]=%.2f",
				i, recommendations[i].Score, i+1, recommendations[i+1].Score)
		}
	}

	t.Logf("推荐分数计算测试通过")
}
