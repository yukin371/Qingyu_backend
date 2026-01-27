package ratelimit

import (
	"context"
	"time"
)

// LimiterStats 限流器统计信息
type LimiterStats struct {
	// 总请求数
	TotalRequests int64
	// 被拒绝的请求数
	RejectedRequests int64
	// 允许的请求数
	AllowedRequests int64
	// 最后请求时间
	LastRequestTime time.Time
}

// Limiter 限流器接口
//
// 定义了限流器的核心行为，支持多种限流算法
type Limiter interface {
	// Allow 检查是否允许请求
	//
	// key: 限流键（如IP、用户ID等）
	// 返回true表示允许请求，false表示被限流
	Allow(key string) bool

	// Wait 等待直到可以处理请求
	//
	// ctx: 上下文，支持超时和取消
	// key: 限流键
	// 返回错误表示：超时、上下文取消等
	Wait(ctx context.Context, key string) error

	// Reset 重置指定key的限流状态
	//
	// key: 限流键
	Reset(key string)

	// GetStats 获取指定key的统计信息
	//
	// key: 限流键
	// 返回统计信息，如果key不存在返回nil
	GetStats(key string) *LimiterStats
}

// LimiterFactory 限流器工厂接口
//
// 用于创建不同类型的限流器
type LimiterFactory interface {
	// CreateLimiter 创建限流器
	//
	// config: 限流配置
	// 返回限流器实例
	CreateLimiter(config *RateLimitConfig) (Limiter, error)
}

// KeyFunc 限流键生成函数
//
// 根据请求上下文生成限流键
type KeyFunc func(*RateLimitContext) string

// RateLimitContext 限流上下文
//
// 封装Gin的Context，提供更便捷的访问
type RateLimitContext struct {
	// Gin的Context
	Context interface{} // 避免循环依赖，使用interface{}

	// 客户端IP
	ClientIP string

	// 用户ID
	UserID string

	// 请求路径
	Path string

	// 请求方法
	Method string

	// 额外的元数据
	Metadata map[string]interface{}
}

// KeyFuncType 限流键类型
type KeyFuncType string

const (
	// KeyFuncIP 基于IP的限流
	KeyFuncIP KeyFuncType = "ip"
	// KeyFuncUser 基于用户ID的限流
	KeyFuncUser KeyFuncType = "user"
	// KeyFuncPath 基于路径的限流
	KeyFuncPath KeyFuncType = "path"
	// KeyFuncIPPath 基于IP+路径的组合限流
	KeyFuncIPPath KeyFuncType = "ip_path"
	// KeyFuncUserPath 基于用户ID+路径的组合限流
	KeyFuncUserPath KeyFuncType = "user_path"
)

// GetKeyFunc 根据类型获取KeyFunc
func GetKeyFunc(keyFuncType KeyFuncType) KeyFunc {
	switch keyFuncType {
	case KeyFuncIP:
		return func(ctx *RateLimitContext) string {
			return ctx.ClientIP
		}
	case KeyFuncUser:
		return func(ctx *RateLimitContext) string {
			if ctx.UserID != "" {
				return "user:" + ctx.UserID
			}
			return ctx.ClientIP
		}
	case KeyFuncPath:
		return func(ctx *RateLimitContext) string {
			return ctx.Path
		}
	case KeyFuncIPPath:
		return func(ctx *RateLimitContext) string {
			return ctx.ClientIP + ":" + ctx.Path
		}
	case KeyFuncUserPath:
		return func(ctx *RateLimitContext) string {
			if ctx.UserID != "" {
				return "user:" + ctx.UserID + ":" + ctx.Path
			}
			return ctx.ClientIP + ":" + ctx.Path
		}
	default:
		return func(ctx *RateLimitContext) string {
			return ctx.ClientIP
		}
	}
}
