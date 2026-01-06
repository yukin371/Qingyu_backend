package reading

import (
	reader2 "Qingyu_backend/models/reader"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/repository/interfaces/infrastructure"
	"Qingyu_backend/service/reader"
)

// ===========================
// Mock实现 - 章节Repository
// ===========================

type MockChapterRepository struct {
	mock.Mock
}

func (m *MockChapterRepository) Create(ctx context.Context, chapter *reader2.Chapter) error {
	args := m.Called(ctx, chapter)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByID(ctx context.Context, id string) (*reader2.Chapter, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockChapterRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockChapterRepository) GetByBookID(ctx context.Context, bookID string) ([]*reader2.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByBookIDWithPagination(ctx context.Context, bookID string, limit, offset int64) ([]*reader2.Chapter, error) {
	args := m.Called(ctx, bookID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetByChapterNum(ctx context.Context, bookID string, chapterNum int) (*reader2.Chapter, error) {
	args := m.Called(ctx, bookID, chapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPrevChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader2.Chapter, error) {
	args := m.Called(ctx, bookID, currentChapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetNextChapter(ctx context.Context, bookID string, currentChapterNum int) (*reader2.Chapter, error) {
	args := m.Called(ctx, bookID, currentChapterNum)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFirstChapter(ctx context.Context, bookID string) (*reader2.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetLastChapter(ctx context.Context, bookID string) (*reader2.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetPublishedChapters(ctx context.Context, bookID string) ([]*reader2.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetVIPChapters(ctx context.Context, bookID string) ([]*reader2.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) GetFreeChapters(ctx context.Context, bookID string) ([]*reader2.Chapter, error) {
	args := m.Called(ctx, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.Chapter), args.Error(1)
}

func (m *MockChapterRepository) CountByBookID(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountByStatus(ctx context.Context, bookID string, status int) (int64, error) {
	args := m.Called(ctx, bookID, status)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) CountVIPChapters(ctx context.Context, bookID string) (int64, error) {
	args := m.Called(ctx, bookID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) BatchCreate(ctx context.Context, chapters []*reader2.Chapter) error {
	args := m.Called(ctx, chapters)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchUpdateStatus(ctx context.Context, chapterIDs []string, status int) error {
	args := m.Called(ctx, chapterIDs, status)
	return args.Error(0)
}

func (m *MockChapterRepository) BatchDelete(ctx context.Context, chapterIDs []string) error {
	args := m.Called(ctx, chapterIDs)
	return args.Error(0)
}

func (m *MockChapterRepository) CheckVIPAccess(ctx context.Context, chapterID string) (bool, error) {
	args := m.Called(ctx, chapterID)
	return args.Bool(0), args.Error(1)
}

func (m *MockChapterRepository) GetChapterPrice(ctx context.Context, chapterID string) (int64, error) {
	args := m.Called(ctx, chapterID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockChapterRepository) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	args := m.Called(ctx, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockChapterRepository) UpdateChapterContent(ctx context.Context, chapterID string, content string) error {
	args := m.Called(ctx, chapterID, content)
	return args.Error(0)
}

func (m *MockChapterRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ===========================
// Mock实现 - 缓存服务
// ===========================
// 注意：MockEventBus 已在 reader_service_test.go 中定义

type MockCacheService struct {
	mock.Mock
}

func (m *MockCacheService) GetChapterContent(ctx context.Context, chapterID string) (string, error) {
	args := m.Called(ctx, chapterID)
	return args.String(0), args.Error(1)
}

func (m *MockCacheService) SetChapterContent(ctx context.Context, chapterID string, content string, expiration time.Duration) error {
	args := m.Called(ctx, chapterID, content, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateChapterContent(ctx context.Context, chapterID string) error {
	args := m.Called(ctx, chapterID)
	return args.Error(0)
}

func (m *MockCacheService) GetChapter(ctx context.Context, chapterID string) (*reader2.Chapter, error) {
	args := m.Called(ctx, chapterID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.Chapter), args.Error(1)
}

func (m *MockCacheService) SetChapter(ctx context.Context, chapterID string, chapter *reader2.Chapter, expiration time.Duration) error {
	args := m.Called(ctx, chapterID, chapter, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateChapter(ctx context.Context, chapterID string) error {
	args := m.Called(ctx, chapterID)
	return args.Error(0)
}

func (m *MockCacheService) GetReadingSettings(ctx context.Context, userID string) (*reader2.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingSettings), args.Error(1)
}

func (m *MockCacheService) SetReadingSettings(ctx context.Context, userID string, settings *reader2.ReadingSettings, expiration time.Duration) error {
	args := m.Called(ctx, userID, settings, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateReadingSettings(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockCacheService) GetReadingProgress(ctx context.Context, userID, bookID string) (*reader2.ReadingProgress, error) {
	args := m.Called(ctx, userID, bookID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingProgress), args.Error(1)
}

func (m *MockCacheService) SetReadingProgress(ctx context.Context, userID, bookID string, progress *reader2.ReadingProgress, expiration time.Duration) error {
	args := m.Called(ctx, userID, bookID, progress, expiration)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateReadingProgress(ctx context.Context, userID, bookID string) error {
	args := m.Called(ctx, userID, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateBookChapters(ctx context.Context, bookID string) error {
	args := m.Called(ctx, bookID)
	return args.Error(0)
}

func (m *MockCacheService) InvalidateUserData(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// ===========================
// Mock实现 - VIP权限服务
// ===========================

type MockVIPPermissionService struct {
	mock.Mock
}

func (m *MockVIPPermissionService) CheckVIPAccess(ctx context.Context, userID, chapterID string, isVIPChapter bool) (bool, error) {
	args := m.Called(ctx, userID, chapterID, isVIPChapter)
	return args.Bool(0), args.Error(1)
}

func (m *MockVIPPermissionService) CheckUserVIPStatus(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVIPPermissionService) CheckChapterPurchased(ctx context.Context, userID, chapterID string) (bool, error) {
	args := m.Called(ctx, userID, chapterID)
	return args.Bool(0), args.Error(1)
}

func (m *MockVIPPermissionService) GrantVIPAccess(ctx context.Context, userID string, duration time.Duration) error {
	args := m.Called(ctx, userID, duration)
	return args.Error(0)
}

func (m *MockVIPPermissionService) GrantChapterAccess(ctx context.Context, userID, chapterID string) error {
	args := m.Called(ctx, userID, chapterID)
	return args.Error(0)
}

// ===========================
// Mock实现 - 设置Repository
// ===========================

type MockSettingsRepository struct {
	mock.Mock
}

func (m *MockSettingsRepository) Create(ctx context.Context, settings *reader2.ReadingSettings) error {
	args := m.Called(ctx, settings)
	return args.Error(0)
}

func (m *MockSettingsRepository) GetByID(ctx context.Context, id string) (*reader2.ReadingSettings, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingSettings), args.Error(1)
}

func (m *MockSettingsRepository) GetByUserID(ctx context.Context, userID string) (*reader2.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingSettings), args.Error(1)
}

func (m *MockSettingsRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockSettingsRepository) UpdateByUserID(ctx context.Context, userID string, settings *reader2.ReadingSettings) error {
	args := m.Called(ctx, userID, settings)
	return args.Error(0)
}

func (m *MockSettingsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSettingsRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSettingsRepository) CreateDefaultSettings(ctx context.Context, userID string) (*reader2.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader2.ReadingSettings), args.Error(1)
}

func (m *MockSettingsRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*reader2.ReadingSettings, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reader2.ReadingSettings), args.Error(1)
}

func (m *MockSettingsRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSettingsRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockSettingsRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// ===========================
// 测试用例 - VIP权限验证
// ===========================

// TestReaderService_GetChapterContent_VIPAccess_Success 测试VIP章节访问成功
func TestReaderService_GetChapterContent_VIPAccess_Success(t *testing.T) {
	// 创建mocks
	mockChapterRepo := new(MockChapterRepository)
	mockCacheService := new(MockCacheService)
	mockVIPService := new(MockVIPPermissionService)
	mockEventBus := new(MockEventBus)

	// 创建service
	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		mockEventBus,
		mockCacheService,
		mockVIPService,
	)

	ctx := context.Background()

	// 场景：VIP章节，用户有VIP权限
	mockCacheService.On("GetChapterContent", ctx, "vip_chapter").Return("", nil)
	mockChapterRepo.On("CheckVIPAccess", ctx, "vip_chapter").Return(true, nil)
	mockVIPService.On("CheckVIPAccess", ctx, "user123", "vip_chapter", true).Return(true, nil)
	mockChapterRepo.On("GetChapterContent", ctx, "vip_chapter").Return("VIP章节内容", nil)
	mockCacheService.On("SetChapterContent", ctx, "vip_chapter", "VIP章节内容", mock.Anything).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	content, err := service.GetChapterContent(ctx, "user123", "vip_chapter")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "VIP章节内容", content)
	mockChapterRepo.AssertExpectations(t)
	mockVIPService.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

// TestReaderService_GetChapterContent_VIPAccess_Denied 测试VIP章节访问被拒绝
func TestReaderService_GetChapterContent_VIPAccess_Denied(t *testing.T) {
	mockChapterRepo := new(MockChapterRepository)
	mockCacheService := new(MockCacheService)
	mockVIPService := new(MockVIPPermissionService)

	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		nil,
		mockCacheService,
		mockVIPService,
	)

	ctx := context.Background()

	// 场景：VIP章节，用户无权限
	mockCacheService.On("GetChapterContent", ctx, "vip_chapter").Return("", nil)
	mockChapterRepo.On("CheckVIPAccess", ctx, "vip_chapter").Return(true, nil)
	mockVIPService.On("CheckVIPAccess", ctx, "user456", "vip_chapter", true).Return(false, nil)

	// 执行测试
	_, err := service.GetChapterContent(ctx, "user456", "vip_chapter")

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "VIP章节")
	mockChapterRepo.AssertExpectations(t)
	mockVIPService.AssertExpectations(t)
}

// TestReaderService_GetChapterContent_FreeChapter 测试免费章节访问
func TestReaderService_GetChapterContent_FreeChapter(t *testing.T) {
	mockChapterRepo := new(MockChapterRepository)
	mockCacheService := new(MockCacheService)
	mockVIPService := new(MockVIPPermissionService)
	mockEventBus := new(MockEventBus)

	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		mockEventBus,
		mockCacheService,
		mockVIPService,
	)

	ctx := context.Background()

	// 场景：免费章节，任何人都可以访问
	mockCacheService.On("GetChapterContent", ctx, "free_chapter").Return("", nil)
	mockChapterRepo.On("CheckVIPAccess", ctx, "free_chapter").Return(false, nil)
	mockChapterRepo.On("GetChapterContent", ctx, "free_chapter").Return("免费章节内容", nil)
	mockCacheService.On("SetChapterContent", ctx, "free_chapter", "免费章节内容", mock.Anything).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	content, err := service.GetChapterContent(ctx, "any_user", "free_chapter")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "免费章节内容", content)
	// VIP服务不应该被调用（因为不是VIP章节）
	mockVIPService.AssertNotCalled(t, "CheckVIPAccess")
}

// ===========================
// 测试用例 - 缓存功能
// ===========================

// TestReaderService_GetChapterContent_CacheHit 测试缓存命中
func TestReaderService_GetChapterContent_CacheHit(t *testing.T) {
	mockChapterRepo := new(MockChapterRepository)
	mockCacheService := new(MockCacheService)
	mockVIPService := new(MockVIPPermissionService)
	mockEventBus := new(MockEventBus)

	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		mockEventBus,
		mockCacheService,
		mockVIPService,
	)

	ctx := context.Background()

	// 场景：缓存命中，非VIP章节
	mockCacheService.On("GetChapterContent", ctx, "chapter1").Return("缓存的内容", nil)
	mockChapterRepo.On("CheckVIPAccess", ctx, "chapter1").Return(false, nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	content, err := service.GetChapterContent(ctx, "user123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "缓存的内容", content)

	// 验证没有调用数据库获取内容
	mockChapterRepo.AssertNotCalled(t, "GetChapterContent")

	// 验证缓存被调用
	mockCacheService.AssertExpectations(t)
}

// TestReaderService_GetChapterContent_CacheMiss 测试缓存未命中
func TestReaderService_GetChapterContent_CacheMiss(t *testing.T) {
	mockChapterRepo := new(MockChapterRepository)
	mockCacheService := new(MockCacheService)
	mockVIPService := new(MockVIPPermissionService)
	mockEventBus := new(MockEventBus)

	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		mockEventBus,
		mockCacheService,
		mockVIPService,
	)

	ctx := context.Background()

	// 场景：缓存未命中
	mockCacheService.On("GetChapterContent", ctx, "chapter2").Return("", nil)
	mockChapterRepo.On("CheckVIPAccess", ctx, "chapter2").Return(false, nil)
	mockChapterRepo.On("GetChapterContent", ctx, "chapter2").Return("数据库内容", nil)
	mockCacheService.On("SetChapterContent", ctx, "chapter2", "数据库内容", mock.Anything).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	content, err := service.GetChapterContent(ctx, "user123", "chapter2")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "数据库内容", content)

	// 验证数据库被调用
	mockChapterRepo.AssertCalled(t, "GetChapterContent", ctx, "chapter2")

	// 验证内容被缓存
	mockCacheService.AssertCalled(t, "SetChapterContent", ctx, "chapter2", "数据库内容", mock.Anything)
}

// TestReaderService_GetChapterContent_CacheHit_VIPCheck 测试缓存命中但仍需验证VIP
func TestReaderService_GetChapterContent_CacheHit_VIPCheck(t *testing.T) {
	mockChapterRepo := new(MockChapterRepository)
	mockCacheService := new(MockCacheService)
	mockVIPService := new(MockVIPPermissionService)

	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		nil,
		mockCacheService,
		mockVIPService,
	)

	ctx := context.Background()

	// 场景：缓存命中，但是VIP章节，需要验证权限
	mockCacheService.On("GetChapterContent", ctx, "vip_cached").Return("缓存的VIP内容", nil)
	mockChapterRepo.On("CheckVIPAccess", ctx, "vip_cached").Return(true, nil)
	mockVIPService.On("CheckVIPAccess", ctx, "no_vip_user", "vip_cached", true).Return(false, nil)

	// 执行测试
	_, err := service.GetChapterContent(ctx, "no_vip_user", "vip_cached")

	// 验证结果：即使缓存命中，也要验证权限
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "VIP")

	mockVIPService.AssertExpectations(t)
}

// ===========================
// 测试用例 - 阅读设置缓存
// ===========================

// TestReaderService_GetReadingSettings_CacheHit 测试设置缓存命中
func TestReaderService_GetReadingSettings_CacheHit(t *testing.T) {
	mockSettingsRepo := new(MockSettingsRepository)
	mockCacheService := new(MockCacheService)

	service := reading.NewReaderService(
		nil, nil, nil,
		mockSettingsRepo,
		nil,
		mockCacheService,
		nil,
	)

	ctx := context.Background()

	// 准备缓存数据
	cachedSettings := &reader2.ReadingSettings{
		UserID:     "user123",
		FontSize:   18,
		FontFamily: "Arial",
		LineHeight: 2.0,
	}

	// 设置期望：缓存命中
	mockCacheService.On("GetReadingSettings", ctx, "user123").Return(cachedSettings, nil)

	// 执行测试
	settings, err := service.GetReadingSettings(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, 18, settings.FontSize)
	assert.Equal(t, "Arial", settings.FontFamily)

	// 验证没有查询数据库
	mockSettingsRepo.AssertNotCalled(t, "GetByUserID")

	mockCacheService.AssertExpectations(t)
}

// TestReaderService_GetReadingSettings_CacheMiss 测试设置缓存未命中
func TestReaderService_GetReadingSettings_CacheMiss(t *testing.T) {
	mockSettingsRepo := new(MockSettingsRepository)
	mockCacheService := new(MockCacheService)

	service := reading.NewReaderService(
		nil, nil, nil,
		mockSettingsRepo,
		nil,
		mockCacheService,
		nil,
	)

	ctx := context.Background()

	// 准备数据库数据
	dbSettings := &reader2.ReadingSettings{
		UserID:     "user123",
		FontSize:   16,
		FontFamily: "Microsoft YaHei",
		LineHeight: 1.8,
	}

	// 设置期望：缓存未命中，从数据库获取
	mockCacheService.On("GetReadingSettings", ctx, "user123").Return(nil, nil)
	mockSettingsRepo.On("GetByUserID", ctx, "user123").Return(dbSettings, nil)
	mockCacheService.On("SetReadingSettings", ctx, "user123", dbSettings, mock.Anything).Return(nil)

	// 执行测试
	settings, err := service.GetReadingSettings(ctx, "user123")

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, 16, settings.FontSize)

	// 验证数据库被查询
	mockSettingsRepo.AssertCalled(t, "GetByUserID", ctx, "user123")

	// 验证结果被缓存
	mockCacheService.AssertCalled(t, "SetReadingSettings", ctx, "user123", dbSettings, mock.Anything)
}

// TestReaderService_SaveReadingSettings_UpdateCache 测试保存设置并更新缓存
func TestReaderService_SaveReadingSettings_UpdateCache(t *testing.T) {
	mockSettingsRepo := new(MockSettingsRepository)
	mockCacheService := new(MockCacheService)

	service := reading.NewReaderService(
		nil, nil, nil,
		mockSettingsRepo,
		nil,
		mockCacheService,
		nil,
	)

	ctx := context.Background()

	newSettings := &reader2.ReadingSettings{
		UserID:     "user123",
		FontSize:   20,
		FontFamily: "SimSun",
		LineHeight: 2.5,
	}

	// 设置期望
	mockSettingsRepo.On("ExistsByUserID", ctx, "user123").Return(true, nil)
	mockSettingsRepo.On("UpdateByUserID", ctx, "user123", newSettings).Return(nil)
	mockCacheService.On("SetReadingSettings", ctx, "user123", newSettings, mock.Anything).Return(nil)

	// 执行测试
	err := service.SaveReadingSettings(ctx, newSettings)

	// 验证结果
	assert.NoError(t, err)

	// 验证数据库和缓存都被更新
	mockSettingsRepo.AssertExpectations(t)
	mockCacheService.AssertExpectations(t)
}

// ===========================
// 测试用例 - 边界情况
// ===========================

// TestReaderService_GetChapterContent_NilServices 测试缺少可选服务时的降级
func TestReaderService_GetChapterContent_NilServices(t *testing.T) {
	mockChapterRepo := new(MockChapterRepository)
	mockEventBus := new(MockEventBus)

	// 创建service（没有缓存和VIP服务）
	service := reading.NewReaderService(
		mockChapterRepo,
		nil, nil, nil,
		mockEventBus,
		nil, // 无缓存服务
		nil, // 无VIP服务
	)

	ctx := context.Background()

	// 场景：没有缓存和VIP服务，应该直接查询数据库
	mockChapterRepo.On("CheckVIPAccess", ctx, "chapter1").Return(false, nil)
	mockChapterRepo.On("GetChapterContent", ctx, "chapter1").Return("直接从DB获取", nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 执行测试
	content, err := service.GetChapterContent(ctx, "user123", "chapter1")

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, "直接从DB获取", content)

	mockChapterRepo.AssertExpectations(t)
}

// TestReaderService_GetReadingSettings_DefaultSettings 测试获取默认设置
func TestReaderService_GetReadingSettings_DefaultSettings(t *testing.T) {
	mockSettingsRepo := new(MockSettingsRepository)
	mockCacheService := new(MockCacheService)

	service := reading.NewReaderService(
		nil, nil, nil,
		mockSettingsRepo,
		nil,
		mockCacheService,
		nil,
	)

	ctx := context.Background()

	// 设置期望：缓存和数据库都没有数据，返回默认设置
	mockCacheService.On("GetReadingSettings", ctx, "new_user").Return(nil, nil)
	mockSettingsRepo.On("GetByUserID", ctx, "new_user").Return(nil, nil)

	// 执行测试
	settings, err := service.GetReadingSettings(ctx, "new_user")

	// 验证结果：应该返回默认设置
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, "new_user", settings.UserID)
	assert.Equal(t, 16, settings.FontSize)        // 默认值
	assert.Equal(t, "serif", settings.FontFamily) // 默认值
}
