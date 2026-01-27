package core

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// 测试辅助函数：创建测试请求上下文
func createTestRequestContext(t *testing.T) (*gin.Context, *RequestContext) {
	t.Helper()

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	gin.SetMode(gin.TestMode)
	
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	return c, NewRequestContext(c, logger)
}

// TestRequestContext_GetRequestID 测试获取请求ID
func TestRequestContext_GetRequestID(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 测试不存在的请求ID
	requestID := rc.GetRequestID()
	assert.Empty(t, requestID, "RequestID should be empty initially")

	// 设置请求ID
	c.Set(string(RequestIDKey), "test-request-id")
	requestID = rc.GetRequestID()
	assert.Equal(t, "test-request-id", requestID, "RequestID should be test-request-id")
}

// TestRequestContext_GetUserID 测试获取用户ID
func TestRequestContext_GetUserID(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 测试不存在的用户ID
	userID := rc.GetUserID()
	assert.Empty(t, userID, "UserID should be empty initially")

	// 设置用户ID
	c.Set(string(UserIDKey), "test-user-id")
	userID = rc.GetUserID()
	assert.Equal(t, "test-user-id", userID, "UserID should be test-user-id")
}

// TestRequestContext_GetUserRoles 测试获取用户角色
func TestRequestContext_GetUserRoles(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 测试不存在的角色
	roles := rc.GetUserRoles()
	assert.Nil(t, roles, "Roles should be nil initially")

	// 设置角色
	testRoles := []string{"admin", "user"}
	c.Set(string(UserRolesKey), testRoles)
	roles = rc.GetUserRoles()
	assert.Equal(t, testRoles, roles, "Roles should match")
}

// TestRequestContext_GetStartTime 测试获取开始时间
func TestRequestContext_GetStartTime(t *testing.T) {
	_, rc := createTestRequestContext(t)

	// 测试获取开始时间
	startTime := rc.GetStartTime()
	// 开始时间应该是当前时间附近
	assert.False(t, startTime.IsZero(), "StartTime should not be zero")
}

// TestRequestContext_SetError_GetError 测试设置和获取错误
func TestRequestContext_SetError_GetError(t *testing.T) {
	_, rc := createTestRequestContext(t)

	// 测试不存在的错误
	err := rc.GetError()
	assert.Nil(t, err, "Error should be nil initially")

	// 设置错误
	testError := assert.AnError
	rc.SetError(testError)
	err = rc.GetError()
	assert.Equal(t, testError, err, "Error should match")
}

// TestRequestContext_WithContext 测试创建带超时的context
func TestRequestContext_WithContext(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 创建一个真实的request
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	// 创建带超时的context
	ctx, cancel := rc.WithContext(5 * time.Second)
	defer cancel()

	assert.NotNil(t, ctx, "Context should not be nil")
	assert.NotNil(t, cancel, "Cancel function should not be nil")

	// 验证context有超时
	deadline, ok := ctx.Deadline()
	assert.True(t, ok, "Context should have deadline")
	assert.True(t, deadline.After(time.Now()), "Deadline should be in the future")
}

// TestRequestContext_LogInfo 测试记录信息日志
func TestRequestContext_LogInfo(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 创建一个真实的request
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	// 这个测试只是确保方法不会panic
	// 实际的日志验证需要更复杂的设置
	rc.LogInfo("test info message")
}

// TestRequestContext_LogError 测试记录错误日志
func TestRequestContext_LogError(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 创建一个真实的request
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	// 这个测试只是确保方法不会panic
	rc.LogError("test error message")
}

// TestRequestContext_LogWarn 测试记录警告日志
func TestRequestContext_LogWarn(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 创建一个真实的request
	req := httptest.NewRequest("GET", "/test", nil)
	c.Request = req

	// 这个测试只是确保方法不会panic
	rc.LogWarn("test warn message")
}

// TestRequestContext_InvalidTypeConversion 测试无效类型转换
func TestRequestContext_InvalidTypeConversion(t *testing.T) {
	c, rc := createTestRequestContext(t)

	// 设置错误类型的请求ID
	c.Set(string(RequestIDKey), 123) // 设置为int而不是string

	// 应该返回空字符串而不是panic
	requestID := rc.GetRequestID()
	assert.Empty(t, requestID, "RequestID should be empty for invalid type")

	// 设置错误类型的用户ID
	c.Set(string(UserIDKey), 456) // 设置为int而不是string

	// 应该返回空字符串而不是panic
	userID := rc.GetUserID()
	assert.Empty(t, userID, "UserID should be empty for invalid type")

	// 设置错误类型的角色
	c.Set(string(UserRolesKey), "invalid") // 设置为string而不是[]string

	// 应该返回nil而不是panic
	roles := rc.GetUserRoles()
	assert.Nil(t, roles, "Roles should be nil for invalid type")
}

// TestContextKey_ConstantValues 测试context键常量值
func TestContextKey_ConstantValues(t *testing.T) {
	// 验证常量值符合预期
	assert.Equal(t, ContextKey("request_id"), RequestIDKey)
	assert.Equal(t, ContextKey("user_id"), UserIDKey)
	assert.Equal(t, ContextKey("user_roles"), UserRolesKey)
	assert.Equal(t, ContextKey("start_time"), StartTimeKey)
	assert.Equal(t, ContextKey("middleware_error"), MiddlewareErrorKey)
}
