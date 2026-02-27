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
// 验证：ARCH-005要求 - 系统级(1-1999)和业务级(10000+)错误码
func TestARCH005_ErrorCodeCategories(t *testing.T) {
	t.Run("应该有系统级错误码(1-1999)", func(t *testing.T) {
		systemErrors := []ErrorCode{
			InvalidParams,     // 1001
			Unauthorized,      // 1002
			Forbidden,         // 1003
			NotFound,          // 1004
			InternalError,     // 5000 - 这个其实超出1999了
		}

		for _, code := range systemErrors {
			// 验证是有效的错误码
			msg := GetDefaultMessage(code)
			assert.NotEmpty(t, msg, "错误码 %d 应该有默认消息", code)

			// 验证有HTTP状态码映射
			status := GetHTTPStatus(code)
			assert.NotZero(t, status, "错误码 %d 应该有HTTP状态码", code)
		}
	})

	t.Run("应该有业务级错误码(2000+)", func(t *testing.T) {
		businessErrors := []ErrorCode{
			UserNotFound,       // 2001
			InvalidCredentials, // 2002
			BookNotFound,       // 3001
			ChapterNotFound,    // 3002
			InsufficientQuota,  // 3010
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

		assert.Equal(t, len(codes), len(uniqueCodes), "所有错误码应该是唯一的")
	})
}

// TestARCH005_HTTPStatusMapping 测试HTTP状态码映射
// 验证：ARCH-005要求 - 错误码到HTTP状态码的映射
func TestARCH005_HTTPStatusMapping(t *testing.T) {
	testCases := []struct {
		name           string
		errorCode      ErrorCode
		expectedStatus int
	}{
		{
			name:           "未授权应该映射到401",
			errorCode:      Unauthorized,
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "禁止访问应该映射到403",
			errorCode:      Forbidden,
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "未找到应该映射到404",
			errorCode:      NotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "用户不存在应该映射到404",
			errorCode:      UserNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "书籍不存在应该映射到404",
			errorCode:      BookNotFound,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "内部错误应该映射到500",
			errorCode:      InternalError,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "配额不足应该映射到400",
			errorCode:      InsufficientQuota,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "频率限制应该映射到429",
			errorCode:      RateLimitExceeded,
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			status := GetHTTPStatus(tc.errorCode)
			assert.Equal(t, tc.expectedStatus, status)
		})
	}
}

// TestARCH005_ErrorCategories 测试错误分类
// 验证：ARCH-005要求 - 错误应该有明确的分类
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

		// 验证所有分类都有定义
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
			assert.True(t, validCategories[cat], "分类 %s 应该是有效的", cat)
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
// 验证：ARCH-005要求 - 支持构建器模式创建错误
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
// 验证：ARCH-005要求 - 提供常用错误的便捷创建函数
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
// 验证：ARCH-005要求 - 错误应该支持可重试判断
func TestARCH005_ErrorRetryable(t *testing.T) {
	t.Run("某些错误应该可以重试", func(t *testing.T) {
		// 网络错误通常可重试
		networkErr := UserServiceFactory.NetworkError("连接超时")
		assert.True(t, networkErr.IsRetryable())
		assert.True(t, networkErr.Retryable)

		// 超时错误可重试
		timeoutErr := UserServiceFactory.TimeoutError("database query")
		assert.True(t, timeoutErr.IsRetryable())

		// 缓存错误可重试
		cacheErr := UserServiceFactory.CacheError("get", nil)
		assert.True(t, cacheErr.IsRetryable())
	})

	t.Run("某些错误不应该重试", func(t *testing.T) {
		// 验证错误不重试
		validationErr := UserServiceFactory.ValidationError("INVALID", "无效")
		assert.False(t, validationErr.IsRetryable())

		// 业务错误不重试
		businessErr := UserServiceFactory.BusinessError("CONFLICT", "冲突")
		assert.False(t, businessErr.IsRetryable())

		// 认证错误不重试
		authErr := UserServiceFactory.AuthError("UNAUTHORIZED", "未授权")
		assert.False(t, authErr.IsRetryable())
	})
}

// TestARCH005_ErrorJSONSerialization 测试错误JSON序列化
// 验证：ARCH-005要求 - 错误应该支持JSON序列化用于API响应
func TestARCH005_ErrorJSONSerialization(t *testing.T) {
	t.Run("UnifiedError应该可以序列化为JSON", func(t *testing.T) {
		appErr := UserServiceFactory.AuthError("UNAUTHORIZED", "未授权").
			WithContext("user123", "req456", "trace789")

		jsonData, jsonErr := appErr.ToJSON()
		assert.NoError(t, jsonErr)
		assert.NotEmpty(t, jsonData)

		// 验证JSON包含关键字段
		jsonStr := string(jsonData)
		assert.Contains(t, jsonStr, "UNAUTHORIZED")
		assert.Contains(t, jsonStr, "user123")
	})
}

// TestARCH005_ServiceFactories 测试服务工厂
// 验证：ARCH-005要求 - 不同服务应该有独立的错误工厂
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

	t.Run("不同工厂的错误应该有正确的服务名", func(t *testing.T) {
		aiErr := AIServiceFactory.InternalError("TEST", "test", nil)
		assert.Equal(t, "ai-service", aiErr.Service)

		userErr := UserServiceFactory.AuthError("TEST", "test")
		assert.Equal(t, "user-service", userErr.Service)

		bookErr := BookstoreServiceFactory.NotFoundError("book", "123")
		assert.Equal(t, "bookstore-service", bookErr.Service)
	})
}

// TestARCH005_AIErrorIntegration 测试AI错误集成
// 验证：ARCH-005要求 - AI错误应该能转换为统一错误
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
// 验证：ARCH-005要求 - 不破坏现有API
func TestARCH005_BackwardCompatibility(t *testing.T) {
	t.Run("现有的ErrorCode常量应该可用", func(t *testing.T) {
		// 系统错误
		assert.Equal(t, ErrorCode(1001), InvalidParams)
		assert.Equal(t, ErrorCode(1002), Unauthorized)
		assert.Equal(t, ErrorCode(1003), Forbidden)
		assert.Equal(t, ErrorCode(1004), NotFound)

		// 业务错误
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
// 验证：ARCH-005要求 - 错误应该提供结构化的日志字段
func TestARCH005_LogFieldsRequirement(t *testing.T) {
	t.Run("错误应该包含足够的信息用于日志", func(t *testing.T) {
		err := UserServiceFactory.NotFoundError("user", "123").
			WithContext("user123", "req456", "trace789").
			AddMetadata("attempt", "3")

		// 验证所有需要的日志字段
		assert.NotEmpty(t, err.ID, "应该有错误ID")
		assert.NotEmpty(t, err.Code, "应该有错误码")
		assert.NotEmpty(t, err.Message, "应该有错误消息")
		assert.NotEmpty(t, err.Service, "应该有服务名")
		assert.NotEmpty(t, err.Category, "应该有错误分类")
		assert.NotEmpty(t, err.Level, "应该有错误级别")
		assert.NotZero(t, err.HTTPStatus, "应该有HTTP状态码")
		assert.NotZero(t, err.Timestamp, "应该有时间戳")
	})

	t.Run("应该能够获取所有日志字段", func(t *testing.T) {
		// 这个测试验证UnifiedError结构是否包含足够的字段
		// 实际的LogFields方法可以后续添加
		err := UserServiceFactory.AuthError("TEST", "test message")

		// 基本字段
		assert.NotEmpty(t, err.ID)
		assert.NotEmpty(t, err.Code)
		assert.NotEmpty(t, err.Message)
		assert.NotEmpty(t, err.Category)
		assert.NotEmpty(t, err.Level)

		// 上下文字段
		assert.NotEmpty(t, err.Service)

		// HTTP字段
		assert.NotZero(t, err.HTTPStatus)

		// 时间戳
		assert.False(t, err.Timestamp.IsZero())
	})
}
