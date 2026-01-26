package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	modelsMessaging "Qingyu_backend/models/messaging"
)

// MockConversationRepository 模拟会话Repository
type MockConversationRepository struct {
	mock.Mock
}

// Create 模拟创建会话
func (m *MockConversationRepository) Create(ctx context.Context, conversation *modelsMessaging.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

// FindByID 模拟根据ID查找会话
func (m *MockConversationRepository) FindByID(ctx context.Context, id string) (*modelsMessaging.Conversation, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*modelsMessaging.Conversation), args.Error(1)
}

// FindByParticipants 模拟根据参与者查找会话
func (m *MockConversationRepository) FindByParticipants(ctx context.Context, participantIDs []string) (*modelsMessaging.Conversation, error) {
	args := m.Called(ctx, participantIDs)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*modelsMessaging.Conversation), args.Error(1)
}

// FindByUserID 模拟根据用户ID查找会话列表
func (m *MockConversationRepository) FindByUserID(ctx context.Context, userID string, page, pageSize int) ([]*modelsMessaging.Conversation, int64, error) {
	args := m.Called(ctx, userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*modelsMessaging.Conversation), args.Get(1).(int64), args.Error(2)
}

// Update 模拟更新会话
func (m *MockConversationRepository) Update(ctx context.Context, conversation *modelsMessaging.Conversation) error {
	args := m.Called(ctx, conversation)
	return args.Error(0)
}

// IncrementUnreadCount 模拟增加未读计数
func (m *MockConversationRepository) IncrementUnreadCount(ctx context.Context, conversationID, userID string) error {
	args := m.Called(ctx, conversationID, userID)
	return args.Error(0)
}
