package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	aiModels "Qingyu_backend/models/ai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockChatRepository 模拟聊天仓库
type MockChatRepository struct {
	mock.Mock
}

func (m *MockChatRepository) CreateSession(ctx context.Context, session *aiModels.ChatSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockChatRepository) GetSessionByID(ctx context.Context, sessionID string) (*aiModels.ChatSession, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*aiModels.ChatSession), args.Error(1)
}

func (m *MockChatRepository) UpdateSession(ctx context.Context, session *aiModels.ChatSession) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *MockChatRepository) DeleteSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockChatRepository) GetSessionsByProjectID(ctx context.Context, projectID string, limit, offset int) ([]*aiModels.ChatSession, error) {
	args := m.Called(ctx, projectID, limit, offset)
	return args.Get(0).([]*aiModels.ChatSession), args.Error(1)
}

func (m *MockChatRepository) CreateMessage(ctx context.Context, message *aiModels.ChatMessage) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}

func (m *MockChatRepository) GetSessionStatistics(ctx context.Context, projectID string) (*ChatStatistics, error) {
	args := m.Called(ctx, projectID)
	return args.Get(0).(*ChatStatistics), args.Error(1)
}

// MockAIService 模拟AI服务
type MockAIService struct {
	mock.Mock
}

func (m *MockAIService) GenerateContent(ctx context.Context, req *GenerateContentRequest) (*GenerateContentResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*GenerateContentResponse), args.Error(1)
}

// setupTestChatService 设置测试聊天服务
func setupTestChatService() (*ChatService, *MockChatRepository, *MockAIService) {
	mockRepo := &MockChatRepository{}
	mockAI := &MockAIService{}
	
	service := &ChatService{
		repository:           mockRepo,
		aiService:           mockAI,
		novelContextService: nil, // 暂时设为nil
	}
	
	return service, mockRepo, mockAI
}

// TestChatService_StartChat_NewSession 测试开始新聊天会话
func TestChatService_StartChat_NewSession(t *testing.T) {
	service, mockRepo, mockAI := setupTestChatService()
	
	ctx := context.Background()
	req := &ChatRequest{
		ProjectID:  "test-project",
		Message:    "Hello, AI!",
		UseContext: false,
	}
	
	// Mock AI响应
	aiResponse := &GenerateContentResponse{
		Content:     "Hello! How can I help you?",
		TokensUsed:  50,
		Model:       "gpt-3.5-turbo",
		GeneratedAt: time.Now(),
	}
	
	// 设置mock期望
	mockRepo.On("GetSessionByID", ctx, mock.AnythingOfType("string")).Return((*aiModels.ChatSession)(nil), errors.New("session not found")).Maybe()
	mockRepo.On("CreateSession", ctx, mock.MatchedBy(func(session *aiModels.ChatSession) bool {
		return session != nil && session.ProjectID == "test-project"
	})).Return(nil)
	mockAI.On("GenerateContent", ctx, mock.MatchedBy(func(req *GenerateContentRequest) bool {
		return req != nil && req.ProjectID == "test-project"
	})).Return(aiResponse, nil)
	mockRepo.On("CreateMessage", ctx, mock.MatchedBy(func(message *aiModels.ChatMessage) bool {
		return message != nil && (message.Role == "user" || message.Role == "assistant")
	})).Return(nil).Twice() // 用户消息和AI消息
	mockRepo.On("UpdateSession", ctx, mock.MatchedBy(func(session *aiModels.ChatSession) bool {
		return session != nil && session.ProjectID == "test-project"
	})).Return(nil)
	
	// 执行测试
	response, err := service.StartChat(ctx, req)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.SessionID)
	assert.Equal(t, "Hello! How can I help you?", response.Message.Content)
	assert.Equal(t, "assistant", response.Message.Role)
	assert.Equal(t, 50, response.TokensUsed)
	assert.Equal(t, "gpt-3.5-turbo", response.Model)
	assert.False(t, response.ContextUsed)
	
	// 验证mock调用
	mockRepo.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

// TestChatService_StartChat_ExistingSession 测试继续现有聊天会话
func TestChatService_StartChat_ExistingSession(t *testing.T) {
	service, mockRepo, mockAI := setupTestChatService()
	
	ctx := context.Background()
	sessionID := primitive.NewObjectID()
	req := &ChatRequest{
		SessionID: sessionID.Hex(),
		ProjectID: "test-project",
		Message:   "Continue chat",
	}
	
	// 准备现有会话
	sessionObjID := primitive.NewObjectID()
	projectObjID := primitive.NewObjectID()
	existingSession := &aiModels.ChatSession{
		ID:        sessionObjID,
		SessionID: sessionID.Hex(), // 设置SessionID字段
		ProjectID: projectObjID.Hex(),
		Title:     "Existing Chat",
		Messages: []aiModels.ChatMessage{
			{
				ID:        primitive.NewObjectID(),
				Role:      "user",
				Content:   "Previous message",
				Timestamp: time.Now().Add(-time.Hour),
			},
		},
		Settings:  &aiModels.ChatSettings{Model: "gpt-3.5-turbo"},
		Metadata:  &aiModels.ChatMetadata{},
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Hour),
	}
	
	aiResponse := &GenerateContentResponse{
		Content:     "Continuing our conversation...",
		TokensUsed:  30,
		Model:       "gpt-3.5-turbo",
		GeneratedAt: time.Now(),
	}
	
	// 设置mock期望
	mockRepo.On("GetSessionByID", ctx, sessionID.Hex()).Return(existingSession, nil)
	mockAI.On("GenerateContent", ctx, mock.MatchedBy(func(req *GenerateContentRequest) bool {
		return req != nil && req.ProjectID == "test-project"
	})).Return(aiResponse, nil)
	mockRepo.On("CreateMessage", ctx, mock.MatchedBy(func(message *aiModels.ChatMessage) bool {
		return message != nil && (message.Role == "user" || message.Role == "assistant")
	})).Return(nil).Twice()
	mockRepo.On("UpdateSession", ctx, mock.MatchedBy(func(session *aiModels.ChatSession) bool {
		return session != nil
	})).Return(nil)
	
	// 执行测试
	response, err := service.StartChat(ctx, req)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, sessionID.Hex(), response.SessionID) // 使用请求中的SessionID
	assert.Equal(t, "Continuing our conversation...", response.Message.Content)
	assert.Equal(t, 30, response.TokensUsed) // 修正类型为int
	
	mockRepo.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

// TestChatService_StartChat_EmptyMessage 测试空消息
func TestChatService_StartChat_EmptyMessage(t *testing.T) {
	service, _, _ := setupTestChatService()
	
	ctx := context.Background()
	req := &ChatRequest{
		ProjectID: "test-project",
		Message:   "", // 空消息
	}
	
	response, err := service.StartChat(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "消息内容不能为空")
}

// TestChatService_StartChat_AIServiceError 测试AI服务错误
func TestChatService_StartChat_AIServiceError(t *testing.T) {
	service, mockRepo, mockAI := setupTestChatService()
	
	ctx := context.Background()
	req := &ChatRequest{
		ProjectID: "test-project",
		Message:   "Hello",
	}
	
	// 设置mock期望
	mockRepo.On("GetSessionByID", ctx, mock.AnythingOfType("string")).Return((*aiModels.ChatSession)(nil), errors.New("session not found")).Maybe()
	mockRepo.On("CreateSession", ctx, mock.MatchedBy(func(session *aiModels.ChatSession) bool {
		return session != nil && session.ProjectID == "test-project"
	})).Return(nil)
	mockRepo.On("CreateMessage", ctx, mock.MatchedBy(func(message *aiModels.ChatMessage) bool {
		return message != nil && message.Role == "user"
	})).Return(nil)
	mockAI.On("GenerateContent", ctx, mock.MatchedBy(func(req *GenerateContentRequest) bool {
		return req != nil && req.ProjectID == "test-project"
	})).Return((*GenerateContentResponse)(nil), errors.New("AI service error"))
	
	response, err := service.StartChat(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "AI service error")
	
	mockRepo.AssertExpectations(t)
	mockAI.AssertExpectations(t)
}

// TestChatService_GetChatHistory_Success 测试获取聊天历史成功
func TestChatService_GetChatHistory_Success(t *testing.T) {
	service, mockRepo, _ := setupTestChatService()
	
	ctx := context.Background()
	sessionID := "test-session-123"
	
	// 准备测试数据
	sessionObjID := primitive.NewObjectID()
	dbSession := &aiModels.ChatSession{
		ID:        sessionObjID,
		SessionID: sessionID, // 设置SessionID字段
		ProjectID: "test-project",
		Title:     "Test Chat",
		Messages: []aiModels.ChatMessage{
			{
				ID:        primitive.NewObjectID(),
				Role:      "user",
				Content:   "Hello",
				Timestamp: time.Now().Add(-time.Hour),
			},
			{
				ID:        primitive.NewObjectID(),
				Role:      "assistant",
				Content:   "Hi there!",
				Timestamp: time.Now().Add(-time.Minute * 59),
			},
		},
		Settings:  &aiModels.ChatSettings{Model: "gpt-3.5-turbo"},
		Metadata:  &aiModels.ChatMetadata{},
		CreatedAt: time.Now().Add(-time.Hour),
		UpdatedAt: time.Now().Add(-time.Minute * 59),
	}
	
	mockRepo.On("GetSessionByID", ctx, sessionID).Return(dbSession, nil)
	
	// 执行测试
	result, err := service.GetChatHistory(ctx, sessionID)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, sessionID, result.ID)
	assert.Equal(t, "test-project", result.ProjectID)
	assert.Equal(t, "Test Chat", result.Title)
	assert.Len(t, result.Messages, 2)
	assert.Equal(t, "user", result.Messages[0].Role)
	assert.Equal(t, "Hello", result.Messages[0].Content)
	assert.Equal(t, "assistant", result.Messages[1].Role)
	assert.Equal(t, "Hi there!", result.Messages[1].Content)
	
	mockRepo.AssertExpectations(t)
}

// TestChatService_GetChatHistory_NotFound 测试会话不存在
func TestChatService_GetChatHistory_NotFound(t *testing.T) {
	service, mockRepo, _ := setupTestChatService()
	
	ctx := context.Background()
	sessionID := "non-existent-session"
	
	mockRepo.On("GetSessionByID", ctx, sessionID).Return((*aiModels.ChatSession)(nil), errors.New("session not found"))
	
	result, err := service.GetChatHistory(ctx, sessionID)
	
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "session not found")
	
	mockRepo.AssertExpectations(t)
}

// TestChatService_ListChatSessions_Success 测试列出聊天会话成功
func TestChatService_ListChatSessions_Success(t *testing.T) {
	service, mockRepo, _ := setupTestChatService()
	
	ctx := context.Background()
	projectID := "test-project"
	limit := 10
	offset := 0
	
	// 准备测试数据
	dbSessions := []*aiModels.ChatSession{
		{
			ID:        primitive.NewObjectID(),
			SessionID: "session-1", // 设置SessionID字段
			ProjectID: projectID,
			Title:     "Chat 1",
			Messages:  []aiModels.ChatMessage{},
			Settings:  &aiModels.ChatSettings{},
			Metadata:  &aiModels.ChatMetadata{},
			CreatedAt: time.Now().Add(-time.Hour * 2),
			UpdatedAt: time.Now().Add(-time.Hour * 2),
		},
		{
			ID:        primitive.NewObjectID(),
			SessionID: "session-2", // 设置SessionID字段
			ProjectID: projectID,
			Title:     "Chat 2",
			Messages:  []aiModels.ChatMessage{},
			Settings:  &aiModels.ChatSettings{},
			Metadata:  &aiModels.ChatMetadata{},
			CreatedAt: time.Now().Add(-time.Hour),
			UpdatedAt: time.Now().Add(-time.Hour),
		},
	}
	
	mockRepo.On("GetSessionsByProjectID", ctx, projectID, limit, offset).Return(dbSessions, nil)
	
	// 执行测试
	result, err := service.ListChatSessions(ctx, projectID, limit, offset)
	
	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	assert.Equal(t, "session-1", result[0].ID)
	assert.Equal(t, "Chat 1", result[0].Title)
	assert.Equal(t, "session-2", result[1].ID)
	assert.Equal(t, "Chat 2", result[1].Title)
	
	mockRepo.AssertExpectations(t)
}

// TestChatService_ListChatSessions_Empty 测试空会话列表
func TestChatService_ListChatSessions_Empty(t *testing.T) {
	service, mockRepo, _ := setupTestChatService()
	
	ctx := context.Background()
	projectID := "empty-project"
	
	mockRepo.On("GetSessionsByProjectID", ctx, projectID, 20, 0).Return([]*aiModels.ChatSession{}, nil)
	
	result, err := service.ListChatSessions(ctx, projectID, 20, 0)
	
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 0)
	
	mockRepo.AssertExpectations(t)
}

// TestChatService_DeleteChatSession_Success 测试删除聊天会话成功
func TestChatService_DeleteChatSession_Success(t *testing.T) {
	service, mockRepo, _ := setupTestChatService()
	
	ctx := context.Background()
	sessionID := "test-session-123"
	
	mockRepo.On("DeleteSession", ctx, sessionID).Return(nil)
	
	err := service.DeleteChatSession(ctx, sessionID)
	
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestChatService_DeleteChatSession_NotFound 测试删除不存在的会话
func TestChatService_DeleteChatSession_NotFound(t *testing.T) {
	service, mockRepo, _ := setupTestChatService()
	
	ctx := context.Background()
	sessionID := "non-existent-session"
	
	mockRepo.On("DeleteSession", ctx, sessionID).Return(errors.New("session not found"))
	
	err := service.DeleteChatSession(ctx, sessionID)
	
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
	mockRepo.AssertExpectations(t)
}