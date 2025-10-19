package reading

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"Qingyu_backend/models/reading/reader"
	readingRepo "Qingyu_backend/repository/interfaces/reading"
	"Qingyu_backend/service/reading"
)

// MockReadingSettingsRepository Mock实现
type MockReadingSettingsRepository struct {
	mock.Mock
	readingRepo.ReadingSettingsRepository // 嵌入接口避免实现所有方法
}

func (m *MockReadingSettingsRepository) GetByUserID(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReadingSettingsRepository) UpdateByUserID(ctx context.Context, userID string, settings *reader.ReadingSettings) error {
	args := m.Called(ctx, userID, settings)
	return args.Error(0)
}

func (m *MockReadingSettingsRepository) CreateDefaultSettings(ctx context.Context, userID string) (*reader.ReadingSettings, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reader.ReadingSettings), args.Error(1)
}

func (m *MockReadingSettingsRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockReadingSettingsRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// 测试辅助函数
func createTestSettings(userID string) *reader.ReadingSettings {
	return &reader.ReadingSettings{
		ID:          "settings123",
		UserID:      userID,
		FontFamily:  "Arial",
		FontSize:    16,
		LineHeight:  1.5,
		Theme:       "light",
		Background:  "#FFFFFF",
		PageMode:    1,
		AutoScroll:  false,
		ScrollSpeed: 50,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// TestSettingService_GetSetting_Success 测试成功获取现有设置
func TestSettingService_GetSetting_Success(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	expectedSettings := createTestSettings(userID)

	// 设置Mock期望：返回现有设置
	mockRepo.On("GetByUserID", ctx, userID).Return(expectedSettings, nil)

	// 执行测试
	settings, err := service.GetSetting(ctx, userID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	assert.Equal(t, "Arial", settings.FontFamily)
	mockRepo.AssertExpectations(t)
}

// TestSettingService_GetSetting_CreateDefault 测试不存在时创建默认设置
func TestSettingService_GetSetting_CreateDefault(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user456"
	defaultSettings := createTestSettings(userID)

	// 设置Mock期望：首次查询返回nil，然后创建默认设置
	mockRepo.On("GetByUserID", ctx, userID).Return(nil, nil)
	mockRepo.On("CreateDefaultSettings", ctx, userID).Return(defaultSettings, nil)

	// 执行测试
	settings, err := service.GetSetting(ctx, userID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	mockRepo.AssertExpectations(t)
}

// TestSettingService_GetSetting_EmptyUserID 测试空userID参数验证
func TestSettingService_GetSetting_EmptyUserID(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()

	// 执行测试
	settings, err := service.GetSetting(ctx, "")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, settings)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// TestSettingService_GetSetting_RepositoryError 测试Repository返回错误
func TestSettingService_GetSetting_RepositoryError(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user789"

	// 设置Mock期望：Repository返回错误
	mockRepo.On("GetByUserID", ctx, userID).Return(nil, errors.New("database error"))

	// 执行测试
	settings, err := service.GetSetting(ctx, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, settings)
	mockRepo.AssertExpectations(t)
}

// TestSettingService_UpdateSetting_Success 测试成功更新设置
func TestSettingService_UpdateSetting_Success(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	settings := createTestSettings(userID)
	settings.FontSize = 18 // 修改字号

	// 设置Mock期望
	mockRepo.On("UpdateByUserID", ctx, userID, settings).Return(nil)

	// 执行测试
	err := service.UpdateSetting(ctx, userID, settings)

	// 验证结果
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestSettingService_UpdateSetting_EmptyUserID 测试空userID参数验证
func TestSettingService_UpdateSetting_EmptyUserID(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	settings := createTestSettings("user123")

	// 执行测试
	err := service.UpdateSetting(ctx, "", settings)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// TestSettingService_UpdateSetting_NilSettings 测试空settings参数验证
func TestSettingService_UpdateSetting_NilSettings(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()

	// 执行测试
	err := service.UpdateSetting(ctx, "user123", nil)

	// 验证结果
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "设置信息不能为空")
}

// TestSettingService_UpdateSetting_EnsureUserID 测试确保userID正确设置
func TestSettingService_UpdateSetting_EnsureUserID(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	settings := createTestSettings("wrongUser") // 错误的userID

	// 设置Mock期望，验证settings的UserID被更正
	mockRepo.On("UpdateByUserID", ctx, userID, mock.MatchedBy(func(s *reader.ReadingSettings) bool {
		return s.UserID == userID // 验证UserID被设置为正确的值
	})).Return(nil)

	// 执行测试
	err := service.UpdateSetting(ctx, userID, settings)

	// 验证结果
	assert.NoError(t, err)
	assert.Equal(t, userID, settings.UserID) // 确认UserID被更正
	mockRepo.AssertExpectations(t)
}

// TestSettingService_ResetSetting_FirstTime 测试首次重置（创建默认设置）
func TestSettingService_ResetSetting_FirstTime(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	defaultSettings := createTestSettings(userID)

	// 设置Mock期望：用户不存在设置
	mockRepo.On("ExistsByUserID", ctx, userID).Return(false, nil)
	mockRepo.On("CreateDefaultSettings", ctx, userID).Return(defaultSettings, nil)

	// 执行测试
	settings, err := service.ResetSetting(ctx, userID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	mockRepo.AssertExpectations(t)
}

// TestSettingService_ResetSetting_DeleteAndRecreate 测试删除现有设置并重建
func TestSettingService_ResetSetting_DeleteAndRecreate(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user456"
	existingSettings := createTestSettings(userID)
	newDefaultSettings := createTestSettings(userID)

	// 设置Mock期望：存在设置，需要删除后重建
	mockRepo.On("ExistsByUserID", ctx, userID).Return(true, nil)
	mockRepo.On("GetByUserID", ctx, userID).Return(existingSettings, nil)
	mockRepo.On("Delete", ctx, existingSettings.ID).Return(nil)
	mockRepo.On("CreateDefaultSettings", ctx, userID).Return(newDefaultSettings, nil)

	// 执行测试
	settings, err := service.ResetSetting(ctx, userID)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, userID, settings.UserID)
	mockRepo.AssertExpectations(t)
}

// TestSettingService_ResetSetting_EmptyUserID 测试空userID参数验证
func TestSettingService_ResetSetting_EmptyUserID(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()

	// 执行测试
	settings, err := service.ResetSetting(ctx, "")

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, settings)
	assert.Contains(t, err.Error(), "用户ID不能为空")
}

// TestSettingService_ResetSetting_ExistsCheckError 测试ExistsByUserID返回错误
func TestSettingService_ResetSetting_ExistsCheckError(t *testing.T) {
	mockRepo := new(MockReadingSettingsRepository)
	service := reading.NewSettingService(mockRepo)

	ctx := context.Background()
	userID := "user789"

	// 设置Mock期望：ExistsByUserID返回错误
	mockRepo.On("ExistsByUserID", ctx, userID).Return(false, errors.New("database error"))

	// 执行测试
	settings, err := service.ResetSetting(ctx, userID)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, settings)
	mockRepo.AssertExpectations(t)
}
