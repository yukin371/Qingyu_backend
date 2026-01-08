// Code generated for testing purposes
// Package mock provides mock implementations for testing
package mock

import (
	"context"
	"time"

	aiModels "Qingyu_backend/models/ai"
	"Qingyu_backend/repository/interfaces/ai"

	"github.com/stretchr/testify/mock"
)

// MockQuotaRepository is a mock implementation of QuotaRepository interface
type MockQuotaRepository struct {
	mock.Mock
}

// CreateQuota mocks base method
func (m *MockQuotaRepository) CreateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

// GetQuotaByUserID mocks base method
func (m *MockQuotaRepository) GetQuotaByUserID(ctx context.Context, userID string, quotaType aiModels.QuotaType) (*aiModels.UserQuota, error) {
	args := m.Called(ctx, userID, quotaType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*aiModels.UserQuota), args.Error(1)
}

// UpdateQuota mocks base method
func (m *MockQuotaRepository) UpdateQuota(ctx context.Context, quota *aiModels.UserQuota) error {
	args := m.Called(ctx, quota)
	return args.Error(0)
}

// DeleteQuota mocks base method
func (m *MockQuotaRepository) DeleteQuota(ctx context.Context, userID string, quotaType aiModels.QuotaType) error {
	args := m.Called(ctx, userID, quotaType)
	return args.Error(0)
}

// GetAllQuotasByUserID mocks base method
func (m *MockQuotaRepository) GetAllQuotasByUserID(ctx context.Context, userID string) ([]*aiModels.UserQuota, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*aiModels.UserQuota), args.Error(1)
}

// BatchResetQuotas mocks base method
func (m *MockQuotaRepository) BatchResetQuotas(ctx context.Context, quotaType aiModels.QuotaType) error {
	args := m.Called(ctx, quotaType)
	return args.Error(0)
}

// CreateTransaction mocks base method
func (m *MockQuotaRepository) CreateTransaction(ctx context.Context, transaction *aiModels.QuotaTransaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

// GetTransactionsByUserID mocks base method
func (m *MockQuotaRepository) GetTransactionsByUserID(ctx context.Context, userID string, limit, offset int) ([]*aiModels.QuotaTransaction, error) {
	args := m.Called(ctx, userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*aiModels.QuotaTransaction), args.Error(1)
}

// GetTransactionsByTimeRange mocks base method
func (m *MockQuotaRepository) GetTransactionsByTimeRange(ctx context.Context, userID string, startTime, endTime time.Time) ([]*aiModels.QuotaTransaction, error) {
	args := m.Called(ctx, userID, startTime, endTime)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*aiModels.QuotaTransaction), args.Error(1)
}

// GetQuotaStatistics mocks base method
func (m *MockQuotaRepository) GetQuotaStatistics(ctx context.Context, userID string) (*ai.QuotaStatistics, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ai.QuotaStatistics), args.Error(1)
}

// GetTotalConsumption mocks base method
func (m *MockQuotaRepository) GetTotalConsumption(ctx context.Context, userID string, quotaType aiModels.QuotaType, startTime, endTime time.Time) (int, error) {
	args := m.Called(ctx, userID, quotaType, startTime, endTime)
	return args.Int(0), args.Error(1)
}

// Health mocks base method
func (m *MockQuotaRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
