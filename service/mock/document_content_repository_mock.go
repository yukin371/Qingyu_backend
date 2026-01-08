// Code generated for testing purposes
// Package mock provides mock implementations for testing
package mock

import (
	"context"

	"Qingyu_backend/models/writer"
	infrastructure "Qingyu_backend/repository/interfaces/infrastructure"

	"github.com/stretchr/testify/mock"
)

// MockDocumentContentRepository is a mock implementation of DocumentContentRepository interface
type MockDocumentContentRepository struct {
	mock.Mock
}

// Create mocks base method
func (m *MockDocumentContentRepository) Create(ctx context.Context, content *writer.DocumentContent) error {
	args := m.Called(ctx, content)
	return args.Error(0)
}

// GetByID mocks base method
func (m *MockDocumentContentRepository) GetByID(ctx context.Context, id string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

// Update mocks base method
func (m *MockDocumentContentRepository) Update(ctx context.Context, id string, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

// Delete mocks base method
func (m *MockDocumentContentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// List mocks base method
func (m *MockDocumentContentRepository) List(ctx context.Context, filter infrastructure.Filter) ([]*writer.DocumentContent, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*writer.DocumentContent), args.Error(1)
}

// Exists mocks base method
func (m *MockDocumentContentRepository) Exists(ctx context.Context, id string) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

// GetByDocumentID mocks base method
func (m *MockDocumentContentRepository) GetByDocumentID(ctx context.Context, documentID string) (*writer.DocumentContent, error) {
	args := m.Called(ctx, documentID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*writer.DocumentContent), args.Error(1)
}

// UpdateWithVersion mocks base method
func (m *MockDocumentContentRepository) UpdateWithVersion(ctx context.Context, documentID string, content string, expectedVersion int) error {
	args := m.Called(ctx, documentID, content, expectedVersion)
	return args.Error(0)
}

// BatchUpdateContent mocks base method
func (m *MockDocumentContentRepository) BatchUpdateContent(ctx context.Context, updates map[string]string) error {
	args := m.Called(ctx, updates)
	return args.Error(0)
}

// GetContentStats mocks base method
func (m *MockDocumentContentRepository) GetContentStats(ctx context.Context, documentID string) (int, int, error) {
	args := m.Called(ctx, documentID)
	return args.Int(0), args.Int(1), args.Error(2)
}

// StoreToGridFS mocks base method
func (m *MockDocumentContentRepository) StoreToGridFS(ctx context.Context, documentID string, content []byte) (string, error) {
	args := m.Called(ctx, documentID, content)
	return args.String(0), args.Error(1)
}

// LoadFromGridFS mocks base method
func (m *MockDocumentContentRepository) LoadFromGridFS(ctx context.Context, gridFSID string) ([]byte, error) {
	args := m.Called(ctx, gridFSID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]byte), args.Error(1)
}

// CreateWithTransaction mocks base method
func (m *MockDocumentContentRepository) CreateWithTransaction(ctx context.Context, content *writer.DocumentContent, callback func(ctx context.Context) error) error {
	args := m.Called(ctx, content, callback)
	return args.Error(0)
}

// CheckHealth mocks base method
func (m *MockDocumentContentRepository) CheckHealth(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Health mocks base method
func (m *MockDocumentContentRepository) Health(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// Count mocks base method
func (m *MockDocumentContentRepository) Count(ctx context.Context, filter infrastructure.Filter) (int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).(int64), args.Error(1)
}
