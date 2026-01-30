package messages_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	messageAPI "Qingyu_backend/api/v1/messages"
	"Qingyu_backend/models/social"
	"Qingyu_backend/service/interfaces"
)

// MockMessageService 模拟私信服务接口
type MockMessageService struct {
	mock.Mock
}

// GetConversations 模拟获取会话列表
func (m *MockMessageService) GetConversations(ctx context.Context, userID string, page, size int) ([]*social.Conversation, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Conversation), args.Get(1).(int64), args.Error(2)
}

// GetConversationMessages 模拟获取会话消息
func (m *MockMessageService) GetConversationMessages(ctx context.Context, userID, conversationID string, page, size int) ([]*social.Message, int64, error) {
	args := m.Called(ctx, userID, conversationID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Message), args.Get(1).(int64), args.Error(2)
}

// SendMessage 模拟发送消息
func (m *MockMessageService) SendMessage(ctx context.Context, senderID, receiverID, content, messageType string) (*social.Message, error) {
	args := m.Called(ctx, senderID, receiverID, content, messageType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*social.Message), args.Error(1)
}

// MarkMessageAsRead 模拟标记消息已读
func (m *MockMessageService) MarkMessageAsRead(ctx context.Context, userID, messageID string) error {
	args := m.Called(ctx, userID, messageID)
	return args.Error(0)
}

// DeleteMessage 模拟删除消息
func (m *MockMessageService) DeleteMessage(ctx context.Context, userID, messageID string) error {
	args := m.Called(ctx, userID, messageID)
	return args.Error(0)
}

// CreateMention 模拟创建@提醒
func (m *MockMessageService) CreateMention(ctx context.Context, senderID, userID, contentType, contentID, content string) error {
	args := m.Called(ctx, senderID, userID, contentType, contentID, content)
	return args.Error(0)
}

// GetMentions 模拟获取@提醒列表
func (m *MockMessageService) GetMentions(ctx context.Context, userID string, page, size int) ([]*social.Mention, int64, error) {
	args := m.Called(ctx, userID, page, size)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*social.Mention), args.Get(1).(int64), args.Error(2)
}

// MarkMentionAsRead 模拟标记@提醒已读
func (m *MockMessageService) MarkMentionAsRead(ctx context.Context, userID, mentionID string) error {
	args := m.Called(ctx, userID, mentionID)
	return args.Error(0)
}

// setupMessageTestRouter 设置测试路由
func setupMessageTestRouter(messageService interfaces.MessageService, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userId
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	api := messageAPI.NewMessageAPI(messageService)

	v1 := r.Group("/api/v1/social")
	{
		// 会话管理
		v1.GET("/messages/conversations", api.GetConversations)
		v1.GET("/messages/:conversationId", api.GetConversationMessages)

		// 消息管理
		v1.POST("/messages", api.SendMessage)
		v1.PUT("/messages/:id/read", api.MarkMessageAsRead)
		v1.DELETE("/messages/:id", api.DeleteMessage)

		// @提醒
		v1.POST("/mentions", api.CreateMention)
		v1.GET("/mentions", api.GetMentions)
		v1.PUT("/mentions/:id/read", api.MarkMentionAsRead)
	}

	return r
}

// =========================
// 会话管理测试
// =========================

// TestMessageAPI_GetConversations_Success 测试获取会话列表成功
func TestMessageAPI_GetConversations_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	expectedConversations := []*social.Conversation{
		{
			ID:           primitive.NewObjectID(),
			Participants: []string{userID, "user2"},
			UnreadCount:  map[string]int{userID: 2},
		},
	}

	mockService.On("GetConversations", mock.Anything, userID, 1, 20).
		Return(expectedConversations, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "获取会话列表成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_GetConversationMessages_Success 测试获取会话消息成功
func TestMessageAPI_GetConversationMessages_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	conversationID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	expectedMessages := []*social.Message{
		{
			ID:             primitive.NewObjectID(),
			ConversationID: conversationID,
			SenderID:       userID,
			ReceiverID:     "user2",
			Content:        "你好",
			MessageType:    "text",
		},
	}

	mockService.On("GetConversationMessages", mock.Anything, userID, conversationID, 1, 50).
		Return(expectedMessages, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/"+conversationID+"?page=1&size=50", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "获取消息成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_GetConversationMessages_EmptyConversationID 测试获取会话消息-会话ID为空
func TestMessageAPI_GetConversationMessages_EmptyConversationID(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/?page=1&size=50", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// =========================
// 消息管理测试
// =========================

// TestMessageAPI_SendMessage_Success 测试发送消息成功
func TestMessageAPI_SendMessage_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	receiverID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	expectedMessage := &social.Message{
		ID:          primitive.NewObjectID(),
		SenderID:    userID,
		ReceiverID:  receiverID,
		Content:     "你好，这是测试消息",
		MessageType: "text",
	}

	mockService.On("SendMessage", mock.Anything, userID, receiverID, "你好，这是测试消息", "text").
		Return(expectedMessage, nil)

	reqBody := map[string]interface{}{
		"receiver_id":  receiverID,
		"content":      "你好，这是测试消息",
		"message_type": "text",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "创建成功", response["message"])
	assert.NotNil(t, response["data"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_SendMessage_MissingReceiverID 测试发送消息-缺少接收者ID
func TestMessageAPI_SendMessage_MissingReceiverID(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"content":      "你好，这是测试消息",
		"message_type": "text",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_InvalidMessageType 测试发送消息-无效的消息类型
func TestMessageAPI_SendMessage_InvalidMessageType(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	receiverID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"receiver_id":  receiverID,
		"content":      "你好",
		"message_type": "invalid",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_MarkMessageAsRead_Success 测试标记消息已读成功
func TestMessageAPI_MarkMessageAsRead_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	messageID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	mockService.On("MarkMessageAsRead", mock.Anything, userID, messageID).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/v1/social/messages/"+messageID+"/read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "标记消息已读成功", response["message"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_DeleteMessage_Success 测试删除消息成功
func TestMessageAPI_DeleteMessage_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	messageID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	mockService.On("DeleteMessage", mock.Anything, userID, messageID).Return(nil)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/messages/"+messageID, nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "删除消息成功", response["message"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_DeleteMessage_EmptyMessageID 测试删除消息-消息ID为空
func TestMessageAPI_DeleteMessage_EmptyMessageID(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	req, _ := http.NewRequest("DELETE", "/api/v1/social/messages/", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// =========================
// @提醒测试
// =========================

// TestMessageAPI_CreateMention_Success 测试创建@提醒成功
func TestMessageAPI_CreateMention_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	targetUserID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	mockService.On("CreateMention", mock.Anything, userID, targetUserID, "comment", "comment123", "测试内容").
		Return(nil)

	reqBody := map[string]interface{}{
		"user_id":      targetUserID,
		"content_type": "comment",
		"content_id":   "comment123",
		"content":      "测试内容",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/mentions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "创建成功", response["message"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_CreateMention_MissingUserID 测试创建@提醒-缺少用户ID
func TestMessageAPI_CreateMention_MissingUserID(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	reqBody := map[string]interface{}{
		"content_type": "comment",
		"content_id":   "comment123",
		"content":      "测试内容",
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/v1/social/mentions", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_GetMentions_Success 测试获取@提醒列表成功
func TestMessageAPI_GetMentions_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	expectedMentions := []*social.Mention{
		{
			ID:          primitive.NewObjectID(),
			UserID:      userID,
			SenderID:    "user2",
			ContentType: "comment",
			ContentID:   "comment123",
			Content:     "测试内容",
		},
	}

	mockService.On("GetMentions", mock.Anything, userID, 1, 20).
		Return(expectedMentions, int64(1), nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/mentions?page=1&size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "获取@提醒列表成功", response["message"])

	data := response["data"].(map[string]interface{})
	assert.Equal(t, float64(1), data["total"])
	assert.NotNil(t, data["list"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_MarkMentionAsRead_Success 测试标记@提醒已读成功
func TestMessageAPI_MarkMentionAsRead_Success(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	mentionID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	mockService.On("MarkMentionAsRead", mock.Anything, userID, mentionID).Return(nil)

	req, _ := http.NewRequest("PUT", "/api/v1/social/mentions/"+mentionID+"/read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "标记@提醒已读成功", response["message"])

	mockService.AssertExpectations(t)
}

// TestMessageAPI_MarkMentionAsRead_EmptyMentionID 测试标记@提醒已读-提醒ID为空
func TestMessageAPI_MarkMentionAsRead_EmptyMentionID(t *testing.T) {
	// Given
	mockService := new(MockMessageService)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageTestRouter(mockService, userID)

	req, _ := http.NewRequest("PUT", "/api/v1/social/mentions//read", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
