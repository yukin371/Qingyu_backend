package errors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestARCH005_UnifiedErrorStructure 测试统一错误结构
// 验证：ARCH-005要求 - 统一的错误处理结构
func TestARCH005_UnifiedErrorStructure(t *testing.T) {
	t.Run("UnifiedError应该包含完整的错误信息", func(t *testing.T) {
		err := UserServiceFactory.AuthError("UNAUTHORIZED", "未授权访问").
			WithContext("user123", "req-456", "trace-789").
			WithService("user-service", "login")

		assert.Equal(t, "UNAUTHORIZED", err.Code)
		assert.Equal(t, "未授权访问", err.Message)
		assert.Equal(t, "user123", err.UserID)
		assert.Equal(t, "req-456", err.RequestID)
		assert.Equal(t, "trace-789", err.TraceID)
		assert.Equal(t, "user-service", err.Service)
		assert.Equal(t, "login", err.Operation)
	})

	t.Run("应该支持错误链", func(t *testing.T) {
		underlyingErr := errors.New("database connection failed")
		err := UserServiceFactory.InternalError("DB_ERROR", "数据库错误", underlyingErr)

		assert.Equal(t, underlyingErr, err.Unwrap())
		assert.Contains(t, err.Error(), "数据库错误")
	})

	t.Run("应该支持添加元数据", func(t *testing.T) {
		err := UserServiceFactory.NotFoundError("user", "123").
			AddMetadata("attempt", "3").
			AddMetadata("ip", "192.168.1.1")

		assert.Equal(t, "3", err.Metadata["attempt"])
		assert.Equal(t, "192.168.1.1", err.Metadata["ip"])
	})
}

// TestARCH005_ErrorCodeCategories 测试错误码分类
// 验证：ARCH-005要求 - 系统级和业务级错误码
func TestARCH005_ErrorCodeCategories(t *testing.T) {
	t.Run("应该有系统级错误码", func(t *testing.T) {
		systemErrors := []ErrorCode{
			InvalidParams,
			Unauthorized,
			Forbidden,
			NotFound,
		}

		for _, code := range systemErrors {
			msg := GetDefaultMessage(code)
			assert.NotEmpty(t, msg, "错误码 %d 应该有默认消息", code)

			status := GetHTTPStatus(code)
			assert.NotZero(t, status, "错误码 %d 应该有HTTP状态码", code)
		}
	})

	t.Run("应该有业务级错误码", func(t *testing.T) {
		businessErrors := []ErrorCode{
			UserNotFound,
			InvalidCredentials,
			BookNotFound,
			ChapterNotFound,
			InsufficientQuota,
		}

		for _, code := range businessErrors {
			msg := GetDefaultMessage(code)
			assert.NotEmpty(t, msg)
			assert.GreaterOrEqual(t, int(code), 2000, "应该是业务级错误码")
		}
	})

	t.Run("错误码应该是唯一的", func(t *testing.T) {
		codes := []ErrorCode{
			InvalidParams,
			Unauthorized,
			Forbidden,
			NotFound,
			UserNotFound,
			BookNotFound,
			InsufficientQuota,
		}

		uniqueCodes := make(map[ErrorCode]bool)
		for _, code := range codes {
			if uniqueCodes[code] {
				t.Errorf("错误码 %d 重复", code)
			}
			uniqueCodes[code] = true
		}

		assert.Equal(t, len(codes), len(uniqueCodes))
	})
}

// TestARCH005_HTTPStatusMapping 测试HTTP状态码映射
func TestARCH005_HTTPStatusMapping(t *testing.T) {
	testCases := []struct {
		name           string
		errorCode      ErrorCode
		expectedStatus int
	}{
		{"未授权应该映射到401", Unauthorized, http.StatusUnauthorized},
		{"禁止访问应该映射到403", Forbidden, http.StatusForbidden},
		{"未找到应该映射到404", NotFound, http.StatusNotFound},
		{"用户不存在应该映射到404", UserNotFound, http.StatusNotFound},
		{"书籍不存在应该映射到404", BookNotFound, http.StatusNotFound},
		{"内部错误应该映射到500", InternalError, http.StatusInternalServerError},
		{"配额不足应该映射到400", InsufficientQuota, http.StatusBadRequest},
		{"频率限制应该映射到429", RateLimitExceeded, http.StatusTooManyRequests},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := GetHTTPStatus(tc.errorCode)
			assert.Equal(t, tc.expectedStatus, status)
		})
	}
}

// TestARCH005_ErrorCategories 测试错误分类
func TestARCH005_ErrorCategories(t *testing.T) {
	t.Run("UnifiedError应该有错误分类", func(t *testing.T) {
		categories := []ErrorCategory{
			CategoryValidation,
			CategoryBusiness,
			CategorySystem,
			CategoryExternal,
			CategoryNetwork,
			CategoryAuth,
			CategoryDatabase,
			CategoryCache,
		}

		validCategories := map[ErrorCategory]bool{
			CategoryValidation: true,
			CategoryBusiness:   true,
			CategorySystem:     true,
			CategoryExternal:   true,
			CategoryNetwork:    true,
			CategoryAuth:       true,
			CategoryDatabase:   true,
			CategoryCache:      true,
		}

		for _, cat := range categories {
			assert.True(t, validCategories[cat])
		}
	})

	t.Run("不同类型的错误应该有正确的分类", func(t *testing.T) {
		validationErr := UserServiceFactory.ValidationError("INVALID_PARAMS", "参数无效")
		assert.Equal(t, CategoryValidation, validationErr.Category)

		businessErr := UserServiceFactory.BusinessError("CONFLICT", "业务冲突")
		assert.Equal(t, CategoryBusiness, businessErr.Category)

		authErr := UserServiceFactory.AuthError("UNAUTHORIZED", "未授权")
		assert.Equal(t, CategoryAuth, authErr.Category)

		internalErr := UserServiceFactory.InternalError("INTERNAL", "内部错误", nil)
		assert.Equal(t, CategorySystem, internalErr.Category)
	})
}

// TestARCH005_ErrorBuilder 测试错误构建器模式
func TestARCH005_ErrorBuilder(t *testing.T) {
	t.Run("应该支持构建器模式", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("CUSTOM_ERROR").
			WithCategory(CategoryBusiness).
			WithLevel(LevelError).
			WithMessage("自定义错误").
			WithDetails("这是详细信息").
			WithHTTPStatus(400).
			WithRetryable(false).
			Build()

		assert.Equal(t, "CUSTOM_ERROR", err.Code)
		assert.Equal(t, CategoryBusiness, err.Category)
		assert.Equal(t, LevelError, err.Level)
		assert.Equal(t, "自定义错误", err.Message)
		assert.Equal(t, "这是详细信息", err.Details)
		assert.Equal(t, 400, err.HTTPStatus)
		assert.False(t, err.Retryable)
	})

	t.Run("构建器应该支持链式调用", func(t *testing.T) {
		err := NewErrorBuilder().
			WithCode("CHAIN_ERROR").
			WithMessage("链式调用测试").
			WithContext("user1", "req1", "trace1").
			WithService("test-service", "test-operation").
			WithMetadata("key1", "value1").
			Build()

		assert.Equal(t, "CHAIN_ERROR", err.Code)
		assert.Equal(t, "user1", err.UserID)
		assert.Equal(t, "test-service", err.Service)
		assert.Equal(t, "value1", err.Metadata["key1"])
	})
}

// TestARCH005_ConvenienceFunctions 测试便捷函数
func TestARCH005_ConvenienceFunctions(t *testing.T) {
	t.Run("New函数应该创建基础错误", func(t *testing.T) {
		err := New(InvalidParams, "参数无效")
		assert.Equal(t, "1001", err.Code)
		assert.Equal(t, "参数无效", err.Message)
		assert.Equal(t, 400, err.HTTPStatus)
	})

	t.Run("NewUnauthorized应该创建401错误", func(t *testing.T) {
		err := NewUnauthorized("请先登录")
		assert.Equal(t, 401, err.HTTPStatus)
		assert.Equal(t, "请先登录", err.Message)
	})

	t.Run("NewForbidden应该创建403错误", func(t *testing.T) {
		err := NewForbidden("权限不足")
		assert.Equal(t, 403, err.HTTPStatus)
		assert.Equal(t, "权限不足", err.Message)
	})

	t.Run("NewNotFound应该创建404错误", func(t *testing.T) {
		err := NewNotFound("用户")
		assert.Equal(t, 404, err.HTTPStatus)
		assert.Contains(t, err.Message, "not found")
	})

	t.Run("NewInternal应该创建500错误", func(t *testing.T) {
		underlying := errors.New("something failed")
		err := NewInternal(underlying, "服务器内部错误")
		assert.Equal(t, 500, err.HTTPStatus)
		assert.Equal(t, underlying, err.Cause)
	})

	t.Run("NewRateLimit应该创建429错误", func(t *testing.T) {
		err := NewRateLimit()
		assert.Equal(t, 429, err.HTTPStatus)
		assert.NotEmpty(t, err.Message)
	})
}

// TestARCH005_ErrorRetryable 测试错误可重试判断
func TestARCH005_ErrorRetryable(t *testing.T) {
	t.Run("某些错误应该可以重试", func(t *testing.T) {
		networkErr := UserServiceFactory.NetworkError("连接超时")
		assert.True(t, networkErr.IsRetryable())

		timeoutErr := UserServiceFactory.TimeoutError("database query")
		assert.True(t, timeoutErr.IsRetryable())

		cacheErr := UserServiceFactory.CacheError("get", nil)
		assert.True(t, cacheErr.IsRetryable())
	})

	t.Run("某些错误不应该重试", func(t *testing.T) {
		validationErr := UserServiceFactory.ValidationError("INVALID", "无效")
		assert.False(t, validationErr.IsRetryable())

		businessErr := UserServiceFactory.BusinessError("CONFLICT", "冲突")
		assert.False(t, businessErr.IsRetryable())

		authErr := UserServiceFactory.AuthError("UNAUTHORIZED", "未授权")
		assert.False(t, authErr.IsRetryable())
	})
}

// TestARCH005_ErrorJSONSerialization 测试错误JSON序列化
func TestARCH005_ErrorJSONSerialization(t *testing.T) {
	t.Run("UnifiedError应该可以序列化为JSON", func(t *testing.T) {
		appErr := UserServiceFactory.AuthError("UNAUTHORIZED", "未授权").
			WithContext("user123", "req456", "trace789")

		jsonData, err := appErr.ToJSON()
		assert.NoError(t, err)
		assert.NotEmpty(t, jsonData)

		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "UNAUTHORIZED")
		assert.Contains(t, jsonStr, "user123")
	})
}

// TestARCH005_ServiceFactories 测试服务工厂
func TestARCH005_ServiceFactories(t *testing.T) {
	t.Run("应该有不同服务的错误工厂", func(t *testing.T) {
		factories := []*ErrorFactory{
			AIServiceFactory,
			UserServiceFactory,
			DocumentServiceFactory,
			ProjectFactory,
			BookstoreServiceFactory,
			ReaderServiceFactory,
			WriterServiceFactory,
		}

		for _, factory := range factories {
			assert.NotNil(t, factory)
			err := factory.InternalError("TEST", "test error", nil)
			assert.Equal(t, factory.service, err.Service)
		}
	})
}

// TestARCH005_AIErrorIntegration 测试AI错误集成
func TestARCH005_AIErrorIntegration(t *testing.T) {
	t.Run("AIError应该能转换为ErrorCode", func(t *testing.T) {
		quotaErr := NewAIError(ErrQuotaExhausted, "配额不足", nil)
		code := quotaErr.ToErrorCode()
		assert.Equal(t, InsufficientQuota, code)

		invalidErr := NewAIError(ErrInvalidRequest, "请求无效", nil)
		code = invalidErr.ToErrorCode()
		assert.Equal(t, InvalidParams, code)
	})

	t.Run("AIError应该能转换为gRPC status", func(t *testing.T) {
		err := NewAIError(ErrQuotaExhausted, "配额不足", nil)
		status := err.ToStatus()
		assert.NotNil(t, status)
		assert.Equal(t, "配额不足", status.Message())
	})
}

// TestARCH005_BackwardCompatibility 测试向后兼容性
func TestARCH005_BackwardCompatibility(t *testing.T) {
	t.Run("现有的ErrorCode常量应该可用", func(t *testing.T) {
		assert.Equal(t, ErrorCode(1001), InvalidParams)
		assert.Equal(t, ErrorCode(1002), Unauthorized)
		assert.Equal(t, ErrorCode(1003), Forbidden)
		assert.Equal(t, ErrorCode(1004), NotFound)
		assert.Equal(t, ErrorCode(2001), UserNotFound)
		assert.Equal(t, ErrorCode(3001), BookNotFound)
		assert.Equal(t, ErrorCode(3002), ChapterNotFound)
		assert.Equal(t, ErrorCode(3010), InsufficientQuota)
	})

	t.Run("现有的便捷函数应该可用", func(t *testing.T) {
		err1 := New(InvalidParams, "test")
		assert.NotNil(t, err1)

		err2 := NewUnauthorized("test")
		assert.NotNil(t, err2)

		err3 := NewForbidden("test")
		assert.NotNil(t, err3)

		err4 := NewNotFound("test")
		assert.NotNil(t, err4)
	})

	t.Run("GetDefaultMessage应该返回正确的消息", func(t *testing.T) {
		msg := GetDefaultMessage(UserNotFound)
		assert.Equal(t, "用户不存在", msg)

		msg = GetDefaultMessage(InsufficientQuota)
		assert.Equal(t, "配额不足", msg)
	})

	t.Run("GetHTTPStatus应该返回正确的状态码", func(t *testing.T) {
		status := GetHTTPStatus(Unauthorized)
		assert.Equal(t, 401, status)

		status = GetHTTPStatus(UserNotFound)
		assert.Equal(t, 404, status)

		status = GetHTTPStatus(InsufficientQuota)
		assert.Equal(t, 400, status)
	})
}

// TestARCH005_LogFieldsRequirement 测试日志字段需求
func TestARCH005_LogFieldsRequirement(t *testing.T) {
	t.Run("错误应该包含足够的信息用于日志", func(t *testing.T) {
		err := UserServiceFactory.NotFoundError("user", "123").
			WithContext("user123", "req456", "trace789").
			AddMetadata("attempt", "3")

		assert.NotEmpty(t, err.ID)
		assert.NotEmpty(t, err.Code)
		assert.NotEmpty(t, err.Message)
		assert.NotEmpty(t, err.Service)
		assert.NotEmpty(t, err.Category)
		assert.NotEmpty(t, err.Level)
		assert.NotZero(t, err.HTTPStatus)
		assert.NotZero(t, err.Timestamp)
	})

	t.Run("应该能够获取所有日志字段", func(t *testing.T) {
		err := UserServiceFactory.AuthError("TEST", "test message")

		assert.NotEmpty(t, err.ID)
		assert.NotEmpty(t, err.Code)
		assert.NotEmpty(t, err.Message)
		assert.NotEmpty(t, err.Category)
		assert.NotEmpty(t, err.Level)
		assert.NotEmpty(t, err.Service)
		assert.NotZero(t, err.HTTPStatus)
		assert.False(t, err.Timestamp.IsZero())
	})
}
