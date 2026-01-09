package recommendation_test

import (
	recommendation2 "Qingyu_backend/models/recommendation"
	recommendationService "Qingyu_backend/service/recommendation"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ===========================
// Mock实现
// 注意: MockItemFeatureRepo和MockHotRecommendationRepo在recommendation_service_enhanced_test.go中定义
// ===========================

type MockBehaviorRepo struct {
	mock.Mock
}

func (m *MockBehaviorRepo) Create(ctx context.Context, b *recommendation2.Behavior) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBehaviorRepo) BatchCreate(ctx context.Context, bs []*recommendation2.Behavior) error {
	args := m.Called(ctx, bs)
	return args.Error(0)
}

func (m *MockBehaviorRepo) GetByUser(ctx context.Context, userID string, limit int) ([]*recommendation2.Behavior, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation2.Behavior), args.Error(1)
}

type MockProfileRepo struct {
	mock.Mock
}

func (m *MockProfileRepo) Upsert(ctx context.Context, p *recommendation2.UserProfile) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProfileRepo) GetByUserID(ctx context.Context, userID string) (*recommendation2.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*recommendation2.UserProfile), args.Error(1)
}

// ===========================
// 测试用例
// ===========================

func TestRecommendationService_RecordBehavior(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendationService.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	t.Run("成功记录行为", func(t *testing.T) {
		ctx := context.Background()
		behavior := &recommendation2.Behavior{
			UserID:       "user123",
			ItemID:       "book456",
			BehaviorType: "view",
			Value:        1.0,
			OccurredAt:   time.Now(),
		}

		mockBehaviorRepo.On("Create", ctx, behavior).Return(nil)

		err := service.RecordBehavior(ctx, behavior)

		assert.NoError(t, err)
		mockBehaviorRepo.AssertExpectations(t)
	})

	t.Run("缺少必需字段", func(t *testing.T) {
		ctx := context.Background()
		behavior := &recommendation2.Behavior{
			BehaviorType: "view",
		}

		err := service.RecordBehavior(ctx, behavior)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "required")
	})
}

func TestRecommendationService_GetPersonalizedRecommendations(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendationService.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	t.Run("获取推荐-用户有画像", func(t *testing.T) {
		ctx := context.Background()
		profile := &recommendation2.UserProfile{
			UserID: "user123",
			Tags: map[string]float64{
				"玄幻": 0.8,
				"修仙": 0.6,
			},
		}

		mockProfileRepo.On("GetByUserID", ctx, "user123").Return(profile, nil)

		// Mock GetByTags返回推荐的物品特征
		mockItems := []*recommendation2.ItemFeature{
			{ItemID: "item1", Tags: map[string]float64{"玄幻": 0.9}},
			{ItemID: "item2", Tags: map[string]float64{"修仙": 0.7}},
		}
		mockItemFeatureRepo.On("GetByTags", ctx, profile.Tags, 20).Return(mockItems, nil)

		// Mock GetHotBooks用于补充推荐结果
		mockHotRecommendations := []string{"hot1", "hot2", "hot3"}
		mockHotRepo.On("GetHotBooks", ctx, 8, 7).Return(mockHotRecommendations, nil)

		recommendations, err := service.GetPersonalizedRecommendations(ctx, "user123", 10)

		assert.NoError(t, err)
		assert.NotNil(t, recommendations)
		mockProfileRepo.AssertExpectations(t)
		mockItemFeatureRepo.AssertExpectations(t)
	})

	t.Run("获取推荐-用户无画像", func(t *testing.T) {
		ctx := context.Background()

		mockProfileRepo.On("GetByUserID", ctx, "new_user").Return(nil, nil)

		// 无画像时返回热门推荐
		mockHotRecommendations := []string{"hot1", "hot2", "hot3"}
		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return(mockHotRecommendations, nil)

		recommendations, err := service.GetPersonalizedRecommendations(ctx, "new_user", 10)

		assert.NoError(t, err)
		assert.NotNil(t, recommendations)
		assert.NotEmpty(t, recommendations) // 无画像返回热门推荐
		mockProfileRepo.AssertExpectations(t)
		mockHotRepo.AssertExpectations(t)
	})
}

func TestRecommendationService_GetSimilarItems(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendationService.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	t.Run("获取相似物品", func(t *testing.T) {
		ctx := context.Background()

		// Mock GetByItemID返回nil（物品特征不存在）
		mockItemFeatureRepo.On("GetByItemID", ctx, "book123").Return(nil, nil)
		// Mock热门推荐作为fallback
		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return([]string{}, nil)

		similarItems, err := service.GetSimilarItems(ctx, "book123", 10)

		assert.NoError(t, err)
		assert.NotNil(t, similarItems)
		// 当前是占位实现，返回空列表
		assert.Empty(t, similarItems)

		mockItemFeatureRepo.AssertExpectations(t)
		mockHotRepo.AssertExpectations(t)
	})
}

func TestRecommendationService_ServiceInfo(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendationService.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	name, version := service.ServiceInfo()

	assert.Equal(t, "RecommendationService", name)
	assert.Equal(t, "1.0.0", version)
}
