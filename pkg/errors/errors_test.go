package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestErrorBuilder 测试错误构建器
func TestErrorBuilder(t *testing.T) {
	tests := []struct {
		name      string
		buildErr  func() *UnifiedError
		wantCode  string
		wantMsg   string
		wantHTTP  int
	}{
		{
			name: "无效参数错误",
			buildErr: func() *UnifiedError {
				return NewErrorBuilder().
					WithCode("INVALID_PARAMS").
					WithCategory(CategoryValidation).
					WithLevel(LevelWarning).
					WithMessage("测试参数无效").
					WithHTTPStatus(http.StatusBadRequest).
					Build()
			},
			wantCode: "INVALID_PARAMS",
			wantMsg:  "测试参数无效",
			wantHTTP: http.StatusBadRequest,
		},
		{
			name: "资源不存在",
			buildErr: func() *UnifiedError {
				return NewErrorBuilder().
					WithCode("NOT_FOUND").
					WithCategory(CategoryBusiness).
					WithLevel(LevelWarning).
					WithMessage("用户不存在").
					WithHTTPStatus(http.StatusNotFound).
					Build()
			},
			wantCode: "NOT_FOUND",
			wantMsg:  "用户不存在",
			wantHTTP: http.StatusNotFound,
		},
		{
			name: "内部错误",
			buildErr: func() *UnifiedError {
				return NewErrorBuilder().
					WithCode("INTERNAL_ERROR").
					WithCategory(CategorySystem).
					WithLevel(LevelError).
					WithMessage("内部错误").
					WithHTTPStatus(http.StatusInternalServerError).
					Build()
			},
			wantCode: "INTERNAL_ERROR",
			wantMsg:  "内部错误",
			wantHTTP: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.buildErr()

			if err.Code != tt.wantCode {
				t.Errorf("Code = %v, want %v", err.Code, tt.wantCode)
			}

			if err.Message != tt.wantMsg {
				t.Errorf("Message = %v, want %v", err.Message, tt.wantMsg)
			}

			if err.HTTPStatus != tt.wantHTTP {
				t.Errorf("HTTPStatus = %v, want %v", err.HTTPStatus, tt.wantHTTP)
			}
		})
	}
}

// TestErrorFactory 测试错误工厂
func TestErrorFactory(t *testing.T) {
	factory := NewErrorFactory("test-service")

	t.Run("ValidationError", func(t *testing.T) {
		err := factory.ValidationError("INVALID_PARAMS", "参数无效", "详情信息")

		assert.Equal(t, "INVALID_PARAMS", err.Code)
		assert.Equal(t, "参数无效", err.Message)
		assert.Equal(t, "详情信息", err.Details)
		assert.Equal(t, CategoryValidation, err.Category)
		assert.Equal(t, LevelWarning, err.Level)
		assert.Equal(t, 400, err.HTTPStatus)
		assert.False(t, err.Retryable)
	})

	t.Run("BusinessError", func(t *testing.T) {
		err := factory.BusinessError("BUSINESS_ERROR", "业务错误")

		assert.Equal(t, "BUSINESS_ERROR", err.Code)
		assert.Equal(t, "业务错误", err.Message)
		assert.Equal(t, CategoryBusiness, err.Category)
		assert.Equal(t, 409, err.HTTPStatus)
	})

	t.Run("NotFoundError", func(t *testing.T) {
		err := factory.NotFoundError("User", "123")

		assert.Equal(t, "NOT_FOUND", err.Code)
		assert.Contains(t, err.Message, "not found")
		assert.Equal(t, 404, err.HTTPStatus)
	})

	t.Run("AuthError", func(t *testing.T) {
		err := factory.AuthError("UNAUTHORIZED", "未授权")

		assert.Equal(t, "UNAUTHORIZED", err.Code)
		assert.Equal(t, "未授权", err.Message)
		assert.Equal(t, CategoryAuth, err.Category)
		assert.Equal(t, 401, err.HTTPStatus)
	})

	t.Run("InternalError", func(t *testing.T) {
		originalErr := errors.New("original error")
		err := factory.InternalError("DB_ERROR", "数据库错误", originalErr)

		assert.Equal(t, "DB_ERROR", err.Code)
		assert.Equal(t, "数据库错误", err.Message)
		assert.Equal(t, originalErr, err.Cause)
		assert.Equal(t, CategorySystem, err.Category)
		assert.Equal(t, 500, err.HTTPStatus)
	})

	t.Run("ExternalError", func(t *testing.T) {
		err := factory.ExternalError("EXTERNAL_ERROR", "外部服务错误", true)

		assert.Equal(t, "EXTERNAL_ERROR", err.Code)
		assert.Equal(t, "外部服务错误", err.Message)
		assert.Equal(t, CategoryExternal, err.Category)
		assert.True(t, err.Retryable)
	})

	t.Run("NetworkError", func(t *testing.T) {
		err := factory.NetworkError("网络连接失败")

		assert.Equal(t, "NETWORK_ERROR", err.Code)
		assert.Equal(t, "网络连接失败", err.Message)
		assert.Equal(t, CategoryNetwork, err.Category)
		assert.True(t, err.Retryable)
		assert.Equal(t, 503, err.HTTPStatus)
	})

	t.Run("TimeoutError", func(t *testing.T) {
		err := factory.TimeoutError("查询数据")

		assert.Equal(t, "TIMEOUT", err.Code)
		assert.Contains(t, err.Message, "timed out")
		assert.Equal(t, "查询数据", err.Operation)
		assert.True(t, err.Retryable)
		assert.Equal(t, 408, err.HTTPStatus)
	})

	t.Run("RateLimitError", func(t *testing.T) {
		err := factory.RateLimitError(100)

		assert.Equal(t, "RATE_LIMIT_EXCEEDED", err.Code)
		assert.Equal(t, "Rate limit exceeded", err.Message)
		assert.Contains(t, err.Details, "100")
		assert.True(t, err.Retryable)
		assert.Equal(t, 429, err.HTTPStatus)
	})

	t.Run("DatabaseError", func(t *testing.T) {
		originalErr := errors.New("connection failed")
		err := factory.DatabaseError("insert", originalErr)

		assert.Equal(t, "DATABASE_ERROR", err.Code)
		assert.Contains(t, err.Message, "insert")
		assert.Equal(t, originalErr, err.Cause)
		assert.Equal(t, CategoryDatabase, err.Category)
		assert.False(t, err.Retryable)
	})

	t.Run("CacheError", func(t *testing.T) {
		originalErr := errors.New("redis error")
		err := factory.CacheError("get", originalErr)

		assert.Equal(t, "CACHE_ERROR", err.Code)
		assert.Contains(t, err.Message, "get")
		assert.Equal(t, originalErr, err.Cause)
		assert.Equal(t, CategoryCache, err.Category)
		assert.True(t, err.Retryable)
	})
}

// TestUnifiedError_Methods 测试UnifiedError的方法
func TestUnifiedError_Methods(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithMessage("测试错误").
			WithDetails("详细信息").
			Build()

		errStr := err.Error()
		assert.Contains(t, errStr, "TEST_ERROR")
		assert.Contains(t, errStr, "测试错误")
		assert.Contains(t, errStr, "详细信息")
	})

	t.Run("Error_WithoutDetails", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithMessage("测试错误").
			Build()

		errStr := err.Error()
		assert.Contains(t, errStr, "TEST_ERROR")
		assert.Contains(t, errStr, "测试错误")
	})

	t.Run("Unwrap", func(t *testing.T) {
		originalErr := errors.New("original")
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithCause(originalErr).
			Build()

		assert.Equal(t, originalErr, err.Unwrap())
	})

	t.Run("GetHTTPStatus", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithCategory(CategoryValidation).
			Build()

		assert.Equal(t, http.StatusBadRequest, err.GetHTTPStatus())
	})

	t.Run("GetHTTPStatus_WithCustom", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithHTTPStatus(418).
			Build()

		assert.Equal(t, 418, err.GetHTTPStatus())
	})

	t.Run("IsRetryable", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithRetryable(true).
			Build()

		assert.True(t, err.IsRetryable())
	})

	t.Run("AddMetadata", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			Build()

		err.AddMetadata("key1", "value1")
		err.AddMetadata("key2", 123)

		assert.Equal(t, "value1", err.Metadata["key1"])
		assert.Equal(t, 123, err.Metadata["key2"])
	})

	t.Run("WithContext", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			Build()

		err = err.WithContext("user123", "req456", "trace789")

		assert.Equal(t, "user123", err.UserID)
		assert.Equal(t, "req456", err.RequestID)
		assert.Equal(t, "trace789", err.TraceID)
	})

	t.Run("WithService", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			Build()

		err = err.WithService("test-service", "test-operation")

		assert.Equal(t, "test-service", err.Service)
		assert.Equal(t, "test-operation", err.Operation)
	})

	t.Run("WithMetadata", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("TEST_ERROR").
			Build()

		metadata := map[string]interface{}{
			"key1": "value1",
			"key2": 123,
		}
		err = err.WithMetadata(metadata)

		assert.Equal(t, "value1", err.Metadata["key1"])
		assert.Equal(t, 123, err.Metadata["key2"])
	})
}

// TestPredefinedFactories 测试预定义的工厂
func TestPredefinedFactories(t *testing.T) {
	t.Run("AIServiceFactory", func(t *testing.T) {
		err := AIServiceFactory.ValidationError("INVALID_AI_REQUEST", "AI请求无效")

		assert.Equal(t, "ai-service", err.Service)
	})

	t.Run("UserServiceFactory", func(t *testing.T) {
		err := UserServiceFactory.NotFoundError("User", "123")

		assert.Equal(t, "user-service", err.Service)
	})

	t.Run("DocumentServiceFactory", func(t *testing.T) {
		err := DocumentServiceFactory.AuthError("TOKEN_EXPIRED", "Token过期")

		assert.Equal(t, "document-service", err.Service)
	})

	t.Run("ProjectFactory", func(t *testing.T) {
		err := ProjectFactory.BusinessError("PROJECT_LOCKED", "项目已锁定")

		assert.Equal(t, "project-service", err.Service)
	})

	t.Run("BookstoreServiceFactory", func(t *testing.T) {
		err := BookstoreServiceFactory.DatabaseError("purchase", errors.New("db error"))

		assert.Equal(t, "bookstore-service", err.Service)
	})

	t.Run("ReaderServiceFactory", func(t *testing.T) {
		err := ReaderServiceFactory.NetworkError("网络错误")

		assert.Equal(t, "reader-service", err.Service)
	})

	t.Run("WriterServiceFactory", func(t *testing.T) {
		err := WriterServiceFactory.RateLimitError(100)

		assert.Equal(t, "writer-service", err.Service)
	})
}

// TestErrorBuilder_Chaining 测试错误构建器链式调用
func TestErrorBuilder_Chaining(t *testing.T) {
	originalErr := errors.New("original error")

	err := NewErrorBuilder().
		WithID("test-error-123").
		WithCode("TEST_ERROR").
		WithCategory(CategoryValidation).
		WithLevel(LevelWarning).
		WithMessage("测试错误").
		WithDetails("详细信息").
		WithCause(originalErr).
		WithHTTPStatus(http.StatusBadRequest).
		WithRetryable(false).
		WithService("test-service", "test-operation").
		WithContext("user123", "req456", "trace789").
		WithMetadata("key1", "value1").
		Build()

	assert.Equal(t, "test-error-123", err.ID)
	assert.Equal(t, "TEST_ERROR", err.Code)
	assert.Equal(t, CategoryValidation, err.Category)
	assert.Equal(t, LevelWarning, err.Level)
	assert.Equal(t, "测试错误", err.Message)
	assert.Equal(t, "详细信息", err.Details)
	assert.Equal(t, originalErr, err.Cause)
	assert.Equal(t, http.StatusBadRequest, err.HTTPStatus)
	assert.False(t, err.Retryable)
	assert.Equal(t, "test-service", err.Service)
	assert.Equal(t, "test-operation", err.Operation)
	assert.Equal(t, "user123", err.UserID)
	assert.Equal(t, "req456", err.RequestID)
	assert.Equal(t, "trace789", err.TraceID)
	assert.Equal(t, "value1", err.Metadata["key1"])
}

// BenchmarkNew 性能测试 - 创建错误
func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewErrorBuilder().
			WithCode("TEST_ERROR").
			WithMessage("测试错误").
			Build()
	}
}

// BenchmarkFactory 性能测试 - 使用工厂创建错误
func BenchmarkFactory(b *testing.B) {
	factory := NewErrorFactory("test-service")
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = factory.ValidationError("TEST_ERROR", "测试错误")
	}
}
