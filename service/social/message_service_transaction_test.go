package social

import (
	"context"
	"errors"
	"testing"
	"time"

	socialModel "Qingyu_backend/models/social"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type messageRepoState struct {
	conversations         map[string]*socialModel.Conversation
	messages              map[string]*socialModel.Message
	failUpdateLastMessage error
	failIncrementUnread   error
	failClearUnread       error
}

func newMessageRepoState() *messageRepoState {
	return &messageRepoState{
		conversations: make(map[string]*socialModel.Conversation),
		messages:      make(map[string]*socialModel.Message),
	}
}

func (m *messageRepoState) CreateConversation(ctx context.Context, conversation *socialModel.Conversation) error {
	if conversation.ID.IsZero() {
		conversation.ID = primitive.NewObjectID()
	}
	cloned := *conversation
	cloned.Participants = append([]string(nil), conversation.Participants...)
	cloned.UnreadCount = cloneUnreadMap(conversation.UnreadCount)
	m.conversations[conversation.ID.Hex()] = &cloned
	return nil
}
func (m *messageRepoState) GetConversationByID(ctx context.Context, conversationID string) (*socialModel.Conversation, error) {
	conversation, ok := m.conversations[conversationID]
	if !ok {
		return nil, errors.New("conversation not found")
	}
	return cloneConversation(conversation), nil
}
func (m *messageRepoState) GetConversationByParticipants(ctx context.Context, participantIDs []string) (*socialModel.Conversation, error) {
	for _, conversation := range m.conversations {
		if sameParticipants(conversation.Participants, participantIDs) {
			return cloneConversation(conversation), nil
		}
	}
	return nil, errors.New("conversation not found")
}
func (m *messageRepoState) GetUserConversations(ctx context.Context, userID string, page, size int) ([]*socialModel.Conversation, int64, error) {
	return nil, 0, nil
}
func (m *messageRepoState) UpdateConversation(ctx context.Context, conversationID string, updates map[string]interface{}) error {
	return nil
}
func (m *messageRepoState) DeleteConversation(ctx context.Context, conversationID string) error {
	return nil
}
func (m *messageRepoState) UpdateLastMessage(ctx context.Context, conversationID string, message *socialModel.Message) error {
	if m.failUpdateLastMessage != nil {
		return m.failUpdateLastMessage
	}
	conversation, ok := m.conversations[conversationID]
	if !ok {
		return errors.New("conversation not found")
	}
	copyMsg := *message
	conversation.LastMessage = &copyMsg
	conversation.UpdatedAt = time.Now()
	return nil
}
func (m *messageRepoState) IncrementUnreadCount(ctx context.Context, conversationID, userID string) error {
	if m.failIncrementUnread != nil {
		return m.failIncrementUnread
	}
	conversation, ok := m.conversations[conversationID]
	if !ok {
		return errors.New("conversation not found")
	}
	if conversation.UnreadCount == nil {
		conversation.UnreadCount = make(map[string]int)
	}
	conversation.UnreadCount[userID]++
	return nil
}
func (m *messageRepoState) ClearUnreadCount(ctx context.Context, conversationID, userID string) error {
	if m.failClearUnread != nil {
		return m.failClearUnread
	}
	conversation, ok := m.conversations[conversationID]
	if !ok {
		return errors.New("conversation not found")
	}
	if conversation.UnreadCount == nil {
		conversation.UnreadCount = make(map[string]int)
	}
	conversation.UnreadCount[userID] = 0
	return nil
}
func (m *messageRepoState) CreateMessage(ctx context.Context, message *socialModel.Message) error {
	if message.ID.IsZero() {
		message.ID = primitive.NewObjectID()
	}
	copyMsg := *message
	m.messages[message.ID.Hex()] = &copyMsg
	return nil
}
func (m *messageRepoState) GetMessageByID(ctx context.Context, messageID string) (*socialModel.Message, error) {
	message, ok := m.messages[messageID]
	if !ok {
		return nil, errors.New("message not found")
	}
	copyMsg := *message
	return &copyMsg, nil
}
func (m *messageRepoState) GetMessagesByConversation(ctx context.Context, conversationID string, page, size int) ([]*socialModel.Message, int64, error) {
	return nil, 0, nil
}
func (m *messageRepoState) GetMessagesBetweenUsers(ctx context.Context, userID1, userID2 string, page, size int) ([]*socialModel.Message, int64, error) {
	return nil, 0, nil
}
func (m *messageRepoState) MarkMessageAsRead(ctx context.Context, messageID string) error {
	message, ok := m.messages[messageID]
	if !ok {
		return errors.New("message not found")
	}
	message.IsRead = true
	now := time.Now()
	message.ReadAt = &now
	message.UpdatedAt = now
	return nil
}
func (m *messageRepoState) MarkConversationMessagesAsRead(ctx context.Context, conversationID, userID string) error {
	return nil
}
func (m *messageRepoState) DeleteMessage(ctx context.Context, messageID string) error { return nil }
func (m *messageRepoState) CountUnreadMessages(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
func (m *messageRepoState) CreateMention(ctx context.Context, mention *socialModel.Mention) error {
	return nil
}
func (m *messageRepoState) GetMentionByID(ctx context.Context, mentionID string) (*socialModel.Mention, error) {
	return nil, nil
}
func (m *messageRepoState) GetUserMentions(ctx context.Context, userID string, page, size int) ([]*socialModel.Mention, int64, error) {
	return nil, 0, nil
}
func (m *messageRepoState) GetUnreadMentions(ctx context.Context, userID string, page, size int) ([]*socialModel.Mention, int64, error) {
	return nil, 0, nil
}
func (m *messageRepoState) MarkMentionAsRead(ctx context.Context, mentionID string) error { return nil }
func (m *messageRepoState) MarkAllMentionsAsRead(ctx context.Context, userID string) error {
	return nil
}
func (m *messageRepoState) CountUnreadMentions(ctx context.Context, userID string) (int, error) {
	return 0, nil
}
func (m *messageRepoState) DeleteMention(ctx context.Context, mentionID string) error { return nil }
func (m *messageRepoState) RunInTransaction(ctx context.Context, fn func(context.Context) error) error {
	conversationSnapshot := cloneConversationMap(m.conversations)
	messageSnapshot := cloneMessageMap(m.messages)
	if err := fn(ctx); err != nil {
		m.conversations = conversationSnapshot
		m.messages = messageSnapshot
		return err
	}
	return nil
}
func (m *messageRepoState) Health(ctx context.Context) error { return nil }

func TestSendMessageRollbackOnUnreadCountFailure(t *testing.T) {
	repo := newMessageRepoState()
	conversation := &socialModel.Conversation{
		ID:           primitive.NewObjectID(),
		Participants: []string{"sender-1", "receiver-1"},
		UnreadCount:  map[string]int{"receiver-1": 0},
	}
	repo.conversations[conversation.ID.Hex()] = cloneConversation(conversation)
	repo.failIncrementUnread = errors.New("mock unread failure")

	service := NewMessageService(repo, nil)

	message, err := service.SendMessage(context.Background(), "sender-1", "receiver-1", "hello", "text")
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Empty(t, repo.messages)
	assert.Nil(t, repo.conversations[conversation.ID.Hex()].LastMessage)
	assert.Equal(t, 0, repo.conversations[conversation.ID.Hex()].UnreadCount["receiver-1"])
}

func TestMarkMessageAsReadRollbackOnUnreadClearFailure(t *testing.T) {
	repo := newMessageRepoState()
	conversation := &socialModel.Conversation{
		ID:           primitive.NewObjectID(),
		Participants: []string{"sender-1", "receiver-1"},
		UnreadCount:  map[string]int{"receiver-1": 1},
	}
	repo.conversations[conversation.ID.Hex()] = cloneConversation(conversation)
	message := &socialModel.Message{
		ID:             primitive.NewObjectID(),
		ConversationID: conversation.ID.Hex(),
		SenderID:       "sender-1",
		ReceiverID:     "receiver-1",
		Content:        "hello",
		MessageType:    "text",
	}
	message.IsRead = false
	repo.messages[message.ID.Hex()] = message
	repo.failClearUnread = errors.New("mock clear unread failure")

	service := NewMessageService(repo, nil)

	err := service.MarkMessageAsRead(context.Background(), "receiver-1", message.ID.Hex())
	assert.Error(t, err)
	assert.False(t, repo.messages[message.ID.Hex()].IsRead)
	assert.Equal(t, 1, repo.conversations[conversation.ID.Hex()].UnreadCount["receiver-1"])
}

func sameParticipants(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	seen := make(map[string]int, len(left))
	for _, value := range left {
		seen[value]++
	}
	for _, value := range right {
		seen[value]--
		if seen[value] < 0 {
			return false
		}
	}
	return true
}

func cloneUnreadMap(source map[string]int) map[string]int {
	if source == nil {
		return nil
	}
	cloned := make(map[string]int, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}

func cloneConversation(conversation *socialModel.Conversation) *socialModel.Conversation {
	if conversation == nil {
		return nil
	}
	cloned := *conversation
	cloned.Participants = append([]string(nil), conversation.Participants...)
	cloned.UnreadCount = cloneUnreadMap(conversation.UnreadCount)
	if conversation.LastMessage != nil {
		lastMessage := *conversation.LastMessage
		cloned.LastMessage = &lastMessage
	}
	return &cloned
}

func cloneConversationMap(source map[string]*socialModel.Conversation) map[string]*socialModel.Conversation {
	cloned := make(map[string]*socialModel.Conversation, len(source))
	for key, value := range source {
		cloned[key] = cloneConversation(value)
	}
	return cloned
}

func cloneMessageMap(source map[string]*socialModel.Message) map[string]*socialModel.Message {
	cloned := make(map[string]*socialModel.Message, len(source))
	for key, value := range source {
		if value == nil {
			cloned[key] = nil
			continue
		}
		copyValue := *value
		cloned[key] = &copyValue
	}
	return cloned
}
