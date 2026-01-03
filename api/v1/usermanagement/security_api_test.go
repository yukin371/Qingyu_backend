package usermanagement

import (
	user2 "Qingyu_backend/service/interfaces/user"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService 模拟用户服务
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) SendEmailVerification(ctx context.Context, req *user2.SendEmailVerificationRequest) (*user2.SendEmailVerificationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.SendEmailVerificationResponse), args.Error(1)
}

func (m *MockUserService) VerifyEmail(ctx context.Context, req *user2.VerifyEmailRequest) (*user2.VerifyEmailResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.VerifyEmailResponse), args.Error(1)
}

func (m *MockUserService) RequestPasswordReset(ctx context.Context, req *user2.RequestPasswordResetRequest) (*user2.RequestPasswordResetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.RequestPasswordResetResponse), args.Error(1)
}

func (m *MockUserService) ConfirmPasswordReset(ctx context.Context, req *user2.ConfirmPasswordResetRequest) (*user2.ConfirmPasswordResetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user2.ConfirmPasswordResetResponse), args.Error(1)
}

// 辅助函数：创建测试上下文
func setupTestContext(method, path string, body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	var reqBody *bytes.Buffer
	if body != nil {
		jsonData, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(jsonData)
	} else {
		reqBody = bytes.NewBuffer([]byte{})
	}

	c.Request = httptest.NewRequest(method, path, reqBody)
	c.Request.Header.Set("Content-Type", "application/json")

	return c, w
}

// TestSecurityAPI_SendEmailVerification 测试发送邮箱验证码API
func TestSecurityAPI_SendEmailVerification(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "成功发送验证码",
			requestBody: map[string]interface{}{
				"user_id": "user123",
				"email":   "test@example.com",
			},
			setupMock: func(m *MockUserService) {
				m.On("SendEmailVerification", mock.Anything, mock.AnythingOfType("*user.SendEmailVerificationRequest")).Return(
					&user2.SendEmailVerificationResponse{
						Success:   true,
						Message:   "验证码已发送到您的邮箱",
						ExpiresIn: 1800,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "验证码已发送到您的邮箱",
			},
		},
		{
			name: "请求参数错误_缺少user_id",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "请求参数错误_缺少email",
			requestBody: map[string]interface{}{
				"user_id": "user123",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "服务错误_用户不存在",
			requestBody: map[string]interface{}{
				"user_id": "nonexistent",
				"email":   "test@example.com",
			},
			setupMock: func(m *MockUserService) {
				m.On("SendEmailVerification", mock.Anything, mock.Anything).Return(
					nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code": float64(500),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			api := NewSecurityAPI(mockService)

			c, w := setupTestContext("POST", "/api/v1/user-management/email/send-code", tt.requestBody)

			// Act
			api.SendEmailVerification(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					actualValue, exists := response[key]
					assert.True(t, exists, "响应中应包含字段: %s", key)
					assert.Equal(t, expectedValue, actualValue, "字段 %s 的值不匹配", key)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestSecurityAPI_VerifyEmail 测试验证邮箱API
func TestSecurityAPI_VerifyEmail(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "验证成功",
			requestBody: map[string]interface{}{
				"user_id": "user123",
				"code":    "123456",
			},
			setupMock: func(m *MockUserService) {
				m.On("VerifyEmail", mock.Anything, mock.AnythingOfType("*user.VerifyEmailRequest")).Return(
					&user2.VerifyEmailResponse{
						Success: true,
						Message: "邮箱验证成功",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "邮箱验证成功",
			},
		},
		{
			name: "验证失败_验证码错误",
			requestBody: map[string]interface{}{
				"user_id": "user123",
				"code":    "000000",
			},
			setupMock: func(m *MockUserService) {
				m.On("VerifyEmail", mock.Anything, mock.Anything).Return(
					nil, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code": float64(400),
			},
		},
		{
			name: "请求参数错误_缺少user_id",
			requestBody: map[string]interface{}{
				"code": "123456",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "请求参数错误_缺少code",
			requestBody: map[string]interface{}{
				"user_id": "user123",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "请求参数错误_空JSON",
			requestBody: map[string]interface{}{},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			api := NewSecurityAPI(mockService)

			c, w := setupTestContext("POST", "/api/v1/user-management/email/verify", tt.requestBody)

			// Act
			api.VerifyEmail(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					actualValue, exists := response[key]
					assert.True(t, exists, "响应中应包含字段: %s", key)
					assert.Equal(t, expectedValue, actualValue, "字段 %s 的值不匹配", key)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestSecurityAPI_RequestPasswordReset 测试请求密码重置API
func TestSecurityAPI_RequestPasswordReset(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "成功发送重置邮件",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			setupMock: func(m *MockUserService) {
				m.On("RequestPasswordReset", mock.Anything, mock.AnythingOfType("*user.RequestPasswordResetRequest")).Return(
					&user2.RequestPasswordResetResponse{
						Success:   true,
						Message:   "密码重置邮件已发送",
						ExpiresIn: 3600,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "密码重置邮件已发送",
			},
		},
		{
			name: "请求参数错误_缺少email",
			requestBody: map[string]interface{}{
				"invalid": "field",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "请求参数错误_空JSON",
			requestBody: map[string]interface{}{},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "服务错误_数据库失败",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
			},
			setupMock: func(m *MockUserService) {
				m.On("RequestPasswordReset", mock.Anything, mock.Anything).Return(
					nil, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"code": float64(500),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			api := NewSecurityAPI(mockService)

			c, w := setupTestContext("POST", "/api/v1/user-management/password/reset-request", tt.requestBody)

			// Act
			api.RequestPasswordReset(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					actualValue, exists := response[key]
					assert.True(t, exists, "响应中应包含字段: %s", key)
					assert.Equal(t, expectedValue, actualValue, "字段 %s 的值不匹配", key)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestSecurityAPI_ConfirmPasswordReset 测试确认密码重置API
func TestSecurityAPI_ConfirmPasswordReset(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockUserService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "成功重置密码",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"token":    "validtoken123",
				"password": "NewPassword123!",
			},
			setupMock: func(m *MockUserService) {
				m.On("ConfirmPasswordReset", mock.Anything, mock.AnythingOfType("*user.ConfirmPasswordResetRequest")).Return(
					&user2.ConfirmPasswordResetResponse{
						Success: true,
						Message: "密码重置成功，请使用新密码登录",
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"code":    float64(0),
				"message": "密码重置成功，请使用新密码登录",
			},
		},
		{
			name: "验证失败_Token无效",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"token":    "invalidtoken",
				"password": "NewPassword123!",
			},
			setupMock: func(m *MockUserService) {
				m.On("ConfirmPasswordReset", mock.Anything, mock.Anything).Return(
					nil, assert.AnError)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"code": float64(400),
			},
		},
		{
			name: "请求参数错误_缺少email",
			requestBody: map[string]interface{}{
				"token":   "validtoken123",
				"password": "NewPassword123!",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "请求参数错误_缺少token",
			requestBody: map[string]interface{}{
				"email":    "test@example.com",
				"password": "NewPassword123!",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "请求参数错误_缺少password",
			requestBody: map[string]interface{}{
				"email": "test@example.com",
				"token": "validtoken123",
			},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "请求参数错误_空JSON",
			requestBody:    map[string]interface{}{},
			setupMock:      func(m *MockUserService) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			tt.setupMock(mockService)
			api := NewSecurityAPI(mockService)

			c, w := setupTestContext("POST", "/api/v1/user-management/password/reset", tt.requestBody)

			// Act
			api.ConfirmPasswordReset(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					actualValue, exists := response[key]
					assert.True(t, exists, "响应中应包含字段: %s", key)
					assert.Equal(t, expectedValue, actualValue, "字段 %s 的值不匹配", key)
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

// TestSecurityAPI_InvalidJSON 测试无效的JSON
func TestSecurityAPI_InvalidJSON(t *testing.T) {
	tests := []struct {
		name       string
		apiMethod  func(*SecurityAPI, *gin.Context)
		endpoint   string
		invalidJSON string
	}{
		{
			name:       "SendEmailVerification_无效JSON",
			apiMethod:  func(api *SecurityAPI, c *gin.Context) { api.SendEmailVerification(c) },
			endpoint:   "/api/v1/user-management/email/send-code",
			invalidJSON: "{invalid json",
		},
		{
			name:       "VerifyEmail_无效JSON",
			apiMethod:  func(api *SecurityAPI, c *gin.Context) { api.VerifyEmail(c) },
			endpoint:   "/api/v1/user-management/email/verify",
			invalidJSON: "{invalid json",
		},
		{
			name:       "RequestPasswordReset_无效JSON",
			apiMethod:  func(api *SecurityAPI, c *gin.Context) { api.RequestPasswordReset(c) },
			endpoint:   "/api/v1/user-management/password/reset-request",
			invalidJSON: "{invalid json",
		},
		{
			name:       "ConfirmPasswordReset_无效JSON",
			apiMethod:  func(api *SecurityAPI, c *gin.Context) { api.ConfirmPasswordReset(c) },
			endpoint:   "/api/v1/user-management/password/reset",
			invalidJSON: "{invalid json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockService := new(MockUserService)
			api := NewSecurityAPI(mockService)

			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			reqBody := bytes.NewBufferString(tt.invalidJSON)
			c.Request = httptest.NewRequest("POST", tt.endpoint, reqBody)
			c.Request.Header.Set("Content-Type", "application/json")

			// Act
			tt.apiMethod(api, c)

			// Assert
			assert.Equal(t, http.StatusBadRequest, w.Code)
		})
	}
}

// TestNewSecurityAPI 测试创建SecurityAPI实例
func TestNewSecurityAPI(t *testing.T) {
	// Arrange
	mockService := new(MockUserService)

	// Act
	api := NewSecurityAPI(mockService)

	// Assert
	assert.NotNil(t, api)
	assert.NotNil(t, api.userService)
}

// Benchmark_SendEmailVerificationAPI 性能测试
func Benchmark_SendEmailVerificationAPI(b *testing.B) {
	mockService := new(MockUserService)
	mockService.On("SendEmailVerification", mock.Anything, mock.Anything).Return(
		&user2.SendEmailVerificationResponse{
			Success:   true,
			Message:   "验证码已发送到您的邮箱",
			ExpiresIn: 1800,
		}, nil)

	api := NewSecurityAPI(mockService)

	requestBody := map[string]interface{}{
		"user_id": "user123",
		"email":   "test@example.com",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c, _ := setupTestContext("POST", "/api/v1/user-management/email/send-code", requestBody)
		api.SendEmailVerification(c)
	}
}
