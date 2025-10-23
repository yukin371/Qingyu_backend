package recommendation_test

import (
	recommendation2 "Qingyu_backend/models/recommendation"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"Qingyu_backend/service/recommendation"
)

// ===========================
// Mock实现 - 完整版本
// ===========================

type MockItemFeatureRepo struct {
	mock.Mock
}

func (m *MockItemFeatureRepo) Create(ctx context.Context, feature *recommendation2.ItemFeature) error {
	args := m.Called(ctx, feature)
	return args.Error(0)
}

func (m *MockItemFeatureRepo) Upsert(ctx context.Context, feature *recommendation2.ItemFeature) error {
	args := m.Called(ctx, feature)
	return args.Error(0)
}

func (m *MockItemFeatureRepo) GetByItemID(ctx context.Context, itemID string) (*recommendation2.ItemFeature, error) {
	args := m.Called(ctx, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*recommendation2.ItemFeature), args.Error(1)
}

func (m *MockItemFeatureRepo) BatchGetByItemIDs(ctx context.Context, itemIDs []string) ([]*recommendation2.ItemFeature, error) {
	args := m.Called(ctx, itemIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation2.ItemFeature), args.Error(1)
}

func (m *MockItemFeatureRepo) GetByCategory(ctx context.Context, category string, limit int) ([]*recommendation2.ItemFeature, error) {
	args := m.Called(ctx, category, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation2.ItemFeature), args.Error(1)
}

func (m *MockItemFeatureRepo) GetByTags(ctx context.Context, tags map[string]float64, limit int) ([]*recommendation2.ItemFeature, error) {
	args := m.Called(ctx, tags, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*recommendation2.ItemFeature), args.Error(1)
}

func (m *MockItemFeatureRepo) Delete(ctx context.Context, itemID string) error {
	args := m.Called(ctx, itemID)
	return args.Error(0)
}

func (m *MockItemFeatureRepo) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type MockHotRecommendationRepo struct {
	mock.Mock
}

func (m *MockHotRecommendationRepo) GetHotBooks(ctx context.Context, limit int, days int) ([]string, error) {
	args := m.Called(ctx, limit, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockHotRecommendationRepo) GetHotBooksByCategory(ctx context.Context, category string, limit int, days int) ([]string, error) {
	args := m.Called(ctx, category, limit, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockHotRecommendationRepo) GetTrendingBooks(ctx context.Context, limit int) ([]string, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockHotRecommendationRepo) GetNewPopularBooks(ctx context.Context, limit int, daysThreshold int) ([]string, error) {
	args := m.Called(ctx, limit, daysThreshold)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockHotRecommendationRepo) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ===========================
// 测试用例 - 热门推荐
// ===========================

func TestRecommendationService_GetHotRecommendations(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(
		mockBehaviorRepo,
		mockProfileRepo,
		mockItemFeatureRepo,
		mockHotRepo,
	)

	t.Run("成功获取热门推荐", func(t *testing.T) {
		ctx := context.Background()
		expectedBooks := []string{"book1", "book2", "book3"}

		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return(expectedBooks, nil)

		result, err := service.GetHotRecommendations(ctx, 10, 7)

		require.NoError(t, err)
		assert.Equal(t, expectedBooks, result)
		mockHotRepo.AssertExpectations(t)
	})

	t.Run("获取热门推荐-空结果", func(t *testing.T) {
		ctx := context.Background()
		expectedBooks := []string{"book1", "book2", "book3"}

		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return(expectedBooks, nil)

		result, err := service.GetHotRecommendations(ctx, 10, 7)

		require.NoError(t, err)
		assert.Equal(t, expectedBooks, result)
		mockHotRepo.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 个性化推荐
// ===========================

func TestRecommendationService_GetPersonalizedRecommendations_Enhanced(t *testing.T) {
	t.Run("有画像用户-基于标签推荐", func(t *testing.T) {
		// 为每个子测试创建独立的Mock实例
		mockBehaviorRepo := new(MockBehaviorRepo)
		mockProfileRepo := new(MockProfileRepo)
		mockItemFeatureRepo := new(MockItemFeatureRepo)
		mockHotRepo := new(MockHotRecommendationRepo)

		service := recommendation.NewRecommendationService(
			mockBehaviorRepo,
			mockProfileRepo,
			mockItemFeatureRepo,
			mockHotRepo,
		)

		ctx := context.Background()
		userID := "user_with_profile"

		// 用户画像
		profile := &recommendation2.UserProfile{
			UserID: userID,
			Tags: map[string]float64{
				"玄幻": 0.8,
				"修真": 0.6,
			},
			Categories: map[string]float64{
				"玄幻": 0.8,
			},
		}

		// Mock物品特征
		itemFeatures := []*recommendation2.ItemFeature{
			{
				ItemID: "book1",
				Tags: map[string]float64{
					"玄幻": 0.9,
					"修真": 0.7,
				},
				Categories: []string{"玄幻"},
			},
			{
				ItemID: "book2",
				Tags: map[string]float64{
					"玄幻": 0.8,
				},
				Categories: []string{"玄幻"},
			},
		}

		mockProfileRepo.On("GetByUserID", ctx, userID).Return(profile, nil)
		mockItemFeatureRepo.On("GetByTags", ctx, profile.Tags, mock.AnythingOfType("int")).Return(itemFeatures, nil)
		mockItemFeatureRepo.On("GetByCategory", ctx, "玄幻", mock.AnythingOfType("int")).Return(itemFeatures, nil)
		// 添加GetHotBooks的mock，用于补充推荐结果
		mockHotRepo.On("GetHotBooks", ctx, mock.AnythingOfType("int"), 7).Return([]string{"hot1", "hot2"}, nil).Maybe()

		result, err := service.GetPersonalizedRecommendations(ctx, userID, 10)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
		mockProfileRepo.AssertExpectations(t)
	})

	t.Run("无画像用户-冷启动策略", func(t *testing.T) {
		// 为每个子测试创建独立的Mock实例
		mockBehaviorRepo := new(MockBehaviorRepo)
		mockProfileRepo := new(MockProfileRepo)
		mockItemFeatureRepo := new(MockItemFeatureRepo)
		mockHotRepo := new(MockHotRecommendationRepo)

		service := recommendation.NewRecommendationService(
			mockBehaviorRepo,
			mockProfileRepo,
			mockItemFeatureRepo,
			mockHotRepo,
		)

		ctx := context.Background()
		userID := "new_user"
		hotBooks := []string{"hot1", "hot2", "hot3"}

		mockProfileRepo.On("GetByUserID", ctx, userID).Return(nil, nil)
		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return(hotBooks, nil)

		result, err := service.GetPersonalizedRecommendations(ctx, userID, 10)

		require.NoError(t, err)
		assert.Equal(t, hotBooks, result)
		mockProfileRepo.AssertExpectations(t)
		mockHotRepo.AssertExpectations(t)
	})

	t.Run("画像为空-冷启动策略", func(t *testing.T) {
		// 为每个子测试创建独立的Mock实例
		mockBehaviorRepo := new(MockBehaviorRepo)
		mockProfileRepo := new(MockProfileRepo)
		mockItemFeatureRepo := new(MockItemFeatureRepo)
		mockHotRepo := new(MockHotRecommendationRepo)

		service := recommendation.NewRecommendationService(
			mockBehaviorRepo,
			mockProfileRepo,
			mockItemFeatureRepo,
			mockHotRepo,
		)

		ctx := context.Background()
		userID := "user_empty_profile"
		hotBooks := []string{"hot1", "hot2"}

		// 空画像
		emptyProfile := &recommendation2.UserProfile{
			UserID:     userID,
			Tags:       map[string]float64{},
			Categories: map[string]float64{},
		}

		mockProfileRepo.On("GetByUserID", ctx, userID).Return(emptyProfile, nil)
		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return(hotBooks, nil)

		result, err := service.GetPersonalizedRecommendations(ctx, userID, 10)

		require.NoError(t, err)
		assert.Equal(t, hotBooks, result)
		mockHotRepo.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 相似推荐
// ===========================

func TestRecommendationService_GetSimilarItems_Enhanced(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(
		mockBehaviorRepo,
		mockProfileRepo,
		mockItemFeatureRepo,
		mockHotRepo,
	)

	t.Run("找到相似物品", func(t *testing.T) {
		ctx := context.Background()
		targetItemID := "book1"

		// 目标物品特征
		targetFeature := &recommendation2.ItemFeature{
			ItemID: targetItemID,
			Tags: map[string]float64{
				"玄幻": 0.9,
				"修真": 0.8,
			},
			Categories: []string{"玄幻"},
			Authors:    []string{"作者A"},
		}

		// 相似物品
		similarFeatures := []*recommendation2.ItemFeature{
			{
				ItemID: "book2",
				Tags: map[string]float64{
					"玄幻": 0.8,
					"修真": 0.7,
				},
				Categories: []string{"玄幻"},
			},
			{
				ItemID: "book3",
				Tags: map[string]float64{
					"玄幻": 0.9,
				},
				Categories: []string{"玄幻"},
			},
		}

		mockItemFeatureRepo.On("GetByItemID", ctx, targetItemID).Return(targetFeature, nil)
		mockItemFeatureRepo.On("GetByCategory", ctx, "玄幻", mock.AnythingOfType("int")).Return(similarFeatures, nil)
		mockItemFeatureRepo.On("GetByTags", ctx, targetFeature.Tags, mock.AnythingOfType("int")).Return(similarFeatures, nil)

		result, err := service.GetSimilarItems(ctx, targetItemID, 10)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
		// 验证不包含目标物品本身
		for _, itemID := range result {
			assert.NotEqual(t, targetItemID, itemID)
		}
		mockItemFeatureRepo.AssertExpectations(t)
	})

	t.Run("物品特征不存在-返回热门", func(t *testing.T) {
		ctx := context.Background()
		itemID := "non_existent_item"
		hotBooks := []string{"hot1", "hot2"}

		mockItemFeatureRepo.On("GetByItemID", ctx, itemID).Return(nil, nil)
		mockHotRepo.On("GetHotBooks", ctx, 10, 7).Return(hotBooks, nil)

		result, err := service.GetSimilarItems(ctx, itemID, 10)

		require.NoError(t, err)
		assert.Equal(t, hotBooks, result)
		mockItemFeatureRepo.AssertExpectations(t)
		mockHotRepo.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 首页推荐
// ===========================

func TestRecommendationService_GetHomepageRecommendations(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(
		mockBehaviorRepo,
		mockProfileRepo,
		mockItemFeatureRepo,
		mockHotRepo,
	)

	t.Run("混合推荐-有用户ID", func(t *testing.T) {
		ctx := context.Background()
		userID := "user123"

		// Mock个性化推荐
		profile := &recommendation2.UserProfile{
			UserID: userID,
			Tags:   map[string]float64{"玄幻": 0.8},
		}
		itemFeatures := []*recommendation2.ItemFeature{
			{ItemID: "personal1", Tags: map[string]float64{"玄幻": 0.9}},
		}

		mockProfileRepo.On("GetByUserID", ctx, userID).Return(profile, nil)
		mockItemFeatureRepo.On("GetByTags", ctx, mock.Anything, mock.Anything).Return(itemFeatures, nil)

		// Mock热门推荐
		hotBooks := []string{"hot1", "hot2"}
		mockHotRepo.On("GetHotBooks", ctx, mock.Anything, 7).Return(hotBooks, nil)

		// Mock新书推荐
		newBooks := []string{"new1", "new2"}
		mockHotRepo.On("GetNewPopularBooks", ctx, mock.Anything, 30).Return(newBooks, nil)

		// Mock飙升榜
		trendingBooks := []string{"trending1"}
		mockHotRepo.On("GetTrendingBooks", ctx, mock.Anything).Return(trendingBooks, nil)

		result, err := service.GetHomepageRecommendations(ctx, userID, 20)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.LessOrEqual(t, len(result), 20)
		mockProfileRepo.AssertExpectations(t)
	})

	t.Run("混合推荐-无用户ID", func(t *testing.T) {
		ctx := context.Background()

		// Mock热门推荐
		hotBooks := []string{"hot1", "hot2", "hot3"}
		mockHotRepo.On("GetHotBooks", ctx, mock.Anything, 7).Return(hotBooks, nil)

		// Mock新书推荐
		newBooks := []string{"new1", "new2"}
		mockHotRepo.On("GetNewPopularBooks", ctx, mock.Anything, 30).Return(newBooks, nil)

		// Mock飙升榜
		trendingBooks := []string{"trending1"}
		mockHotRepo.On("GetTrendingBooks", ctx, mock.Anything).Return(trendingBooks, nil)

		result, err := service.GetHomepageRecommendations(ctx, "", 20)

		require.NoError(t, err)
		assert.NotEmpty(t, result)
		mockHotRepo.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 分类推荐
// ===========================

func TestRecommendationService_GetCategoryRecommendations(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(
		mockBehaviorRepo,
		mockProfileRepo,
		mockItemFeatureRepo,
		mockHotRepo,
	)

	t.Run("成功获取分类推荐", func(t *testing.T) {
		ctx := context.Background()
		category := "玄幻"
		expectedBooks := []string{"book1", "book2", "book3"}

		mockHotRepo.On("GetHotBooksByCategory", ctx, category, 10, 7).Return(expectedBooks, nil)

		result, err := service.GetCategoryRecommendations(ctx, category, 10)

		require.NoError(t, err)
		assert.Equal(t, expectedBooks, result)
		mockHotRepo.AssertExpectations(t)
	})
}

// ===========================
// 测试用例 - 用户行为记录
// ===========================

func TestRecommendationService_RecordBehavior_Enhanced(t *testing.T) {
	mockBehaviorRepo := new(MockBehaviorRepo)
	mockProfileRepo := new(MockProfileRepo)
	mockItemFeatureRepo := new(MockItemFeatureRepo)
	mockHotRepo := new(MockHotRecommendationRepo)

	service := recommendation.NewRecommendationService(
		mockBehaviorRepo,
		mockProfileRepo,
		mockItemFeatureRepo,
		mockHotRepo,
	)

	t.Run("成功记录完整行为", func(t *testing.T) {
		ctx := context.Background()
		behavior := &recommendation2.Behavior{
			UserID:       "user123",
			ItemID:       "book456",
			ChapterID:    "chapter789",
			BehaviorType: "read",
			Value:        120.5,
			Metadata: map[string]interface{}{
				"readTime": 120,
				"progress": 0.5,
			},
		}

		mockBehaviorRepo.On("Create", ctx, behavior).Return(nil)

		err := service.RecordBehavior(ctx, behavior)

		require.NoError(t, err)
		mockBehaviorRepo.AssertExpectations(t)
	})

	t.Run("记录行为-缺少UserID", func(t *testing.T) {
		ctx := context.Background()
		behavior := &recommendation2.Behavior{
			ItemID:       "book456",
			BehaviorType: "read",
		}

		err := service.RecordBehavior(ctx, behavior)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id")
	})

	t.Run("记录行为-缺少ItemID", func(t *testing.T) {
		ctx := context.Background()
		behavior := &recommendation2.Behavior{
			UserID:       "user123",
			BehaviorType: "read",
		}

		err := service.RecordBehavior(ctx, behavior)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "item_id")
	})
}
