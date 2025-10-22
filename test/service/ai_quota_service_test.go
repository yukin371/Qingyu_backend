package service_test

import (
	"context"
	"testing"
	"time"

	"Qingyu_backend/models/ai"
	"Qingyu_backend/repository/interfaces"
	aiService "Qingyu_backend/service/ai"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockQuotaRepository 模拟配额Repository
type MockQuotaRepository struct {
	mock.Mock
}

func (m *MockQuotaRepository) CreateQuota(ctx context.Context, quota *ai.UserQuota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) GetQuotaByUserID(ctx context.Context, userID string, quotaType ai.QuotaType) (*ai.UserQuota, error) {
	args := m.Called(ctx, userID, quotaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ai.UserQuota), args.Error(1)
}

func (m *MockQuotaRepository) UpdateQuota(ctx context.Context, quota *ai.UserQuota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

func (m *MockQuotaRepository) DeleteQuota(ctx context.Context, userID string, quotaType ai.QuotaType) error {
	args := m.Called(ctx, userID, quotaType)
	return args.Error(0)
}

func (m *MockQuotaRepository) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*ai.UserQuota, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ai.UserQuota), args.Error(1)
}

func (m *MockQuotaRepository) BatchResetQuotas(ctx context.Context, quotaType ai.QuotaType) error {
	args := m.Called(ctx, quotaType)
	return args.Error(0)
}

func (m *MockQuotaRepository) CreateTransaction(ctx context.Context, transaction *ai.QuotaTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockQuotaRepository) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*ai.QuotaTransaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ai.QuotaTransaction), args.Error(1)
}

func (m *MockQuotaRepository) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*ai.QuotaTransaction, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*ai.QuotaTransaction), args.Error(1)
}

func (m *MockQuotaRepository) GetQuotaStatistics(ctx context.Context, userID string) (*interfaces.QuotaStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*interfaces.QuotaStatistics), args.Error(1)
}

func (m *MockQuotaRepository) GetTotalConsumption(ctx context.Context, userID string, quotaType ai.QuotaType, startTime, endTime time.Time) (int, error) {
	args := m.Called(ctx, userID, quotaType, startTime, endTime)
	return args.Int(0), args.Error(1)
}

func (m *MockQuotaRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// TestQuotaService_InitializeUserQuota 测试初始化用户配额
func TestQuotaService_InitializeUserQuota(t *testing.T) {
	mockRepo := new(MockQuotaRepository)
	service := aiService.NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user_123"
	userRole := "writer"
	membershipLevel := "signed"

	// 模拟不存在的配额
	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).
		Return(nil, ai.ErrQuotaNotFound)

	// 模拟创建配额成功
	mockRepo.On("CreateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).
		Return(nil)

	err := service.InitializeUserQuota(ctx, userID, userRole, membershipLevel)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_CheckQuota 测试检查配额
func TestQuotaService_CheckQuota(t *testing.T) {
	mockRepo := new(MockQuotaRepository)
	service := aiService.NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user_123"

	// 创建测试配额
	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      200,
		RemainingQuota: 800,
		Status:         ai.QuotaStatusActive,
		ResetAt:        time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).
		Return(quota, nil)

	// 测试配额充足的情况
	err := service.CheckQuota(ctx, userID, 500)
	assert.NoError(t, err)

	// 测试配额不足的情况
	err = service.CheckQuota(ctx, userID, 1000)
	assert.Error(t, err)
	assert.Equal(t, ai.ErrInsufficientQuota, err)

	mockRepo.AssertExpectations(t)
}

// TestQuotaService_ConsumeQuota 测试消费配额
func TestQuotaService_ConsumeQuota(t *testing.T) {
	mockRepo := new(MockQuotaRepository)
	service := aiService.NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user_123"
	amount := 100

	// 创建测试配额
	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      200,
		RemainingQuota: 800,
		Status:         ai.QuotaStatusActive,
		ResetAt:        time.Now().Add(24 * time.Hour),
		Metadata:       &ai.QuotaMetadata{},
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).
		Return(quota, nil)

	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).
		Return(nil)

	mockRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*ai.QuotaTransaction")).
		Return(nil)

	err := service.ConsumeQuota(ctx, userID, amount, "test_service", "gpt-4", "req_123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_RestoreQuota 测试恢复配额
func TestQuotaService_RestoreQuota(t *testing.T) {
	mockRepo := new(MockQuotaRepository)
	service := aiService.NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user_123"
	amount := 100

	// 创建测试配额
	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      500,
		RemainingQuota: 500,
		Status:         ai.QuotaStatusActive,
		ResetAt:        time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).
		Return(quota, nil)

	mockRepo.On("UpdateQuota", ctx, mock.AnythingOfType("*ai.UserQuota")).
		Return(nil)

	mockRepo.On("CreateTransaction", ctx, mock.AnythingOfType("*ai.QuotaTransaction")).
		Return(nil)

	err := service.RestoreQuota(ctx, userID, amount, "测试恢复")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestQuotaService_QuotaExhausted 测试配额用尽
func TestQuotaService_QuotaExhausted(t *testing.T) {
	mockRepo := new(MockQuotaRepository)
	service := aiService.NewQuotaService(mockRepo)

	ctx := context.Background()
	userID := "user_123"

	// 创建配额已用尽的配额
	quota := &ai.UserQuota{
		UserID:         userID,
		QuotaType:      ai.QuotaTypeDaily,
		TotalQuota:     1000,
		UsedQuota:      1000,
		RemainingQuota: 0,
		Status:         ai.QuotaStatusExhausted,
		ResetAt:        time.Now().Add(24 * time.Hour),
	}

	mockRepo.On("GetQuotaByUserID", ctx, userID, ai.QuotaTypeDaily).
		Return(quota, nil)

	// 测试配额用尽的情况
	err := service.CheckQuota(ctx, userID, 100)
	assert.Error(t, err)
	assert.Equal(t, ai.ErrQuotaExhausted, err)

	mockRepo.AssertExpectations(t)
}
