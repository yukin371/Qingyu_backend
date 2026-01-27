package core

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ContextKey 是context键的类型，避免键冲突
type ContextKey string

// 常用的context键
const (
	// RequestIDKey 请求ID
	RequestIDKey ContextKey = "request_id"
	// UserIDKey 用户ID
	UserIDKey ContextKey = "user_id"
	// UserRolesKey 用户角色列表
	UserRolesKey ContextKey = "user_roles"
	// StartTimeKey 请求开始时间
	StartTimeKey ContextKey = "start_time"
	// MiddlewareErrorKey 中间件错误
	MiddlewareErrorKey ContextKey = "middleware_error"
)

// RequestContext 请求上下文
//
// 封装Gin的Context，提供更便捷的访问方法。
type RequestContext struct {
	*gin.Context
	logger *zap.Logger
}

// NewRequestContext 创建请求上下文
func NewRequestContext(c *gin.Context, logger *zap.Logger) *RequestContext {
	return &RequestContext{
		Context: c,
		logger:  logger,
	}
}

// GetRequestID 获取请求ID
func (rc *RequestContext) GetRequestID() string {
	if requestID, exists := rc.Get(string(RequestIDKey)); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID 获取用户ID
func (rc *RequestContext) GetUserID() string {
	if userID, exists := rc.Get(string(UserIDKey)); exists {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserRoles 获取用户角色列表
func (rc *RequestContext) GetUserRoles() []string {
	if roles, exists := rc.Get(string(UserRolesKey)); exists {
		if r, ok := roles.([]string); ok {
			return r
		}
	}
	return nil
}

// GetStartTime 获取请求开始时间
func (rc *RequestContext) GetStartTime() time.Time {
	if startTime, exists := rc.Get(string(StartTimeKey)); exists {
		if t, ok := startTime.(time.Time); ok {
			return t
		}
	}
	return time.Now()
}

// SetError 设置错误
func (rc *RequestContext) SetError(err error) {
	rc.Set(string(MiddlewareErrorKey), err)
}

// GetError 获取错误
func (rc *RequestContext) GetError() error {
	if err, exists := rc.Get(string(MiddlewareErrorKey)); exists {
		if e, ok := err.(error); ok {
			return e
		}
	}
	return nil
}

// WithContext 创建带超时的context
func (rc *RequestContext) WithContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(rc.Request.Context(), timeout)
	return ctx, cancel
}

// LogInfo 记录信息日志
func (rc *RequestContext) LogInfo(msg string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("request_id", rc.GetRequestID()),
		zap.String("path", rc.Request.URL.Path),
		zap.String("method", rc.Request.Method),
	}, fields...)
	rc.logger.Info(msg, allFields...)
}

// LogError 记录错误日志
func (rc *RequestContext) LogError(msg string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("request_id", rc.GetRequestID()),
		zap.String("path", rc.Request.URL.Path),
		zap.String("method", rc.Request.Method),
	}, fields...)
	rc.logger.Error(msg, allFields...)
}

// LogWarn 记录警告日志
func (rc *RequestContext) LogWarn(msg string, fields ...zap.Field) {
	allFields := append([]zap.Field{
		zap.String("request_id", rc.GetRequestID()),
		zap.String("path", rc.Request.URL.Path),
		zap.String("method", rc.Request.Method),
	}, fields...)
	rc.logger.Warn(msg, allFields...)
}
