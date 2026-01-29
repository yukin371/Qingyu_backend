package social_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"

	"Qingyu_backend/models/social"
	socialRepo "Qingyu_backend/repository/interfaces/social"
	socialMongo "Qingyu_backend/repository/mongodb/social"
	"Qingyu_backend/test/testutil"
)

// setupMessageRepo 测试辅助函数
func setupMessageRepo(t *testing.T) (socialRepo.MessageRepository, *mongo.Database, context.Context, func()) {
	db, cleanup := testutil.SetupTestDB(t)
	repo := socialMongo.NewMongoMessageRepository(db)
	ctx := context.Background()
	return repo, db, ctx, cleanup
}

// TestMongoMessageRepository_CreateConversation 测试创建会话
func TestMongoMessageRepository_CreateConversation(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount: map[string]int{
			"user1": 0,
			"user2": 0,
		},
	}

	// Act
	err := repo.CreateConversation(ctx, conversation)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, conversation.ID)
	assert.False(t, conversation.CreatedAt.IsZero())
	assert.False(t, conversation.UpdatedAt.IsZero())
}

// TestMongoMessageRepository_GetConversationByParticipants 测试根据参与者获取会话
func TestMongoMessageRepository_GetConversationByParticipants(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 创建测试会话
	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 0},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	// Act
	found, err := repo.GetConversationByParticipants(ctx, []string{"user1", "user2"})

	// Assert
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, conversation.ID, found.ID)
	assert.ElementsMatch(t, conversation.Participants, found.Participants)
}

// TestMongoMessageRepository_CreateMessage 测试创建消息
func TestMongoMessageRepository_CreateMessage(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 先创建会话
	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 0},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	message := &social.Message{
		ConversationID: conversation.ID.Hex(),
		SenderID:       "user1",
		ReceiverID:     "user2",
		Content:        "Hello, World!",
		MessageType:    "text",
	}

	// Act
	err = repo.CreateMessage(ctx, message)

	// Assert
	require.NoError(t, err)
	assert.NotEmpty(t, message.ID)
	assert.False(t, message.CreatedAt.IsZero())
}

// TestMongoMessageRepository_GetMessageByID 测试根据ID获取消息
func TestMongoMessageRepository_GetMessageByID(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 创建会话和消息
	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 0},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	message := &social.Message{
		ConversationID: conversation.ID.Hex(),
		SenderID:       "user1",
		ReceiverID:     "user2",
		Content:        "Test message",
		MessageType:    "text",
	}
	err = repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	// Act
	found, err := repo.GetMessageByID(ctx, message.ID.Hex())

	// Assert
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, message.Content, found.Content)
	assert.Equal(t, message.SenderID, found.SenderID)
}

// TestMongoMessageRepository_GetMessagesByConversation 测试获取会话的消息列表
func TestMongoMessageRepository_GetMessagesByConversation(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 创建会话
	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 0},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	// 创建多条消息
	for i := 0; i < 3; i++ {
		message := &social.Message{
			ConversationID: conversation.ID.Hex(),
			SenderID:       "user1",
			ReceiverID:     "user2",
			Content:        "Test message",
			MessageType:    "text",
		}
		err = repo.CreateMessage(ctx, message)
		require.NoError(t, err)
	}

	// Act
	messages, total, err := repo.GetMessagesByConversation(ctx, conversation.ID.Hex(), 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, messages, 3)
}

// TestMongoMessageRepository_MarkMessageAsRead 测试标记消息为已读
func TestMongoMessageRepository_MarkMessageAsRead(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 创建会话和消息
	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 1},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	message := &social.Message{
		ConversationID: conversation.ID.Hex(),
		SenderID:       "user1",
		ReceiverID:     "user2",
		Content:        "Test message",
		MessageType:    "text",
	}
	err = repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	// Act
	err = repo.MarkMessageAsRead(ctx, message.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证消息已标记为已读
	found, err := repo.GetMessageByID(ctx, message.ID.Hex())
	require.NoError(t, err)
	assert.True(t, found.IsRead)
	assert.NotNil(t, found.ReadAt)
}

// TestMongoMessageRepository_CountUnreadMessages 测试统计未读消息数
func TestMongoMessageRepository_CountUnreadMessages(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 创建会话
	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 2},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	// 创建未读消息
	for i := 0; i < 2; i++ {
		message := &social.Message{
			ConversationID: conversation.ID.Hex(),
			SenderID:       "user1",
			ReceiverID:     "user2",
			Content:        "Test message",
			MessageType:    "text",
		}
		err = repo.CreateMessage(ctx, message)
		require.NoError(t, err)
	}

	// Act
	count, err := repo.CountUnreadMessages(ctx, "user2")

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

// TestMongoMessageRepository_GetUserConversations 测试获取用户的会话列表
func TestMongoMessageRepository_GetUserConversations(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// 创建多个会话
	conversations := []*social.Conversation{
		{
			Participants: []string{"user1", "user2"},
			UnreadCount:  map[string]int{"user1": 0, "user2": 0},
		},
		{
			Participants: []string{"user1", "user3"},
			UnreadCount:  map[string]int{"user1": 1, "user3": 0},
		},
	}

	for _, conv := range conversations {
		err := repo.CreateConversation(ctx, conv)
		require.NoError(t, err)
	}

	// Act
	result, total, err := repo.GetUserConversations(ctx, "user1", 1, 10)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

// TestMongoMessageRepository_UpdateLastMessage 测试更新会话最后一条消息
func TestMongoMessageRepository_UpdateLastMessage(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 0},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	message := &social.Message{
		ConversationID: conversation.ID.Hex(),
		SenderID:       "user1",
		ReceiverID:     "user2",
		Content:        "Last message",
		MessageType:    "text",
	}
	err = repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	// Act
	err = repo.UpdateLastMessage(ctx, conversation.ID.Hex(), message)

	// Assert
	require.NoError(t, err)

	// 验证会话的LastMessage已更新
	updated, err := repo.GetConversationByID(ctx, conversation.ID.Hex())
	require.NoError(t, err)
	assert.NotNil(t, updated.LastMessage)
	assert.Equal(t, message.Content, updated.LastMessage.Content)
}

// TestMongoMessageRepository_DeleteMessage 测试删除消息（软删除）
func TestMongoMessageRepository_DeleteMessage(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	conversation := &social.Conversation{
		Participants: []string{"user1", "user2"},
		UnreadCount:  map[string]int{"user1": 0, "user2": 0},
	}
	err := repo.CreateConversation(ctx, conversation)
	require.NoError(t, err)

	message := &social.Message{
		ConversationID: conversation.ID.Hex(),
		SenderID:       "user1",
		ReceiverID:     "user2",
		Content:        "To be deleted",
		MessageType:    "text",
	}
	err = repo.CreateMessage(ctx, message)
	require.NoError(t, err)

	// Act
	err = repo.DeleteMessage(ctx, message.ID.Hex())

	// Assert
	require.NoError(t, err)

	// 验证消息已被软删除
	found, err := repo.GetMessageByID(ctx, message.ID.Hex())
	require.NoError(t, err)
	assert.True(t, found.IsDeleted)
	assert.NotNil(t, found.DeletedAt)
}

// TestMongoMessageRepository_Health 测试健康检查
func TestMongoMessageRepository_Health(t *testing.T) {
	// Arrange
	repo, _, ctx, cleanup := setupMessageRepo(t)
	defer cleanup()

	// Act
	err := repo.Health(ctx)

	// Assert
	assert.NoError(t, err)
}
