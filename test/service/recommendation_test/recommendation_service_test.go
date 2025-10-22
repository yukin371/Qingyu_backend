package recommendation_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	reco "Qingyu_backend/models/recommendation/reco"
	"Qingyu_backend/service/recommendation"
)

// ===========================
// Mock实现
// 注意: MockItemFeatureRepo和MockHotRecommendationRepo在recommendation_service_enhanced_test.go中定义
// ===========================

type MockBehaviorRepo struct {
	mock.Mock
}

func (m *MockBehaviorRepo) Create(ctx context.Context, b *reco.Behavior) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

func (m *MockBehaviorRepo) BatchCreate(ctx context.Context, bs []*reco.Behavior) error {
	args := m.Called(ctx, bs)
	return args.Error(0)
}

func (m *MockBehaviorRepo) GetByUser(ctx context.Context, userID string, limit int) ([]*reco.Behavior, error) {
	args := m.Called(ctx, userID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reco.Behavior), args.Error(1)
}

type MockProfileRepo struct {
	mock.Mock
}

func (m *MockProfileRepo) Upsert(ctx context.Context, p *reco.UserProfile) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProfileRepo) GetByUserID(ctx context.Context, userID string) (*reco.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reco.UserProfile), args.Error(1)
}

// ===========================
// 测试用例
// ===========================

func TestRecommendationService_RecordBehavior(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	t.Run("成功记录行为", func(t *testing.T) {
		ctx := context.Background()
		behavior := &reco.Behavior{
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
		behavior := &reco.Behavior{
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

	service := recommendation.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	t.Run("获取推荐-用户有画像", func(t *testing.T) {
		ctx := context.Background()
		profile := &reco.UserProfile{
			UserID: "user123",
			Tags: map[string]float64{
				"玄幻": 0.8,
				"修仙": 0.6,
			},
		}

		mockProfileRepo.On("GetByUserID", ctx, "user123").Return(profile, nil)

		recommendations, err := service.GetPersonalizedRecommendations(ctx, "user123", 10)

		assert.NoError(t, err)
		assert.NotNil(t, recommendations)
		mockProfileRepo.AssertExpectations(t)
	})

	t.Run("获取推荐-用户无画像", func(t *testing.T) {
		ctx := context.Background()

		mockProfileRepo.On("GetByUserID", ctx, "new_user").Return(nil, nil)

		recommendations, err := service.GetPersonalizedRecommendations(ctx, "new_user", 10)

		assert.NoError(t, err)
		assert.NotNil(t, recommendations)
		assert.Empty(t, recommendations) // 无画像返回空
		mockProfileRepo.AssertExpectations(t)
	})
}

func TestRecommendationService_GetSimilarItems(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	t.Run("获取相似物品", func(t *testing.T) {
		ctx := context.Background()

		similarItems, err := service.GetSimilarItems(ctx, "book123", 10)

		assert.NoError(t, err)
		assert.NotNil(t, similarItems)
		// 当前是占位实现，返回空列表
		assert.Empty(t, similarItems)
	})
}

func TestRecommendationService_ServiceInfo(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(mockBehaviorRepo, mockProfileRepo, mockItemFeatureRepo, mockHotRepo)

	name, version := service.ServiceInfo()

	assert.Equal(t, "RecommendationService", name)
	assert.Equal(t, "1.0.0", version)
}
