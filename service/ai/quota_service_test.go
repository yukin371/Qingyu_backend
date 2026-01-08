package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	"Qingyu_backend/models/ai"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	testMock "Qingyu_backend/service/mock"
	"Qingyu_backend/service/base"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient Mock Redis客户端
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	args := m.Called(ctx, key, value, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Delete(ctx context.Context, keys ...string) error {
	args := m.Called(ctx, keys)
	return args.Error(0)
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	args := m.Called(ctx, key, expiration)
	return args.Error(0)
}

func (m *MockRedisClient) Incr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Decr(ctx context.Context, key string) (int64, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) DecrBy(ctx context.Context, key string, value int64) (int64, error) {
	args := m.Called(ctx, key, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisClient) HExists(ctx context.Context, key, field string) (bool, error) {
	args := m.Called(ctx, key, field)
	return args.Bool(0), args.Error(1)
}

func (m *MockRedisClient) HIncrBy(ctx context.Context, key, field string, value int64) (int64, error) {
	args := m.Called(ctx, key, field, value)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	args := m.Called(ctx, keys)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]interface{}), args.Error(1)
}

func (m *MockRedisClient) MSet(ctx context.Context, pairs ...interface{}) error {
	args := m.Called(ctx, pairs)
	return args.Error(0)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisClient) SAdd(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockRedisClient) SRem(ctx context.Context, key string, members ...interface{}) error {
	args := m.Called(ctx, key, members)
	return args.Error(0)
}

func (m *MockRedisClient) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockRedisClient) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRedisClient) GetClient() interface{} {
	args := m.Called()
	return args.Get(0)
}

// MockEventBus Mock事件总线
type MockEventBus struct {
	mock.Mock
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerID string) error {
	args := m.Called(eventType, handlerID)
	return args.Error(0)
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// TestNewQuotaService 测试创建配额服务
func TestNewQuotaService(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)

	service := NewQuotaService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.quotaRepo)
	assert.Nil(t, service.redisClient)
	assert.Equal(t, 5*time.Minute, service.cacheTTL)
	assert.Equal(t, 0.2, service.warningThreshold)
	assert.Equal(t, 0.1, service.criticalThreshold)
}

// TestNewQuotaServiceWithCache 测试创建带缓存的配额服务
func TestNewQuotaServiceWithCache(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	mockRedis := new(MockRedisClient)
	mockEventBus := new(MockEventBus)

	service := NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.quotaRepo)
	assert.Equal(t, mockRedis, service.redisClient)
	assert.Equal(t, mockEventBus, service.eventBus)
	assert.Equal(t, 5*time.Minute, service.cacheTTL)
}

// TestInitializeUserQuota_Success 测试成功初始化用户配额
func TestInitializeUserQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	userRole := "reader"
	membershipLevel := "normal"

	// 模拟配额不存在
	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(nil, ai.ErrQuotaNotFound)
	mockRepo.On("CreateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil)

	err := service.InitializeUserQuota(ctx, userID, userRole, membershipLevel)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestInitializeUserQuota_AlreadyExists 测试配额已存在的情况
func TestInitializeUserQuota_AlreadyExists(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	existingQuota := &ai.UserQuota{
		UserID:    userID,
		QuotaType: ai.QuotaTypeDaily,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(existingQuota, nil)

	err := service.InitializeUserQuota(ctx, userID, "reader", "normal")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCheckQuota_Success 测试配额检查成功
func TestCheckQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      100,
		RemainingQuota: 900,
		Status:         ai.QuotaStatusActive,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)

	err := service.CheckQuota(ctx, userID, 100)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCheckQuota_NotFound_AutoInit 测试配额不存在时自动初始化
func TestCheckQuota_NotFound_AutoInit(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	// 第一次返回不存在
	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(nil, ai.ErrQuotaNotFound).Once()
	// 初始化后返回配额
	mockRepo.On("CreateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil).Once()
	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(&ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     100,
		UsedQuota:      0,
		RemainingQuota: 100,
		Status:         ai.QuotaStatusActive,
	}, nil).Once()

	err := service.CheckQuota(ctx, userID, 50)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCheckQuota_Insufficient 测试配额不足
func TestCheckQuota_Insufficient(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     100,
		UsedQuota:      90,
		RemainingQuota: 10,
		Status:         ai.QuotaStatusActive,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)

	err := service.CheckQuota(ctx, userID, 50)

	assert.Error(t, err)
	assert.Equal(t, ai.ErrInsufficientQuota, err)
	mockRepo.AssertExpectations(t)
}

// TestCheckQuota_Exhausted 测试配额耗尽
func TestCheckQuota_Exhausted(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     100,
		UsedQuota:      100,
		RemainingQuota: 0,
		Status:         ai.QuotaStatusExhausted,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)

	err := service.CheckQuota(ctx, userID, 50)

	assert.Error(t, err)
	assert.Equal(t, ai.ErrQuotaExhausted, err)
	mockRepo.AssertExpectations(t)
}

// TestConsumeQuota_Success 测试成功消费配额
func TestConsumeQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	amount := 100

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      100,
		RemainingQuota: 900,
		Status:         ai.QuotaStatusActive,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil).Once()
	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil).Once()
	mockRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*ai.QuotaTransaction")).Return(nil).Once()

	err := service.ConsumeQuota(ctx, userID, amount, "chat", "claude-3", "req123")

	assert.NoError(t, err)
	assert.Equal(t, 200, quota.UsedQuota)
	assert.Equal(t, 800, quota.RemainingQuota)
	mockRepo.AssertExpectations(t)
}

// TestRestoreQuota_Success 测试成功恢复配额
func TestRestoreQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	amount := 50

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      200,
		RemainingQuota: 800,
		Status:         ai.QuotaStatusActive,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)
	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil)
	mockRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*ai.QuotaTransaction")).Return(nil)

	err := service.RestoreQuota(ctx, userID, amount, "错误回滚")

	assert.NoError(t, err)
	assert.Equal(t, 150, quota.UsedQuota)
	assert.Equal(t, 850, quota.RemainingQuota)
	mockRepo.AssertExpectations(t)
}

// TestGetQuotaInfo_WithCache 测试从缓存获取配额信息
func TestGetQuotaInfo_WithCache(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	mockRedis := new(MockRedisClient)
	service := NewQuotaServiceWithCache(mockRepo, mockRedis, nil)

	ctx := context.Background()
	userID := "user123"

	// 模拟缓存命中
	mockRedis.On("Get", ctx, "quota:user:user123:daily").Return(`{"user_id":"user123","quota_type":"daily","total_quota":1000,"used_quota":100,"remaining_quota":900,"status":"active"}`, nil)
	mockRedis.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	result, err := service.GetQuotaInfo(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	mockRedis.AssertExpectations(t)
}

// TestGetQuotaInfo_CacheMiss 测试缓存未命中
func TestGetQuotaInfo_CacheMiss(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	mockRedis := new(MockRedisClient)
	service := NewQuotaServiceWithCache(mockRepo, mockRedis, nil)

	ctx := context.Background()
	userID := "user123"

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      100,
		RemainingQuota: 900,
		Status:         ai.QuotaStatusActive,
	}

	// 模拟缓存未命中
	mockRedis.On("Get", ctx, "quota:user:user123:daily").Return("", errors.New("cache miss"))
	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)
	mockRedis.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := service.GetQuotaInfo(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, userID, result.UserID)
	mockRedis.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestRechargeQuota_Success 测试成功充值配额
func TestRechargeQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	amount := 500
	reason := "活动赠送"
	operatorID := "admin123"

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      200,
		RemainingQuota: 800,
		Status:         ai.QuotaStatusActive,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)
	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil)
	mockRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*ai.QuotaTransaction")).Return(nil)

	err := service.RechargeQuota(ctx, userID, amount, reason, operatorID)

	assert.NoError(t, err)
	assert.Equal(t, 1500, quota.TotalQuota)
	assert.Equal(t, 1300, quota.RemainingQuota)
	mockRepo.AssertExpectations(t)
}

// TestRechargeQuota_InvalidAmount 测试无效的充值金额
func TestRechargeQuota_InvalidAmount(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()

	err := service.RechargeQuota(ctx, "user123", -100, "测试", "admin123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "充值金额必须大于0")
}

// TestSuspendUserQuota_Success 测试成功暂停用户配额
func TestSuspendUserQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	quotas := []*ai.UserQuota{
		{
			UserID:    userID,
			QuotaType: ai.QuotaTypeDaily,
			Status:    ai.QuotaStatusActive,
		},
		{
			UserID:    userID,
			QuotaType: ai.QuotaTypeMonthly,
			Status:    ai.QuotaStatusActive,
		},
	}

	mockRepo.On("GetAllQuotasByUserID", ctx, userID).Return(quotas, nil)
	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil).Twice()

	err := service.SuspendUserQuota(ctx, userID)

	assert.NoError(t, err)
	for _, quota := range quotas {
		assert.Equal(t, ai.QuotaStatusSuspended, quota.Status)
	}
	mockRepo.AssertExpectations(t)
}

// TestActivateUserQuota_Success 测试成功激活用户配额
func TestActivateUserQuota_Success(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	quotas := []*ai.UserQuota{
		{
			UserID:    userID,
			QuotaType: ai.QuotaTypeDaily,
			Status:    ai.QuotaStatusSuspended,
		},
	}

	mockRepo.On("GetAllQuotasByUserID", ctx, userID).Return(quotas, nil)
	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil)

	err := service.ActivateUserQuota(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, ai.QuotaStatusActive, quotas[0].Status)
	mockRepo.AssertExpectations(t)
}

// TestResetDailyQuotas 测试重置日配额
func TestResetDailyQuotas(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()

	mockRepo.On("BatchResetQuotas", ctx, ai.QuotaTypeDaily).Return(nil)

	err := service.ResetDailyQuotas(ctx)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestResetMonthlyQuotas 测试重置月配额
func TestResetMonthlyQuotas(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()

	mockRepo.On("BatchResetQuotas", ctx, ai.QuotaTypeMonthly).Return(nil)

	err := service.ResetMonthlyQuotas(ctx)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUpdateUserQuota_CreateNew 测试创建新配额
func TestUpdateUserQuota_CreateNew(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	totalQuota := 2000

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(nil, ai.ErrQuotaNotFound)
	mockRepo.On("CreateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil)

	err := service.UpdateUserQuota(ctx, userID, ai.QuotaTypeDaily, totalQuota)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestUpdateUserQuota_Existing 测试更新现有配额
func TestUpdateUserQuota_Existing(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	totalQuota := 2000

	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      300,
		RemainingQuota: 700,
		Status:         ai.QuotaStatusActive,
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).Return(quota, nil)
	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).Return(nil)

	err := service.UpdateUserQuota(ctx, userID, ai.QuotaTypeDaily, totalQuota)

	assert.NoError(t, err)
	assert.Equal(t, totalQuota, quota.TotalQuota)
	assert.Equal(t, 1700, quota.RemainingQuota)
	mockRepo.AssertExpectations(t)
}

// TestGetAllQuotas 测试获取所有配额
func TestGetAllQuotas(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	expectedQuotas := []*ai.UserQuota{
		{
			UserID:    userID,
			QuotaType: ai.QuotaTypeDaily,
		},
		{
			UserID:    userID,
			QuotaType: ai.QuotaTypeMonthly,
		},
	}

	mockRepo.On("GetAllQuotasByUserID", ctx, userID).Return(expectedQuotas, nil)

	quotas, err := service.GetAllQuotas(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(quotas))
	mockRepo.AssertExpectations(t)
}

// TestGetQuotaStatistics 测试获取配额统计
func TestGetQuotaStatistics(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"

	expectedStats := &aiRepo.QuotaStatistics{
		TotalQuota:     10000,
		UsedQuota:      2000,
		RemainingQuota: 8000,
	}

	mockRepo.On("GetQuotaStatistics", ctx, userID).Return(expectedStats, nil)

	stats, err := service.GetQuotaStatistics(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedStats, stats)
	mockRepo.AssertExpectations(t)
}

// TestGetTransactionHistory 测试获取事务历史
func TestGetTransactionHistory(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	service := NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user123"
	limit := 10
	offset := 0

	expectedTransactions := []*ai.QuotaTransaction{
		{
			UserID: userID,
			Amount: 100,
		},
	}

	mockRepo.On("GetTransactionsByUserID", ctx, userID, limit, offset).Return(expectedTransactions, nil)

	transactions, err := service.GetTransactionHistory(ctx, userID, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(transactions))
	mockRepo.AssertExpectations(t)
}

// TestQuotaWarningEvent 测试配额预警事件
func TestQuotaWarningEvent(t *testing.T) {
	event := &QuotaWarningEvent{
		UserID:           "user123",
		QuotaType:        "daily",
		TotalQuota:       1000,
		RemainingQuota:   50,
		UsedQuota:        950,
		RemainingPercent: 5.0,
		Level:            "critical",
		Timestamp:        time.Now(),
	}

	assert.Equal(t, "quota.warning", event.GetEventType())
	assert.Equal(t, event, event.GetEventData())
	assert.Equal(t, event.Timestamp, event.GetTimestamp())
	assert.Equal(t, "QuotaService", event.GetSource())
}

// TestCacheQuota 测试缓存配额
func TestCacheQuota(t *testing.T) {
	mockRedis := new(MockRedisClient)
	service := NewQuotaServiceWithCache(nil, mockRedis, nil)

	ctx := context.Background()
	quota := &ai.UserQuota{
		UserID:         "user123",
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      100,
		RemainingQuota: 900,
		Status:         ai.QuotaStatusActive,
	}

	mockRedis.On("Set", ctx, "quota:user:user123:daily", mock.Anything, 5*time.Minute).Return(nil)

	// 这是一个私有方法，但我们可以通过ConsumeQuota间接测试
	// 或者添加一个公开的测试方法
	// 这里我们直接调用私有方法用于测试
	service.cacheQuota(ctx, quota)

	mockRedis.AssertExpectations(t)
}

// TestInvalidateQuotaCache 测试清除配额缓存
func TestInvalidateQuotaCache(t *testing.T) {
	mockRedis := new(MockRedisClient)
	service := NewQuotaServiceWithCache(nil, mockRedis, nil)

	ctx := context.Background()
	userID := "user123"

	mockRedis.On("Delete", ctx, "quota:user:user123:daily").Return(nil)

	service.invalidateQuotaCache(ctx, userID, ai.QuotaTypeDaily)

	mockRedis.AssertExpectations(t)
}

// TestCheckAndPublishWarning_CriticalLevel 测试严重级别预警
func TestCheckAndPublishWarning_CriticalLevel(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	mockEventBus := new(MockEventBus)
	service := NewQuotaServiceWithCache(mockRepo, nil, mockEventBus)

	ctx := context.Background()
	quota := &ai.UserQuota{
		UserID:         "user123",
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      950,
		RemainingQuota: 50,
		Status:         ai.QuotaStatusActive,
	}

	mockEventBus.On("PublishAsync", ctx, mock.AnythingOfType("*ai.QuotaWarningEvent")).Return(nil)

	// 这是一个私有方法，我们通过ConsumeQuota间接测试
	service.checkAndPublishWarning(ctx, quota)

	mockEventBus.AssertExpectations(t)
}

// TestCheckAndPublishWarning_NoWarning 测试不需要预警的情况
func TestCheckAndPublishWarning_NoWarning(t *testing.T) {
	mockRepo := new(testMock.MockQuotaRepository)
	mockEventBus := new(MockEventBus)
	service := NewQuotaServiceWithCache(mockRepo, nil, mockEventBus)

	ctx := context.Background()
	quota := &ai.UserQuota{
		UserID:         "user123",
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      100,
		RemainingQuota: 900,
		Status:         ai.QuotaStatusActive,
	}

	// 剩余90%，不需要预警
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(errors.New("should not be called")).Maybe()

	service.checkAndPublishWarning(ctx, quota)

	// 验证PublishAsync没有被调用
	mockEventBus.AssertNotCalled(t, "PublishAsync", ctx, mock.Anything)
}
