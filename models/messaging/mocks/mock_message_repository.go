package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	modelsMessaging "Qingyu_backend/models/messaging"
)

// MockMessageRepository 模拟消息Repository
type MockMessageRepository struct {
	mock.Mock
}

// Create 模拟创建消息
func (m *MockMessageRepository) Create(ctx context.Context, message *modelsMessaging.DirectMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

// FindByConversationID 模拟根据会话ID查找消息
func (m *MockMessageRepository) FindByConversationID(ctx context.Context, conversationID, userID string, page, pageSize int, before, after *string) ([]*modelsMessaging.DirectMessage, int, error) {
	args := m.Called(ctx, conversationID, userID, page, pageSize, before, after)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*modelsMessaging.DirectMessage), args.Int(1), args.Error(2)
}

// MarkConversationRead 模拟标记会话已读
func (m *MockMessageRepository) MarkConversationRead(ctx context.Context, conversationID, userID string, readAt time.Time) (int, error) {
	args := m.Called(ctx, conversationID, userID, readAt)
	return args.Int(0), args.Error(1)
}

// CountUnreadInConversation 模拟统计会话未读消息数
func (m *MockMessageRepository) CountUnreadInConversation(ctx context.Context, conversationID, userID string) (int, error) {
	args := m.Called(ctx, conversationID, userID)
	return args.Int(0), args.Error(1)
}

// CountUnread 模拟统计未读消息数 (用于测试的其他方法)
func (m *MockMessageRepository) CountUnread(ctx context.Context, conversationID, userID string) (int64, error) {
	args := m.Called(ctx, conversationID, userID)
	return args.Get(0).(int64), args.Error(1)
}

// MarkMultipleAsRead 模拟批量标记为已读
func (m *MockMessageRepository) MarkMultipleAsRead(ctx context.Context, conversationID, userID string, readAt time.Time) (int64, error) {
	args := m.Called(ctx, conversationID, userID, readAt)
	return args.Get(0).(int64), args.Error(1)
}
