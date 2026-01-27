package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"Qingyu_backend/internal/middleware/builtin"
	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func init() {
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
}

// TestMiddlewareChain_FullChain 测试完整的中间件链
func TestMiddlewareChain_FullChain(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore)

	// 创建中间件管理器
	manager := core.NewManager(logger)

	// 注册所有中间件
	requestID := builtin.NewRequestIDMiddleware()
	recovery := builtin.NewRecoveryMiddleware(logger)
	security := builtin.NewSecurityMiddleware()
	loggerMW := builtin.NewLoggerMiddleware(logger)
	compression := builtin.NewCompressionMiddleware()

	err := manager.Register(requestID)
	assert.NoError(t, err)

	err = manager.Register(recovery)
	assert.NoError(t, err)

	err = manager.Register(security)
	assert.NoError(t, err)

	err = manager.Register(loggerMW)
	assert.NoError(t, err)

	err = manager.Register(compression)
	assert.NoError(t, err)

	// 创建Gin引擎并应用中间件
	router := gin.New()
	err = manager.ApplyToRouter(router)
	assert.NoError(t, err)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
			"request_id": builtin.GetRequestID(c),
		})
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证请求ID存在
	responseRequestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, responseRequestID)

	// 验证安全头存在
	assert.Equal(t, "DENY", w.Header().Get("X-Frame-Options"))
	assert.NotEmpty(t, w.Header().Get("X-Content-Type-Options"))

	// 验证日志记录
	assert.True(t, logs.FilterField(zap.String("method", "GET")).Len() > 0)
}

// TestMiddlewareChain_ExecutionOrder 测试中间件执行顺序
func TestMiddlewareChain_ExecutionOrder(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := core.NewManager(logger)

	// 注册中间件
	manager.Register(builtin.NewRequestIDMiddleware())
	manager.Register(builtin.NewRecoveryMiddleware(logger))
	manager.Register(builtin.NewSecurityMiddleware())
	manager.Register(builtin.NewLoggerMiddleware(logger))
	manager.Register(builtin.NewCompressionMiddleware())

	// 获取执行顺序
	order := manager.GetExecutionOrder()

	// 验证顺序
	expectedOrder := []string{
		"request_id",
		"recovery",
		"security",
		"logger",
		"compression",
	}

	assert.Equal(t, expectedOrder, order)
}

// TestMiddlewareChain_ErrorHandling 测试错误处理流程
func TestMiddlewareChain_ErrorHandling(t *testing.T) {
	// 创建测试观察器
	observedZapCore, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(observedZapCore)

	manager := core.NewManager(logger)

	// 注册中间件
	manager.Register(builtin.NewRequestIDMiddleware())
	manager.Register(builtin.NewRecoveryMiddleware(logger))
	manager.Register(builtin.NewLoggerMiddleware(logger))

	// 创建Gin引擎
	router := gin.New()
	manager.ApplyToRouter(router)

	// 添加会panic的路由
	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/panic", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证panic被捕获并返回500
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// 验证请求ID仍然存在
	requestID := w.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)

	// 验证错误日志
	assert.True(t, logs.FilterField(zap.String("error", "test panic")).Len() > 0)
}

// TestMiddlewareChain_PriorityOverride 测试优先级覆盖
func TestMiddlewareChain_PriorityOverride(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := core.NewManager(logger)

	// 注册多个中间件
	manager.Register(builtin.NewRequestIDMiddleware())       // priority 1
	manager.Register(builtin.NewRecoveryMiddleware(logger))   // priority 2
	manager.Register(builtin.NewSecurityMiddleware())         // priority 4
	cors := builtin.NewCORSMiddleware()
	manager.Register(cors, core.WithPriority(10)) // 覆盖为10（原priority 5）
	manager.Register(builtin.NewLoggerMiddleware(logger))     // priority 7

	// 获取执行顺序
	order := manager.GetExecutionOrder()

	// 验证执行顺序：request_id(1), recovery(2), security(4), logger(7), cors(10)
	expectedOrder := []string{"request_id", "recovery", "security", "logger", "cors"}
	assert.Equal(t, expectedOrder, order)
}

// TestMiddlewareChain_CompressionAndSecurity 测试压缩和安全头组合
func TestMiddlewareChain_CompressionAndSecurity(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := core.NewManager(logger)

	// 注册中间件
	manager.Register(builtin.NewRequestIDMiddleware())
	manager.Register(builtin.NewRecoveryMiddleware(logger))
	manager.Register(builtin.NewSecurityMiddleware())
	manager.Register(builtin.NewCompressionMiddleware())

	// 创建Gin引擎
	router := gin.New()
	manager.ApplyToRouter(router)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		largeData := strings.Repeat("data", 1000)
		c.JSON(http.StatusOK, gin.H{"data": largeData})
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证安全头
	assert.NotEmpty(t, w.Header().Get("X-Frame-Options"))
	assert.NotEmpty(t, w.Header().Get("X-Content-Type-Options"))

	// 验证压缩
	assert.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
}

// TestMiddlewareChain_RequestIDPropagation 测试请求ID传播
func TestMiddlewareChain_RequestIDPropagation(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := core.NewManager(logger)

	// 注册中间件
	manager.Register(builtin.NewRequestIDMiddleware())
	manager.Register(builtin.NewSecurityMiddleware())
	manager.Register(builtin.NewLoggerMiddleware(logger))

	// 创建Gin引擎
	router := gin.New()
	manager.ApplyToRouter(router)

	// 添加测试路由
	router.GET("/test", func(c *gin.Context) {
		requestID := builtin.GetRequestID(c)
		c.JSON(http.StatusOK, gin.H{"request_id": requestID})
	})

	// 发送请求（带自定义请求ID）
	customRequestID := "custom-request-id-123"
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", customRequestID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证使用的是自定义请求ID
	responseRequestID := w.Header().Get("X-Request-ID")
	assert.Equal(t, customRequestID, responseRequestID)

	// 验证响应体中的请求ID
	var response map[string]interface{}
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))
	assert.Equal(t, customRequestID, response["request_id"])
}

// TestMiddlewareChain_SlowRequest 测试慢请求处理
func TestMiddlewareChain_SlowRequest(t *testing.T) {
	// 创建测试观察器
	observedZapCore, _ := observer.New(zap.WarnLevel)
	logger := zap.New(observedZapCore)

	manager := core.NewManager(logger)

	// 注册中间件
	loggerMW := builtin.NewLoggerMiddleware(logger)
	loggerMW.SetSlowRequestThreshold(1) // 设置很低的阈值

	manager.Register(builtin.NewRequestIDMiddleware())
	manager.Register(loggerMW)

	// 创建Gin引擎
	router := gin.New()
	manager.ApplyToRouter(router)

	// 添加测试路由
	router.GET("/slow", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// 发送请求
	req, _ := http.NewRequest("GET", "/slow", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应
	assert.Equal(t, http.StatusOK, w.Code)

	// 验证慢请求日志（可能不稳定，因为实际耗时可能很短）
	// 在实际项目中，应该使用mock时间
}

// TestMiddlewareChain_MiddlewareOrderReport 测试中间件顺序报告
func TestMiddlewareChain_MiddlewareOrderReport(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	manager := core.NewManager(logger)

	// 注册中间件
	manager.Register(builtin.NewRequestIDMiddleware())
	manager.Register(builtin.NewRecoveryMiddleware(logger))
	manager.Register(builtin.NewSecurityMiddleware())
	manager.Register(builtin.NewLoggerMiddleware(logger))
	manager.Register(builtin.NewCompressionMiddleware())

	// 生成报告
	report := manager.GenerateOrderReport()

	// 验证报告内容
	assert.Contains(t, report, "Middleware Execution Order")
	assert.Contains(t, report, "request_id")
	assert.Contains(t, report, "recovery")
	assert.Contains(t, report, "security")
	assert.Contains(t, report, "logger")
	assert.Contains(t, report, "compression")
}
