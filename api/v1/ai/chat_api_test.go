package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	aiService "Qingyu_backend/service/ai"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockChatService 模拟聊天服务
type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) StartChat(ctx context.Context, req *aiService.ChatRequest) (*aiService.ChatResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*aiService.ChatResponse), args.Error(1)
}

func (m *MockChatService) GetChatHistory(ctx context.Context, sessionID string) (*aiService.ChatSession, error) {
	args := m.Called(ctx, sessionID)
	return args.Get(0).(*aiService.ChatSession), args.Error(1)
}

func (m *MockChatService) ListChatSessions(ctx context.Context, projectID string, limit, offset int) ([]*aiService.ChatSession, error) {
	args := m.Called(ctx, projectID, limit, offset)
	return args.Get(0).([]*aiService.ChatSession), args.Error(1)
}

func (m *MockChatService) DeleteChatSession(ctx context.Context, sessionID string) error {
	args := m.Called(ctx, sessionID)
	return args.Error(0)
}

func (m *MockChatService) GetStatistics(projectID string) (*aiService.ChatStatistics, error) {
	args := m.Called(projectID)
	return args.Get(0).(*aiService.ChatStatistics), args.Error(1)
}

func (m *MockChatService) StartChatStream(ctx context.Context, req *aiService.ChatRequest) (<-chan *aiService.StreamChatResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(<-chan *aiService.StreamChatResponse), args.Error(1)
}

// setupTestRouter 设置测试路由
func setupTestChatRouter() (*gin.Engine, *MockChatService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	mockService := &MockChatService{}
	chatApi := NewChatApi(mockService)
	
	v1 := router.Group("/api/v1")
	ai := v1.Group("/ai")
	{
		ai.POST("/chat/start", chatApi.StartChat)
		ai.GET("/chat/history/:sessionId", chatApi.GetChatHistory)
		ai.GET("/chat/sessions", chatApi.ListChatSessions)
		ai.DELETE("/chat/sessions/:sessionId", chatApi.DeleteChatSession)
		ai.GET("/chat/statistics", chatApi.GetChatStatistics)
	}
	
	return router, mockService
}

// TestChatApi_StartChat_Success 测试开始聊天成功
func TestChatApi_StartChat_Success(t *testing.T) {
	router, mockService := setupTestChatRouter()
	
	// 准备请求数据
	reqData := aiService.ChatRequest{
		ProjectID: "test-project",
		Message:   "Hello, AI!",
		UseContext: true,
	}
	
	// 准备响应数据
	expectedResponse := &aiService.ChatResponse{
		SessionID: "test-session-123",
		Message: &aiService.ChatMessage{
			ID:        "msg-123",
			Role:      "assistant",
			Content:   "Hello! How can I help you?",
			Timestamp: time.Now(),
		},
		TokensUsed:   50,
		Model:        "gpt-3.5-turbo",
		ContextUsed:  true,
		ResponseTime: time.Millisecond * 500,
	}
	
	// 设置mock期望
	mockService.On("StartChat", mock.Anything, &reqData).Return(expectedResponse, nil)
	
	// 创建请求
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat/start", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	// 执行请求
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// 验证结果
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["data"])
	
	// 验证mock调用
	mockService.AssertExpectations(t)
}

// TestChatApi_StartChat_BadRequest 测试开始聊天请求参数错误
func TestChatApi_StartChat_BadRequest(t *testing.T) {
	router, _ := setupTestChatRouter()
	
	// 测试无效JSON
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat/start", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "请求参数错误")
}

// TestChatApi_StartChat_EmptyMessage 测试空消息
func TestChatApi_StartChat_EmptyMessage(t *testing.T) {
	router, _ := setupTestChatRouter()
	
	reqData := aiService.ChatRequest{
		ProjectID: "test-project",
		Message:   "", // 空消息
	}
	
	jsonData, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat/start", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "Message")
}

// TestChatApi_GetChatHistory_Success 测试获取聊天历史成功
func TestChatApi_GetChatHistory_Success(t *testing.T) {
	router, mockService := setupTestChatRouter()
	
	sessionID := "test-session-123"
	expectedSession := &aiService.ChatSession{
		ID:        sessionID,
		ProjectID: "test-project",
		Title:     "Test Chat",
		Messages: []*aiService.ChatMessage{
			{
				ID:        "msg-1",
				Role:      "user",
				Content:   "Hello",
				Timestamp: time.Now(),
			},
			{
				ID:        "msg-2",
				Role:      "assistant",
				Content:   "Hi there!",
				Timestamp: time.Now(),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	mockService.On("GetChatHistory", mock.Anything, sessionID).Return(expectedSession, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/ai/chat/history/"+sessionID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	assert.NotNil(t, response["data"])
	
	mockService.AssertExpectations(t)
}

// TestChatApi_GetChatHistory_EmptySessionID 测试空会话ID
func TestChatApi_GetChatHistory_EmptySessionID(t *testing.T) {
	router, _ := setupTestChatRouter()
	
	req, _ := http.NewRequest("GET", "/api/v1/ai/chat//history", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// 由于路由不匹配，应该返回404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestChatApi_ListChatSessions_Success 测试列出聊天会话成功
func TestChatApi_ListChatSessions_Success(t *testing.T) {
	router, mockService := setupTestChatRouter()
	
	projectID := "test-project"
	expectedSessions := []*aiService.ChatSession{
		{
			ID:        "session-1",
			ProjectID: projectID,
			Title:     "Chat 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "session-2",
			ProjectID: projectID,
			Title:     "Chat 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	
	mockService.On("ListChatSessions", mock.Anything, projectID, 20, 0).Return(expectedSessions, nil)
	
	req, _ := http.NewRequest("GET", "/api/v1/ai/chat/sessions?projectId="+projectID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	
	data := response["data"].(map[string]interface{})
	assert.NotNil(t, data["sessions"])
	assert.Equal(t, float64(2), data["total"])
	assert.Equal(t, float64(20), data["limit"])
	assert.Equal(t, float64(0), data["offset"])
	
	mockService.AssertExpectations(t)
}

// TestChatApi_ListChatSessions_InvalidLimit 测试无效的limit参数
func TestChatApi_ListChatSessions_InvalidLimit(t *testing.T) {
	router, _ := setupTestChatRouter()
	
	req, _ := http.NewRequest("GET", "/api/v1/ai/chat/sessions?projectId=test&limit=invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "limit参数格式错误")
}

// TestChatApi_DeleteChatSession_Success 测试删除聊天会话成功
func TestChatApi_DeleteChatSession_Success(t *testing.T) {
	router, mockService := setupTestChatRouter()
	
	sessionID := "test-session-123"
	mockService.On("DeleteChatSession", mock.Anything, sessionID).Return(nil)
	
	req, _ := http.NewRequest("DELETE", "/api/v1/ai/chat/sessions/"+sessionID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(0), response["code"])
	assert.Equal(t, "success", response["message"])
	
	mockService.AssertExpectations(t)
}

// TestChatApi_DeleteChatSession_EmptySessionID 测试删除会话时空会话ID
func TestChatApi_DeleteChatSession_EmptySessionID(t *testing.T) {
	router, _ := setupTestChatRouter()
	
	req, _ := http.NewRequest("DELETE", "/api/v1/ai/chat/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	// 由于路由不匹配，应该返回404
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// TestChatApi_Integration 集成测试
func TestChatApi_Integration(t *testing.T) {
	router, mockService := setupTestChatRouter()
	
	// 1. 开始聊天
	startReq := aiService.ChatRequest{
		ProjectID: "test-project",
		Message:   "Hello",
	}
	
	startResp := &aiService.ChatResponse{
		SessionID: "session-123",
		Message: &aiService.ChatMessage{
			ID:      "msg-1",
			Role:    "assistant",
			Content: "Hello! How can I help?",
		},
		TokensUsed: 25,
		Model:      "gpt-3.5-turbo",
	}
	
	mockService.On("StartChat", mock.Anything, &startReq).Return(startResp, nil)
	
	jsonData, _ := json.Marshal(startReq)
	req, _ := http.NewRequest("POST", "/api/v1/ai/chat/start", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	// 2. 获取聊天历史
	session := &aiService.ChatSession{
		ID:        "session-123",
		ProjectID: "test-project",
		Messages: []*aiService.ChatMessage{
			{ID: "msg-1", Role: "assistant", Content: "Hello! How can I help?"},
		},
	}
	
	mockService.On("GetChatHistory", mock.Anything, "session-123").Return(session, nil)
	
	req2, _ := http.NewRequest("GET", "/api/v1/ai/chat/history/session-123", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	
	assert.Equal(t, http.StatusOK, w2.Code)
	
	// 3. 删除会话
	mockService.On("DeleteChatSession", mock.Anything, "session-123").Return(nil)
	
	req3, _ := http.NewRequest("DELETE", "/api/v1/ai/chat/sessions/session-123", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	
	assert.Equal(t, http.StatusOK, w3.Code)
	
	mockService.AssertExpectations(t)
}