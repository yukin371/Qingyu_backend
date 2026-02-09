package quota

import (
	"context"
	"errors"
	"fmt"
)

var (
	// ErrInsufficientQuota 配额不足错误
	ErrInsufficientQuota = errors.New("配额不足")
	// ErrQuotaExhausted 配额已用尽错误
	ErrQuotaExhausted = errors.New("配额已用尽")
	// ErrQuotaSuspended 配额已暂停错误
	ErrQuotaSuspended = errors.New("配额已暂停")
)

// Checker 配额检查器接口
// 定义了检查用户配额的抽象行为，不依赖具体实现
type Checker interface {
	// Check 检查用户是否有足够的配额执行操作
	//
	// 参数:
	//   ctx - 上下文
	//   userID - 用户ID
	//   amount - 预估消耗的配额数量（Token数或次数）
	//
	// 返回:
	//   error - 配额不足时返回错误，nil表示检查通过
	Check(ctx context.Context, userID string, amount int) error
}

// ErrorCode 配额错误码
type ErrorCode int

const (
	CodeInsufficientQuota ErrorCode = 40001
	CodeQuotaExhausted    ErrorCode = 40002
	CodeQuotaSuspended    ErrorCode = 40003
)

// QuotaError 配额错误类型
type QuotaError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *QuotaError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *QuotaError) Unwrap() error {
	return e.Err
}
