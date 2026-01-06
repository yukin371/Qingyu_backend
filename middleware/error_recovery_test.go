package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	apperrors "Qingyu_backend/pkg/errors"
)

// TestErrorRecoveryMiddleware_PanicRecovery 测试panic恢复
func TestErrorRecoveryMiddleware_PanicRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorRecoveryMiddleware())

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	router.GET("/panic-error", func(c *gin.Context) {
		panic(errors.New("error panic"))
	})

	router.GET("/success", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "success"})
	})

	// 测试panic恢复
	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// 测试error panic
	req, _ = http.NewRequest("GET", "/panic-error", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// 测试正常请求
	req, _ = http.NewRequest("GET", "/success", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

// TestErrorRecoveryMiddleware_WithStacktrace 测试启用堆栈跟踪
func TestErrorRecoveryMiddleware_WithStacktrace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	config := ErrorRecoveryConfig{
		EnableStacktrace: true,
		EnableLogging:    false,
	}

	router := gin.New()
	router.Use(ErrorRecoveryMiddleware(config))

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic with stack")
	})

	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// 检查响应包含堆栈信息
	// 在生产环境中不应该返回堆栈
}

// TestErrorHandlerMiddleware_UnifiedError 测试UnifiedError处理
func TestErrorHandlerMiddleware_UnifiedError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(ErrorHandlerMiddleware())

	router.GET("/unified-error", func(c *gin.Context) {
		factory := apperrors.NewErrorFactory("test-service")
		err := factory.NotFoundError("user", "123")
		c.Error(err)
	})

	router.GET("/standard-error", func(c *gin.Context) {
		c.Error(errors.New("standard error"))
	})

	router.GET("/validation-error", func(c *gin.Context) {
		c.Error(errors.New("invalid parameter: email"))
	})

	router.GET("/unauthorized-error", func(c *gin.Context) {
		c.Error(errors.New("unauthorized access"))
	})

	// 测试UnifiedError
	req, _ := http.NewRequest("GET", "/unified-error", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}

	// 测试标准错误
	req, _ = http.NewRequest("GET", "/standard-error", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 500 {
		t.Errorf("Expected status 500, got %d", w.Code)
	}

	// 测试参数验证错误
	req, _ = http.NewRequest("GET", "/validation-error", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 400 {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	// 测试未授权错误
	req, _ = http.NewRequest("GET", "/unauthorized-error", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != 401 {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

// TestWrapError 测试错误包装
func TestWrapError(t *testing.T) {
	// 测试包装标准错误
	stdErr := errors.New("standard error")
	wrapped := WrapError("test-service", "test-operation", stdErr)

	if wrapped == nil {
		t.Fatal("Wrapped error should not be nil")
	}

	if wrapped.Service != "test-service" {
		t.Errorf("Expected service 'test-service', got '%s'", wrapped.Service)
	}

	if wrapped.Operation != "test-operation" {
		t.Errorf("Expected operation 'test-operation', got '%s'", wrapped.Operation)
	}

	// 测试包装nil错误
	wrappedNil := WrapError("test-service", "test-operation", nil)
	if wrappedNil != nil {
		t.Error("Wrapping nil error should return nil")
	}

	// 测试包装UnifiedError
	factory := apperrors.NewErrorFactory("original-service")
	unifiedErr := factory.NotFoundError("resource", "123")
	wrappedUnified := WrapError("wrapper-service", "wrap-operation", unifiedErr)

	if wrappedUnified.Service != "wrapper-service" {
		t.Errorf("Expected service 'wrapper-service', got '%s'", wrappedUnified.Service)
	}
}

// TestNotFoundError 测试404错误创建
func TestNotFoundError(t *testing.T) {
	err := NotFoundError("user", "123")

	if err == nil {
		t.Fatal("NotFoundError should not return nil")
	}

	if err.Category != apperrors.CategoryBusiness {
		t.Errorf("Expected category '%s', got '%s'", apperrors.CategoryBusiness, err.Category)
	}

	if err.HTTPStatus != 404 {
		t.Errorf("Expected HTTP status 404, got %d", err.HTTPStatus)
	}
}

// TestValidationError 测试验证错误创建
func TestValidationError(t *testing.T) {
	err := ValidationError("email", "invalid email format")

	if err == nil {
		t.Fatal("ValidationError should not return nil")
	}

	if err.Category != apperrors.CategoryValidation {
		t.Errorf("Expected category '%s', got '%s'", apperrors.CategoryValidation, err.Category)
	}

	if err.HTTPStatus != 400 {
		t.Errorf("Expected HTTP status 400, got %d", err.HTTPStatus)
	}

	if err.Details == "" {
		t.Error("ValidationError should have details")
	}
}

// TestBusinessError 测试业务错误创建
func TestBusinessError(t *testing.T) {
	err := BusinessError("INSUFFICIENT_BALANCE", "余额不足")

	if err == nil {
		t.Fatal("BusinessError should not return nil")
	}

	if err.Category != apperrors.CategoryBusiness {
		t.Errorf("Expected category '%s', got '%s'", apperrors.CategoryBusiness, err.Category)
	}

	if err.Code != "INSUFFICIENT_BALANCE" {
		t.Errorf("Expected code 'INSUFFICIENT_BALANCE', got '%s'", err.Code)
	}

	if err.Message != "余额不足" {
		t.Errorf("Expected message '余额不足', got '%s'", err.Message)
	}
}
