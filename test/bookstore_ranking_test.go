package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockRankingRepository 模拟榜单仓储
type MockRankingRepository struct {
	mock.Mock
}

func (m *MockRankingRepository) GetByTypeWithBooks(ctx context.Context, rankingType bookstore.RankingType, period string, limit, offset int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit, offset)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateRealtimeRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateWeeklyRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateMonthlyRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) CalculateNewbieRanking(ctx context.Context, period string) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) UpdateRankings(ctx context.Context, rankingType bookstore.RankingType, period string, items []*bookstore.RankingItem) error {
	args := m.Called(ctx, rankingType, period, items)
	return args.Error(0)
}

// 实现其他必需的接口方法（简化实现）
func (m *MockRankingRepository) Create(ctx context.Context, entity *bookstore.RankingItem) (*bookstore.RankingItem, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).(*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*bookstore.RankingItem, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) Update(ctx context.Context, entity *bookstore.RankingItem) (*bookstore.RankingItem, error) {
	args := m.Called(ctx, entity)
	return args.Get(0).(*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRankingRepository) List(ctx context.Context, limit, offset int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockRankingRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRankingRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	return args.Error(0)
}

// TestGetRealtimeRanking 测试获取实时榜
func TestGetRealtimeRanking(t *testing.T) {
	// 准备测试数据
	mockRankingRepo := new(MockRankingRepository)
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	service := bookstore.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	// 模拟返回数据
	expectedItems := []*bookstore.RankingItem{
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeRealtime,
			Rank:      1,
			Score:     100.0,
			ViewCount: 1000,
			LikeCount: 50,
			Period:    "2024-01-15",
		},
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeRealtime,
			Rank:      2,
			Score:     95.0,
			ViewCount: 800,
			LikeCount: 40,
			Period:    "2024-01-15",
		},
	}

	period := bookstore.GetPeriodString(bookstore.RankingTypeRealtime, time.Now())
	mockRankingRepo.On("GetByTypeWithBooks", mock.Anything, bookstore.RankingTypeRealtime, period, 10, 0).Return(expectedItems, nil)

	// 执行测试
	ctx := context.Background()
	result, err := service.GetRealtimeRanking(ctx, 10)

	// 验证结果
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, expectedItems[0].Rank, result[0].Rank)
	assert.Equal(t, expectedItems[1].Rank, result[1].Rank)

	mockRankingRepo.AssertExpectations(t)
}

// TestGetWeeklyRanking 测试获取周榜
func TestGetWeeklyRanking(t *testing.T) {
	mockRankingRepo := new(MockRankingRepository)
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	service := bookstore.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	expectedItems := []*bookstore.RankingItem{
		{
			ID:        primitive.NewObjectID(),
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeWeekly,
			Rank:      1,
			Score:     200.0,
			ViewCount: 5000,
			LikeCount: 250,
			Period:    "2024-W03",
		},
	}

	period := "2024-W03"
	mockRankingRepo.On("GetByTypeWithBooks", mock.Anything, bookstore.RankingTypeWeekly, period, 20, 0).Return(expectedItems, nil)

	ctx := context.Background()
	result, err := service.GetWeeklyRanking(ctx, period, 20)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, bookstore.RankingTypeWeekly, result[0].Type)
	assert.Equal(t, period, result[0].Period)

	mockRankingRepo.AssertExpectations(t)
}

// TestUpdateRankings 测试更新榜单
func TestUpdateRankings(t *testing.T) {
	mockRankingRepo := new(MockRankingRepository)
	mockBookRepo := new(MockBookRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBannerRepo := new(MockBannerRepository)

	service := bookstore.NewBookstoreService(mockBookRepo, mockCategoryRepo, mockBannerRepo, mockRankingRepo)

	period := "2024-01-15"
	calculatedItems := []*bookstore.RankingItem{
		{
			BookID:    primitive.NewObjectID(),
			Type:      bookstore.RankingTypeRealtime,
			Rank:      1,
			Score:     100.0,
			ViewCount: 1000,
			LikeCount: 50,
			Period:    period,
		},
	}

	mockRankingRepo.On("CalculateRealtimeRanking", mock.Anything, period).Return(calculatedItems, nil)
	mockRankingRepo.On("UpdateRankings", mock.Anything, bookstore.RankingTypeRealtime, period, calculatedItems).Return(nil)

	ctx := context.Background()
	err := service.UpdateRankings(ctx, bookstore.RankingTypeRealtime, period)

	assert.NoError(t, err)
	mockRankingRepo.AssertExpectations(t)
}

// TestRankingTypeValidation 测试榜单类型验证
func TestRankingTypeValidation(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid realtime", "realtime", true},
		{"Valid weekly", "weekly", true},
		{"Valid monthly", "monthly", true},
		{"Valid newbie", "newbie", true},
		{"Invalid type", "invalid", false},
		{"Empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bookstore.IsValidRankingType(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetPeriodString 测试周期字符串生成
func TestGetPeriodString(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // 2024年1月15日，周一

	tests := []struct {
		name        string
		rankingType bookstore.RankingType
		expected    string
	}{
		{"Realtime period", bookstore.RankingTypeRealtime, "2024-01-15"},
		{"Weekly period", bookstore.RankingTypeWeekly, "2024-W03"},
		{"Monthly period", bookstore.RankingTypeMonthly, "2024-01"},
		{"Newbie period", bookstore.RankingTypeNewbie, "2024-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := bookstore.GetPeriodString(tt.rankingType, testTime)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestRankingItemValidation 测试榜单项验证
func TestRankingItemValidation(t *testing.T) {
	validItem := &bookstore.RankingItem{
		BookID:    primitive.NewObjectID(),
		Type:      bookstore.RankingTypeRealtime,
		Rank:      1,
		Score:     100.0,
		ViewCount: 1000,
		LikeCount: 50,
		Period:    "2024-01-15",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 这里可以添加更多的验证逻辑测试
	assert.NotNil(t, validItem.BookID)
	assert.True(t, validItem.Rank > 0)
	assert.True(t, validItem.Score >= 0)
	assert.True(t, validItem.ViewCount >= 0)
	assert.True(t, validItem.LikeCount >= 0)
	assert.NotEmpty(t, validItem.Period)
}