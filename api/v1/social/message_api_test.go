package social_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	messageAPI "Qingyu_backend/api/v1/social"
	"Qingyu_backend/api/v1/social/dto"
	modelsMessaging "Qingyu_backend/models/messaging"
	"Qingyu_backend/models/messaging/mocks"
	serviceMessaging "Qingyu_backend/service/messaging"
)

// setupMessageTestRouter 设置消息测试路由
func setupMessageAPITestRouter(messageRepo *mocks.MockMessageRepository, conversationRepo *mocks.MockConversationRepository, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// 添加middleware来设置userID
	r.Use(func(c *gin.Context) {
		if userID != "" {
			c.Set("userID", userID)
		}
		c.Next()
	})

	// 创建真实的服务实例
	conversationSvc := serviceMessaging.NewConversationService(conversationRepo, messageRepo)
	messageSvc := serviceMessaging.NewMessageService(messageRepo, conversationSvc)

	api := messageAPI.NewMessageAPIV2(messageSvc, conversationSvc, nil)

	v1 := r.Group("/api/v1/social/messages")
	{
		v1.GET("/conversations/:conversationId/messages", api.GetMessages)
		v1.POST("/conversations/:conversationId/messages", api.SendMessage)
		v1.POST("/conversations", api.CreateConversation)
		v1.POST("/conversations/:conversationId/read", api.MarkConversationRead)
	}

	return r
}

// =========================
// GetMessages 测试
// =========================

// TestMessageAPI_GetMessages_MissingUserID 测试缺少用户ID
func TestMessageAPI_GetMessages_MissingUserID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, "")

	conversationID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=1&page_size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestMessageAPI_GetMessages_MissingConversationID 测试缺少会话ID
func TestMessageAPI_GetMessages_MissingConversationID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations//messages?page=1&page_size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_GetMessages_DefaultPage 测试默认page参数
func TestMessageAPI_GetMessages_DefaultPage(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{userID, otherUserID}

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(conv, nil)

	// 正确的签名：7个参数（包括userID和before/after）
	mockMessageRepo.On("FindByConversationID", mock.Anything, conversationID, userID, 1, 20, (*string)(nil), (*string)(nil)).
		Return([]*modelsMessaging.DirectMessage{}, 0, nil)

	// page=0应该被设置为默认值1
	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=0&page_size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockConvRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

// TestMessageAPI_GetMessages_DefaultPageSize 测试默认page_size参数
func TestMessageAPI_GetMessages_DefaultPageSize(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{userID, otherUserID}

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(conv, nil)

	mockMessageRepo.On("FindByConversationID", mock.Anything, conversationID, userID, 1, 20, (*string)(nil), (*string)(nil)).
		Return([]*modelsMessaging.DirectMessage{}, 0, nil)

	// page_size=0应该被设置为默认值20
	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=1&page_size=0", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockConvRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

// TestMessageAPI_GetMessages_PageSizeExceedsMax 测试page_size超过最大值
func TestMessageAPI_GetMessages_PageSizeExceedsMax(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=1&page_size=101", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_GetMessages_ConversationNotFound 测试会话不存在
func TestMessageAPI_GetMessages_ConversationNotFound(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(nil, modelsMessaging.ErrConversationNotFound)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=1&page_size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockConvRepo.AssertExpectations(t)
}

// TestMessageAPI_GetMessages_Forbidden 测试无权访问
func TestMessageAPI_GetMessages_Forbidden(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{otherUserID, primitive.NewObjectID().Hex()}

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(conv, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=1&page_size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockConvRepo.AssertExpectations(t)
}

// TestMessageAPI_GetMessages_Success 测试获取成功
func TestMessageAPI_GetMessages_Success(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{userID, otherUserID}

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(conv, nil)

	msgID := primitive.NewObjectID()
	messages := []*modelsMessaging.DirectMessage{
		func() *modelsMessaging.DirectMessage {
			m := &modelsMessaging.DirectMessage{}
			m.ID = msgID
			m.ConversationID = conversationID
			m.SenderID = otherUserID
			m.ReceiverID = userID
			m.Content = "Hello"
			m.Type = modelsMessaging.MessageTypeText
			m.IsRead = false
			m.CreatedAt = time.Now()
			return m
		}(),
	}

	mockMessageRepo.On("FindByConversationID", mock.Anything, conversationID, userID, 1, 20, (*string)(nil), (*string)(nil)).
		Return(messages, 1, nil)

	req, _ := http.NewRequest("GET", "/api/v1/social/messages/conversations/"+conversationID+"/messages?page=1&page_size=20", nil)

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockConvRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

// =========================
// SendMessage 测试
// =========================

// TestMessageAPI_SendMessage_MissingUserID 测试缺少用户ID
func TestMessageAPI_SendMessage_MissingUserID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, "")

	conversationID := primitive.NewObjectID().Hex()
	body := dto.SendMessageRequest{
		Content: "Hello",
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestMessageAPI_SendMessage_MissingConversationID 测试缺少会话ID
func TestMessageAPI_SendMessage_MissingConversationID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	body := dto.SendMessageRequest{
		Content: "Hello",
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations//messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_MissingContent 测试缺少content
func TestMessageAPI_SendMessage_MissingContent(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	body := dto.SendMessageRequest{
		Type: "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_ContentTooShort 测试content太短
func TestMessageAPI_SendMessage_ContentTooShort(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	body := dto.SendMessageRequest{
		Content: "",
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_ContentTooLong 测试content太长
func TestMessageAPI_SendMessage_ContentTooLong(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	longContent := string(make([]byte, 5001)) // 超过5000字符限制
	body := dto.SendMessageRequest{
		Content: longContent,
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_MissingType 测试缺少type
func TestMessageAPI_SendMessage_MissingType(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	body := map[string]string{
		"content": "Hello",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_InvalidType 测试无效的type
func TestMessageAPI_SendMessage_InvalidType(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	body := map[string]string{
		"content": "Hello",
		"type":    "invalid",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_SendMessage_ConversationNotFound 测试会话不存在
func TestMessageAPI_SendMessage_ConversationNotFound(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	body := dto.SendMessageRequest{
		Content: "Hello",
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(nil, modelsMessaging.ErrConversationNotFound)

	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	mockConvRepo.AssertExpectations(t)
}

// TestMessageAPI_SendMessage_Forbidden 测试无权访问
func TestMessageAPI_SendMessage_Forbidden(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{otherUserID, primitive.NewObjectID().Hex()}

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(conv, nil)

	body := dto.SendMessageRequest{
		Content: "Hello",
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)
	mockConvRepo.AssertExpectations(t)
}

// TestMessageAPI_SendMessage_Success 测试发送成功
func TestMessageAPI_SendMessage_Success(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{userID, otherUserID}

	mockConvRepo.On("FindByID", mock.Anything, conversationID).
		Return(conv, nil)

	msgID := primitive.NewObjectID()
	message := &modelsMessaging.DirectMessage{}
	message.ID = msgID
	message.ConversationID = conversationID
	message.SenderID = userID
	message.ReceiverID = otherUserID
	message.Content = "Hello"
	message.Type = modelsMessaging.MessageTypeText
	message.Status = modelsMessaging.MessageStatusNormal
	message.CreatedAt = time.Now()

	// Create方法在repository中只返回error，但service.Create会返回message
	// 所以我们需要mock repository.Create返回nil（成功），然后service会返回传入的message
	mockMessageRepo.On("Create", mock.Anything, mock.Anything).
		Return(nil)

	// UpdateLastMessage会调用conversationRepo.Update
	mockConvRepo.On("Update", mock.Anything, mock.AnythingOfType("*messaging.Conversation")).
		Return(nil)

	// IncrementUnreadCount会调用conversationRepo.IncrementUnreadCount
	mockConvRepo.On("IncrementUnreadCount", mock.Anything, conversationID, otherUserID).
		Return(nil)

	mockMessageRepo.On("CountUnreadInConversation", mock.Anything, conversationID, otherUserID).
		Return(5, nil)

	body := dto.SendMessageRequest{
		Content: "Hello",
		Type:    "text",
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/messages", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockConvRepo.AssertExpectations(t)
	mockMessageRepo.AssertExpectations(t)
}

// =========================
// CreateConversation 测试
// =========================

// TestMessageAPI_CreateConversation_MissingUserID 测试缺少用户ID
func TestMessageAPI_CreateConversation_MissingUserID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, "")

	body := dto.CreateConversationRequest{
		ParticipantIDs: []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestMessageAPI_CreateConversation_EmptyParticipantIDs 测试空的参与者ID列表
func TestMessageAPI_CreateConversation_EmptyParticipantIDs(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	body := dto.CreateConversationRequest{
		ParticipantIDs: []string{},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_CreateConversation_NotEnoughParticipants 测试参与者数量不足
func TestMessageAPI_CreateConversation_NotEnoughParticipants(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	body := dto.CreateConversationRequest{
		ParticipantIDs: []string{userID},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_CreateConversation_UserNotParticipant 测试用户不在参与者列表中
func TestMessageAPI_CreateConversation_UserNotParticipant(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	body := dto.CreateConversationRequest{
		ParticipantIDs: []string{primitive.NewObjectID().Hex(), primitive.NewObjectID().Hex()},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_CreateConversation_Success 测试创建成功
func TestMessageAPI_CreateConversation_Success(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	otherUserID := primitive.NewObjectID().Hex()

	conv := &modelsMessaging.Conversation{}
	conv.ID = primitive.NewObjectID()
	conv.ParticipantIDs = []string{userID, otherUserID}
	conv.CreatedAt = time.Now()
	conv.CreatedBy = userID

	// FindByParticipants会先检查是否已存在会话
	mockConvRepo.On("FindByParticipants", mock.Anything, []string{userID, otherUserID}).
		Return(nil, modelsMessaging.ErrConversationNotFound)

	mockConvRepo.On("Create", mock.Anything, mock.AnythingOfType("*messaging.Conversation")).
		Return(nil)

	body := dto.CreateConversationRequest{
		ParticipantIDs: []string{userID, otherUserID},
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusCreated, w.Code)
	mockConvRepo.AssertExpectations(t)
}

// =========================
// MarkConversationRead 测试
// =========================

// TestMessageAPI_MarkConversationRead_MissingUserID 测试缺少用户ID
func TestMessageAPI_MarkConversationRead_MissingUserID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, "")

	conversationID := primitive.NewObjectID().Hex()
	body := dto.MarkConversationReadRequest{
		ReadAt: time.Now().Unix(),
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// TestMessageAPI_MarkConversationRead_MissingConversationID 测试缺少会话ID
func TestMessageAPI_MarkConversationRead_MissingConversationID(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	body := dto.MarkConversationReadRequest{
		ReadAt: time.Now().Unix(),
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations//read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_MarkConversationRead_MissingReadAt 测试缺少read_at
func TestMessageAPI_MarkConversationRead_MissingReadAt(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	body := map[string]string{}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestMessageAPI_MarkConversationRead_Success 测试标记成功
func TestMessageAPI_MarkConversationRead_Success(t *testing.T) {
	// Given
	mockMessageRepo := new(mocks.MockMessageRepository)
	mockConvRepo := new(mocks.MockConversationRepository)
	userID := primitive.NewObjectID().Hex()
	router := setupMessageAPITestRouter(mockMessageRepo, mockConvRepo, userID)

	conversationID := primitive.NewObjectID().Hex()
	readAt := time.Now().Unix()

	mockMessageRepo.On("MarkConversationRead", mock.Anything, conversationID, userID, time.Unix(readAt, 0)).
		Return(10, nil)

	body := dto.MarkConversationReadRequest{
		ReadAt: readAt,
	}
	bodyBytes, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/v1/social/messages/conversations/"+conversationID+"/read", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// When
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	mockMessageRepo.AssertExpectations(t)
}
