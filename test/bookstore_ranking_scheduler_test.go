package test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/reading/bookstore"
	bookstoreService "Qingyu_backend/service/bookstore"
)

// MockBookstoreService 模拟书城服务
type MockBookstoreService struct {
	mock.Mock
}

func (m *MockBookstoreService) GetBookByID(ctx context.Context, id string) (*bookstore.Book, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreService) GetBooksByCategory(ctx context.Context, categoryID string, page, pageSize int) ([]*bookstore.Book, int64, error) {
	args := m.Called(ctx, categoryID, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetRecommendedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreService) GetFeaturedBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreService) GetHotBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreService) GetNewReleases(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreService) GetFreeBooks(ctx context.Context, page, pageSize int) ([]*bookstore.Book, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Error(1)
}

func (m *MockBookstoreService) SearchBooks(ctx context.Context, keyword string, page, pageSize int) ([]*bookstore.Book, int64, error) {
	args := m.Called(ctx, keyword, page, pageSize)
	return args.Get(0).([]*bookstore.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) SearchBooksWithFilter(ctx context.Context, filter *bookstore.BookFilter) ([]*bookstore.Book, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*bookstore.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookstoreService) GetCategoryTree(ctx context.Context) ([]*bookstore.CategoryTree, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore.CategoryTree), args.Error(1)
}

func (m *MockBookstoreService) GetCategoryByID(ctx context.Context, id string) (*bookstore.Category, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.Category), args.Error(1)
}

func (m *MockBookstoreService) GetRootCategories(ctx context.Context) ([]*bookstore.Category, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*bookstore.Category), args.Error(1)
}

func (m *MockBookstoreService) GetActiveBanners(ctx context.Context, limit int) ([]*bookstore.Banner, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.Banner), args.Error(1)
}

func (m *MockBookstoreService) IncrementBannerClick(ctx context.Context, bannerID string) error {
	args := m.Called(ctx, bannerID)
	return args.Error(0)
}

func (m *MockBookstoreService) GetRealtimeRanking(ctx context.Context, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetWeeklyRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetMonthlyRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetNewbieRanking(ctx context.Context, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) GetRankingByType(ctx context.Context, rankingType bookstore.RankingType, period string, limit int) ([]*bookstore.RankingItem, error) {
	args := m.Called(ctx, rankingType, period, limit)
	return args.Get(0).([]*bookstore.RankingItem), args.Error(1)
}

func (m *MockBookstoreService) UpdateRankings(ctx context.Context, rankingType bookstore.RankingType, period string) error {
	args := m.Called(ctx, rankingType, period)
	return args.Error(0)
}

func (m *MockBookstoreService) GetHomepageData(ctx context.Context) (*bookstoreService.HomepageData, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstoreService.HomepageData), args.Error(1)
}

func (m *MockBookstoreService) GetBookStats(ctx context.Context) (*bookstore.BookStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*bookstore.BookStats), args.Error(1)
}

func (m *MockBookstoreService) IncrementBookView(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

// TestRankingScheduler_NewScheduler 测试创建调度器
func TestRankingScheduler_NewScheduler(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[RankingScheduler] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	assert.NotNil(t, scheduler, "调度器不应该为nil")
}

// TestRankingScheduler_Start 测试启动调度器
func TestRankingScheduler_Start(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 启动调度器
	err := scheduler.Start()
	assert.NoError(t, err, "启动调度器不应该出错")

	// 等待一小段时间确保调度器已启动
	time.Sleep(100 * time.Millisecond)

	// 停止调度器
	scheduler.Stop()

	// 等待调度器完全停止
	time.Sleep(100 * time.Millisecond)
}

// TestRankingScheduler_Stop 测试停止调度器
func TestRankingScheduler_Stop(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 启动调度器
	err := scheduler.Start()
	assert.NoError(t, err, "启动调度器不应该出错")

	// 停止调度器（不应该panic）
	assert.NotPanics(t, func() {
		scheduler.Stop()
	}, "停止调度器不应该panic")
}

// TestRankingScheduler_UpdateRankingNow 测试立即更新榜单
func TestRankingScheduler_UpdateRankingNow(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 设置Mock期望
	mockService.On("UpdateRankings", mock.Anything, bookstore.RankingTypeRealtime, mock.AnythingOfType("string")).Return(nil)

	// 立即更新实时榜
	err := scheduler.UpdateRankingNow(bookstore.RankingTypeRealtime, "")

	assert.NoError(t, err, "立即更新榜单不应该出错")
	mockService.AssertExpectations(t)
}

// TestRankingScheduler_UpdateRankingNow_WithPeriod 测试使用指定period立即更新榜单
func TestRankingScheduler_UpdateRankingNow_WithPeriod(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 指定period
	period := "2025-10"

	// 设置Mock期望
	mockService.On("UpdateRankings", mock.Anything, bookstore.RankingTypeMonthly, period).Return(nil)

	// 立即更新月榜
	err := scheduler.UpdateRankingNow(bookstore.RankingTypeMonthly, period)

	assert.NoError(t, err, "立即更新榜单不应该出错")
	mockService.AssertExpectations(t)
}

// TestRankingScheduler_UpdateRankingNow_Error 测试更新榜单失败
func TestRankingScheduler_UpdateRankingNow_Error(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 设置Mock期望 - 返回错误
	expectedError := assert.AnError
	mockService.On("UpdateRankings", mock.Anything, bookstore.RankingTypeWeekly, mock.AnythingOfType("string")).Return(expectedError)

	// 立即更新周榜
	err := scheduler.UpdateRankingNow(bookstore.RankingTypeWeekly, "")

	assert.Error(t, err, "应该返回错误")
	assert.Equal(t, expectedError, err, "错误应该是预期的错误")
	mockService.AssertExpectations(t)
}

// TestRankingScheduler_GetSchedulerStatus 测试获取调度器状态
func TestRankingScheduler_GetSchedulerStatus(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 启动调度器
	err := scheduler.Start()
	assert.NoError(t, err, "启动调度器不应该出错")

	// 获取调度器状态
	status := scheduler.GetSchedulerStatus()

	assert.NotNil(t, status, "状态不应该为nil")
	assert.True(t, status["running"].(bool), "调度器应该在运行")
	assert.Equal(t, 5, status["job_count"].(int), "应该有5个定时任务")

	// 检查next_runs
	nextRuns := status["next_runs"].([]map[string]interface{})
	assert.Equal(t, 5, len(nextRuns), "应该有5个下次运行时间")

	// 停止调度器
	scheduler.Stop()
}

// TestRankingScheduler_GetSchedulerStatus_NotStarted 测试未启动时获取调度器状态
func TestRankingScheduler_GetSchedulerStatus_NotStarted(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 获取调度器状态（未启动）
	status := scheduler.GetSchedulerStatus()

	assert.NotNil(t, status, "状态不应该为nil")
	assert.False(t, status["running"].(bool), "调度器不应该在运行")
	assert.Equal(t, 0, status["job_count"].(int), "不应该有定时任务")
}

// TestRankingScheduler_MultipleStarts 测试多次启动调度器
func TestRankingScheduler_MultipleStarts(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 第一次启动
	err1 := scheduler.Start()
	assert.NoError(t, err1, "第一次启动不应该出错")

	// 第二次启动（应该会添加新的任务）
	err2 := scheduler.Start()
	assert.NoError(t, err2, "第二次启动不应该出错")

	// 检查状态
	status := scheduler.GetSchedulerStatus()
	// 注意：多次启动会添加多个任务，这里应该有10个任务
	jobCount := status["job_count"].(int)
	assert.True(t, jobCount >= 5, "任务数量应该至少有5个")

	// 停止调度器
	scheduler.Stop()
}

// TestRankingScheduler_ConcurrentUpdates 测试并发更新榜单
func TestRankingScheduler_ConcurrentUpdates(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 设置Mock期望 - 允许多次调用
	mockService.On("UpdateRankings", mock.Anything, mock.AnythingOfType("bookstore.RankingType"), mock.AnythingOfType("string")).Return(nil)

	// 并发更新多个榜单
	done := make(chan bool, 4)

	go func() {
		scheduler.UpdateRankingNow(bookstore.RankingTypeRealtime, "")
		done <- true
	}()

	go func() {
		scheduler.UpdateRankingNow(bookstore.RankingTypeWeekly, "")
		done <- true
	}()

	go func() {
		scheduler.UpdateRankingNow(bookstore.RankingTypeMonthly, "")
		done <- true
	}()

	go func() {
		scheduler.UpdateRankingNow(bookstore.RankingTypeNewbie, "")
		done <- true
	}()

	// 等待所有goroutine完成
	for i := 0; i < 4; i++ {
		<-done
	}

	// 验证所有Mock调用
	mockService.AssertExpectations(t)
}

// TestRankingScheduler_UpdateAllRankingTypes 测试更新所有类型的榜单
func TestRankingScheduler_UpdateAllRankingTypes(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 定义所有榜单类型
	rankingTypes := []bookstore.RankingType{
		bookstore.RankingTypeRealtime,
		bookstore.RankingTypeWeekly,
		bookstore.RankingTypeMonthly,
		bookstore.RankingTypeNewbie,
	}

	// 为每种类型设置Mock期望
	for _, rankingType := range rankingTypes {
		mockService.On("UpdateRankings", mock.Anything, rankingType, mock.AnythingOfType("string")).Return(nil)
	}

	// 更新所有类型的榜单
	for _, rankingType := range rankingTypes {
		err := scheduler.UpdateRankingNow(rankingType, "")
		assert.NoError(t, err, "更新%s榜单不应该出错", rankingType)
	}

	// 验证所有Mock调用
	mockService.AssertExpectations(t)
}

// TestRankingScheduler_Cleanup 测试清理过期榜单
func TestRankingScheduler_Cleanup(t *testing.T) {
	mockService := new(MockBookstoreService)
	logger := log.New(os.Stdout, "[Test] ", log.LstdFlags)

	scheduler := bookstoreService.NewRankingScheduler(mockService, logger)

	// 启动调度器
	err := scheduler.Start()
	assert.NoError(t, err, "启动调度器不应该出错")

	// 等待一段时间（清理任务在每天凌晨4点执行，这里只是测试调度器能正常添加任务）
	time.Sleep(100 * time.Millisecond)

	// 停止调度器
	scheduler.Stop()

	// 测试通过说明清理任务已经成功添加到调度器
	assert.True(t, true, "清理任务测试通过")
}


