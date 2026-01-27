package builtin

import (
	"bytes"
	"fmt"
	"runtime/debug"

	"Qingyu_backend/internal/middleware/core"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RecoveryMiddleware 异常恢复中间件
//
// 优先级: 2（最早执行，确保能捕获所有后续中间件的panic）
// 用途: 捕获panic，恢复服务，记录错误日志，返回统一错误响应
type RecoveryMiddleware struct {
	config    *RecoveryConfig
	logger    *zap.Logger
	customRecovery gin.RecoveryFunc
}

// RecoveryConfig 恢复配置
type RecoveryConfig struct {
	// StackSize 堆栈缓冲区大小
	// 用于捕获panic时的堆栈信息
	// 默认: 4096
	StackSize int `yaml:"stack_size"`

	// DisablePrint 是否禁用打印堆栈到标准输出
	// 如果为true，堆栈信息只记录到日志，不打印到控制台
	// 默认: true
	DisablePrint bool `yaml:"disable_print"`
}

// DefaultRecoveryConfig 返回默认恢复配置
func DefaultRecoveryConfig() *RecoveryConfig {
	return &RecoveryConfig{
		StackSize:    4096,
		DisablePrint: true,
	}
}

// NewRecoveryMiddleware 创建新的恢复中间件
func NewRecoveryMiddleware(logger *zap.Logger) *RecoveryMiddleware {
	if logger == nil {
		// 如果没有提供logger，创建一个开发环境的logger
		logger, _ = zap.NewDevelopment()
	}

	return &RecoveryMiddleware{
		config: DefaultRecoveryConfig(),
		logger: logger,
		customRecovery: nil,
	}
}

// Name 返回中间件名称
func (m *RecoveryMiddleware) Name() string {
	return "recovery"
}

// Priority 返回执行优先级
//
// 返回2，确保恢复中间件尽早执行
// 这样可以捕获所有后续中间件和处理器的panic
func (m *RecoveryMiddleware) Priority() int {
	return 2
}

// Handler 返回Gin处理函数
func (m *RecoveryMiddleware) Handler() gin.HandlerFunc {
	// 如果有自定义恢复函数，使用自定义的
	if m.customRecovery != nil {
		return func(c *gin.Context) {
			defer func() {
				if err := recover(); err != nil {
					m.customRecovery(c, err)
				}
			}()
			c.Next()
		}
	}

	// 否则使用默认恢复逻辑
	return func(c *gin.Context) {
		// 捕获panic
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息
				request := m.getRequestInfo(c)

				// 获取堆栈信息
				stack := m.getStackInfo()

				// 记录错误日志
				m.logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("request", request),
					zap.String("stack", stack),
				)

				// 如果不禁用打印，打印堆栈到标准输出
				if !m.config.DisablePrint {
					fmt.Printf("Panic recovered: %v\n", err)
					fmt.Printf("Request: %s\n", request)
					fmt.Printf("Stack:\n%s\n", stack)
				}

				// 设置错误到Context
				// 注意: 这里使用string literal而不是常量，避免循环依赖
				c.Set("middleware_error", err)

				// 中断请求处理
				c.AbortWithStatusJSON(500, gin.H{
					"code":    500,
					"message": "Internal Server Error",
				})
			}
		}()

		c.Next()
	}
}

// getRequestInfo 获取请求信息
func (m *RecoveryMiddleware) getRequestInfo(c *gin.Context) string {
	var buf bytes.Buffer

	// 获取请求方法、路径和查询参数
	fmt.Fprintf(&buf, "%s %s", c.Request.Method, c.Request.URL.Path)
	if c.Request.URL.RawQuery != "" {
		fmt.Fprintf(&buf, "?%s", c.Request.URL.RawQuery)
	}

	// 尝试获取请求ID
	if requestID := GetRequestID(c); requestID != "" {
		fmt.Fprintf(&buf, " [Request-ID: %s]", requestID)
	}

	return buf.String()
}

// getStackInfo 获取堆栈信息
func (m *RecoveryMiddleware) getStackInfo() string {
	if m.config.StackSize <= 0 {
		return string(debug.Stack())
	}

	stack := debug.Stack()
	if len(stack) > m.config.StackSize {
		return string(stack[:m.config.StackSize]) + "... (truncated)"
	}

	return string(stack)
}

// SetCustomRecovery 设置自定义恢复函数
//
// 允许用户自定义panic恢复逻辑
//
// 示例:
//
//	middleware.SetCustomRecovery(func(c *gin.Context, recovered interface{}) {
//	    // 自定义恢复逻辑
//	    c.JSON(500, gin.H{"custom_error": "something went wrong"})
//	})
func (m *RecoveryMiddleware) SetCustomRecovery(recoveryFunc gin.RecoveryFunc) {
	m.customRecovery = recoveryFunc
}

// LoadConfig 从配置加载参数
//
// 实现ConfigurableMiddleware接口
func (m *RecoveryMiddleware) LoadConfig(config map[string]interface{}) error {
	if m.config == nil {
		m.config = &RecoveryConfig{}
	}

	// 加载StackSize
	if stackSize, ok := config["stack_size"].(int); ok {
		m.config.StackSize = stackSize
	}

	// 加载DisablePrint
	if disablePrint, ok := config["disable_print"].(bool); ok {
		m.config.DisablePrint = disablePrint
	}

	return nil
}

// ValidateConfig 验证配置有效性
//
// 实现ConfigurableMiddleware接口
func (m *RecoveryMiddleware) ValidateConfig() error {
	if m.config == nil {
		m.config = DefaultRecoveryConfig()
	}

	// 验证StackSize
	if m.config.StackSize < 0 {
		return fmt.Errorf("stack_size不能为负数")
	}

	return nil
}

// 确保RecoveryMiddleware实现了ConfigurableMiddleware接口
var _ core.ConfigurableMiddleware = (*RecoveryMiddleware)(nil)
