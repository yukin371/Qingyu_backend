package ai_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	aiRepo "Qingyu_backend/repository/interfaces/ai"
	aiService "Qingyu_backend/service/ai"
	"Qingyu_backend/service/base"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ============ Mock 实现 ============

// MockQuotaRepository Mock配额仓储
type MockQuotaRepository struct {
	mock.Mock
}

func (m *MockQuotaRepository) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	args := m.Called(ctx, userID, quotaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aiModels.UserQuota), args.Error(1)
}

func (m *MockQuotaRepository) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*aiModels.UserQuota), args.Error(1)
}

func (m *MockQuotaRepository) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	args := m.Called(ctx, userID, quotaType)
	return args.Error(0)
}

func (m *MockQuotaRepository) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	args := m.Called(ctx, quotaType)
	return args.Error(0)
}

func (m *MockQuotaRepository) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockQuotaRepository) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*aiModels.QuotaTransaction), args.Error(1)
}

func (m *MockQuotaRepository) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*aiModels.QuotaTransaction), args.Error(1)
}

func (m *MockQuotaRepository) GetQuotaStatistics(ctx context.Context, userID string) (*aiRepo.QuotaStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aiRepo.QuotaStatistics), args.Error(1)
}

func (m *MockQuotaRepository) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	args := m.Called(ctx, userID, quotaType, startTime, endTime)
	return args.Int(0), args.Error(1)
}

func (m *MockQuotaRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockRedisClient Mock Redis客户端
type MockRedisClient struct {
	mock.Mock
	cache map[string]string // 简单的内存缓存
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		cache: make(map[string]string),
	}
}

func (m *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	_ = m.Called(ctx, key)
	// 实际从cache中获取
	if val, ok := m.cache[key]; ok {
		return val, nil
	}
	return "", fmt.Errorf("key not found")
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	_ = m.Called(ctx, key, value, ttl)
	// 实际写入cache
	if strVal, ok := value.(string); ok {
		m.cache[key] = strVal
	}
	return nil
}

func (m *MockRedisClient) Delete(ctx context.Context, keys ...string) error {
	_ = m.Called(ctx, keys)
	// 实际删除
	for _, key := range keys {
		delete(m.cache, key)
	}
	return nil
}

func (m *MockRedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	args := m.Called(ctx, keys)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRedisClient) Expire(ctx context.Context, key string, ttl time.Duration) error {
	args := m.Called(ctx, key, ttl)
	return args.Error(0)
}

func (m *MockRedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(time.Duration), args.Error(1)
}

func (m *MockRedisClient) Keys(ctx context.Context, pattern string) ([]string, error) {
	args := m.Called(ctx, pattern)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]string), args.Error(1)
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

// Hash operations
func (m *MockRedisClient) HGet(ctx context.Context, key, field string) (string, error) {
	args := m.Called(ctx, key, field)
	return args.String(0), args.Error(1)
}

func (m *MockRedisClient) HSet(ctx context.Context, key string, values ...interface{}) error {
	args := m.Called(ctx, key, values)
	return args.Error(0)
}

func (m *MockRedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]string), args.Error(1)
}

func (m *MockRedisClient) HDel(ctx context.Context, key string, fields ...string) error {
	args := m.Called(ctx, key, fields)
	return args.Error(0)
}

// Set operations
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
	events []base.Event // 存储已发布的事件
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		events: make([]base.Event, 0),
	}
}

func (m *MockEventBus) Subscribe(eventType string, handler base.EventHandler) error {
	args := m.Called(eventType, handler)
	return args.Error(0)
}

func (m *MockEventBus) Unsubscribe(eventType string, handlerName string) error {
	args := m.Called(eventType, handlerName)
	return args.Error(0)
}

func (m *MockEventBus) Publish(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	// 实际记录事件
	m.events = append(m.events, event)
	return args.Error(0)
}

func (m *MockEventBus) PublishAsync(ctx context.Context, event base.Event) error {
	args := m.Called(ctx, event)
	// 实际记录事件
	m.events = append(m.events, event)
	return args.Error(0)
}

// ============ 测试用例 ============

// TestQuotaServiceWithCache 测试带缓存的配额服务
func TestQuotaServiceWithCache(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockQuotaRepository)
	mockRedis := NewMockRedisClient()
	mockEventBus := NewMockEventBus()

	// 创建带缓存的服务
	service := aiService.NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

	// 测试数据
	userID := "test-user-123"
	quota := &aiModels.UserQuota{
		UserID:         userID,
		QuotaType:      aiModels.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      0,
		RemainingQuota: 1000,
		Status:         aiModels.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1),
	}

	// Mock期望
	mockRepo.On("GetQuotaByUserID", ctx, userID, aiModels.QuotaTypeDaily).Return(quota, nil).Once()
	mockRedis.On("Get", ctx, mock.Anything).Return("", fmt.Errorf("not found")).Maybe()
	mockRedis.On("Set", ctx, mock.Anything, mock.Anything, mock.Anything).Return(nil).Maybe()

	// 第一次调用（缓存未命中）
	result1, err1 := service.GetQuotaInfo(ctx, userID)
	assert.NoError(t, err1)
	assert.NotNil(t, result1)
	assert.Equal(t, 1000, result1.TotalQuota)

	// 验证缓存已写入
	cacheKey := fmt.Sprintf("quota:user:%s:daily", userID)
	assert.Contains(t, mockRedis.cache, cacheKey)

	// 第二次调用（从缓存读取）
	mockRedis.On("Get", ctx, cacheKey).Return(mockRedis.cache[cacheKey], nil).Once()
	result2, err2 := service.GetQuotaInfo(ctx, userID)
	assert.NoError(t, err2)
	assert.NotNil(t, result2)
	assert.Equal(t, 1000, result2.TotalQuota)

	mockRepo.AssertExpectations(t)
}

// TestQuotaWarningTrigger 测试配额预警触发
func TestQuotaWarningTrigger(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockQuotaRepository)
	mockRedis := NewMockRedisClient()
	mockEventBus := NewMockEventBus()

	service := aiService.NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

	userID := "test-user-456"
	// 配额剩余70% → 消费后剩余18%（低于20%预警阈值）
	quota := &aiModels.UserQuota{
		UserID:         userID,
		QuotaType:      aiModels.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      300,
		RemainingQuota: 700, // 70%
		Status:         aiModels.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1), // 明天重置
	}

	// Mock期望
	mockRepo.On("GetQuotaByUserID", ctx, userID, aiModels.QuotaTypeDaily).Return(quota, nil)
	mockRepo.On("UpdateQuota", ctx, mock.Anything).Return(nil)
	mockRepo.On("CreateTransaction", ctx, mock.Anything).Return(nil)
	mockRedis.On("Delete", ctx, mock.Anything).Return(nil).Maybe()
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil)

	// 消费配额，触发预警（消费520，剩余180 = 18%）
	err := service.ConsumeQuota(ctx, userID, 520, "test-service", "test-model", "req-123")
	assert.NoError(t, err)

	// 验证预警事件已发布
	assert.Greater(t, len(mockEventBus.events), 0, "应该发布了预警事件")

	// 检查事件类型
	warningEvent := mockEventBus.events[0]
	assert.Equal(t, "quota.warning", warningEvent.GetEventType())

	mockRepo.AssertExpectations(t)
	mockEventBus.AssertExpectations(t)
}

// TestQuotaCacheInvalidation 测试缓存失效
func TestQuotaCacheInvalidation(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockQuotaRepository)
	mockRedis := NewMockRedisClient()
	mockEventBus := NewMockEventBus()

	service := aiService.NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

	userID := "test-user-789"
	quota := &aiModels.UserQuota{
		UserID:         userID,
		QuotaType:      aiModels.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      100,
		RemainingQuota: 900,
		Status:         aiModels.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1), // 明天重置
	}

	// Mock期望
	mockRepo.On("GetQuotaByUserID", ctx, userID, aiModels.QuotaTypeDaily).Return(quota, nil)
	mockRepo.On("UpdateQuota", ctx, mock.Anything).Return(nil)
	mockRepo.On("CreateTransaction", ctx, mock.Anything).Return(nil)
	mockRedis.On("Delete", ctx, mock.Anything).Return(nil)
	mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Maybe()

	// 预先设置缓存
	cacheKey := fmt.Sprintf("quota:user:%s:daily", userID)
	oldData, _ := json.Marshal(quota)
	mockRedis.cache[cacheKey] = string(oldData)

	// 消费配额（应该清除缓存）
	err := service.ConsumeQuota(ctx, userID, 100, "test-service", "test-model", "req-456")
	assert.NoError(t, err)

	// 验证缓存已删除
	_, exists := mockRedis.cache[cacheKey]
	assert.False(t, exists, "缓存应该已被清除")

	mockRepo.AssertExpectations(t)
	mockRedis.AssertExpectations(t)
}

// TestRechargeQuota 测试配额充值
func TestRechargeQuota(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockQuotaRepository)
	mockRedis := NewMockRedisClient()
	mockEventBus := NewMockEventBus()

	service := aiService.NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

	userID := "test-user-recharge"
	quota := &aiModels.UserQuota{
		UserID:         userID,
		QuotaType:      aiModels.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      700,
		RemainingQuota: 300,
		Status:         aiModels.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1), // 明天重置
	}

	// Mock期望
	mockRepo.On("GetQuotaByUserID", ctx, userID, aiModels.QuotaTypeDaily).Return(quota, nil)
	mockRepo.On("UpdateQuota", ctx, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		updatedQuota := args.Get(1).(*aiModels.UserQuota)
		assert.Equal(t, 1500, updatedQuota.TotalQuota, "总配额应增加500")
		assert.Equal(t, 800, updatedQuota.RemainingQuota, "剩余配额应增加500")
	})
	mockRepo.On("CreateTransaction", ctx, mock.Anything).Return(nil)
	mockRedis.On("Delete", ctx, mock.Anything).Return(nil).Maybe()

	// 充值500配额
	err := service.RechargeQuota(ctx, userID, 500, "购买充值", userID)
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

// TestRedisCachePerformance 测试缓存性能提升
func TestRedisCachePerformance(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockQuotaRepository)
	mockRedis := NewMockRedisClient()
	mockEventBus := NewMockEventBus()

	service := aiService.NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

	userID := "test-user-perf"
	quota := &aiModels.UserQuota{
		UserID:         userID,
		QuotaType:      aiModels.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      200,
		RemainingQuota: 800,
		Status:         aiModels.QuotaStatusActive,
		ResetAt:        time.Now().AddDate(0, 0, 1), // 明天重置
	}

	// 准备缓存数据
	cacheKey := fmt.Sprintf("quota:user:%s:daily", userID)
	cachedData, _ := json.Marshal(quota)
	mockRedis.cache[cacheKey] = string(cachedData)

	// Mock期望：只应该调用Redis，不调用数据库
	mockRedis.On("Get", ctx, cacheKey).Return(string(cachedData), nil)
	// 注意：不设置mockRepo的期望，如果调用了数据库会失败

	// 执行查询（应该从缓存读取）
	start := time.Now()
	result, err := service.GetQuotaInfo(ctx, userID)
	duration := time.Since(start)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1000, result.TotalQuota)

	// 缓存读取应该非常快（<1ms）
	assert.Less(t, duration, 1*time.Millisecond, "缓存读取应该非常快")

	// 验证没有调用数据库
	mockRepo.AssertNotCalled(t, "GetQuotaByUserID")
}

// TestMultipleWarningLevels 测试多级预警
func TestMultipleWarningLevels(t *testing.T) {
	ctx := context.Background()

	testCases := []struct {
		name             string
		remainingPercent float64
		shouldTrigger    bool
		expectedLevel    string
	}{
		{"正常状态", 0.50, false, ""},
		{"接近预警", 0.25, false, ""},
		{"预警级别", 0.15, true, "warning"},
		{"严重级别", 0.05, true, "critical"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockQuotaRepository)
			mockRedis := NewMockRedisClient()
			mockEventBus := NewMockEventBus()

			service := aiService.NewQuotaServiceWithCache(mockRepo, mockRedis, mockEventBus)

			userID := "test-user-" + tc.name
			totalQuota := 1000
			remaining := int(float64(totalQuota) * tc.remainingPercent)

			quota := &aiModels.UserQuota{
				UserID:         userID,
				QuotaType:      aiModels.QuotaTypeDaily,
				TotalQuota:     totalQuota,
				UsedQuota:      totalQuota - remaining,
				RemainingQuota: remaining,
				Status:         aiModels.QuotaStatusActive,
				ResetAt:        time.Now().AddDate(0, 0, 1), // 明天重置
			}

			mockRepo.On("GetQuotaByUserID", ctx, userID, aiModels.QuotaTypeDaily).Return(quota, nil)
			mockRepo.On("UpdateQuota", ctx, mock.Anything).Return(nil)
			mockRepo.On("CreateTransaction", ctx, mock.Anything).Return(nil)
			mockRedis.On("Delete", ctx, mock.Anything).Return(nil).Maybe()
			mockEventBus.On("PublishAsync", ctx, mock.Anything).Return(nil).Maybe()

			// 消费少量配额
			err := service.ConsumeQuota(ctx, userID, 10, "test", "model", "req-123")
			assert.NoError(t, err)

			if tc.shouldTrigger {
				assert.Greater(t, len(mockEventBus.events), 0, "应该触发预警")
				if len(mockEventBus.events) > 0 {
					event := mockEventBus.events[0].(*aiService.QuotaWarningEvent)
					assert.Equal(t, tc.expectedLevel, event.Level)
				}
			} else {
				assert.Equal(t, 0, len(mockEventBus.events), "不应该触发预警")
			}
		})
	}
}
