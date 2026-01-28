package errors

import "net/http"

// ErrorCode 错误码类型
type ErrorCode int

// 常用错误码定义
// 格式：6位数字，前2位表示错误类别，后4位为具体错误编号
// 10xxxx: 客户端错误 - 参数/请求相关
// 10xxxx: 客户端错误 - 认证授权相关
// 99xxxx: 服务器错误 - 系统内部错误
// 99xxxx: 服务器错误 - 外部服务错误
const (
	// 成功
	Success ErrorCode = 0

	// 客户端错误 - 参数相关 (10xxxx)
	InvalidParams ErrorCode = 100001 // 无效参数

	// 客户端错误 - 认证授权相关 (10xxxx)
	Unauthorized ErrorCode = 100601 // 未授权
	Forbidden    ErrorCode = 100603 // 禁止访问

	// 客户端错误 - 资源相关 (10xxxx)
	NotFound      ErrorCode = 100401 // 资源不存在
	AlreadyExists ErrorCode = 100201 // 资源已存在
	Conflict      ErrorCode = 100202 // 冲突

	// 认证授权错误 (10xxxx)
	InvalidCredentials ErrorCode = 100611 // 无效凭证
	TokenExpired       ErrorCode = 100612 // Token过期
	TokenInvalid       ErrorCode = 100613 // Token无效
	PasswordTooWeak    ErrorCode = 100114 // 密码强度不足
	AccountLocked      ErrorCode = 100615 // 账户已锁定
	AccountDisabled    ErrorCode = 100616 // 账户已禁用

	// 业务逻辑错误 (10xxxx)
	InsufficientBalance ErrorCode = 100301 // 余额不足
	InsufficientQuota   ErrorCode = 100302 // 配额不足
	WalletFrozen        ErrorCode = 100303 // 钱包已冻结
	ContentNotPublished ErrorCode = 100304 // 内容未发布
	ChapterLocked       ErrorCode = 100305 // 章节已锁定

	// 内容审核错误 (10xxxx)
	ContentPendingReview ErrorCode = 100401 // 内容待审核
	ContentRejected      ErrorCode = 100403 // 内容被拒绝
	ContentViolation     ErrorCode = 100405 // 内容违规

	// 服务器错误 - 系统错误 (99xxxx)
	InternalError ErrorCode = 995001 // 内部错误
	DatabaseError ErrorCode = 995002 // 数据库错误
	RedisError    ErrorCode = 995004 // Redis错误

	// 服务器错误 - 外部服务 (99xxxx/99xxxx)
	ExternalAPIError  ErrorCode = 990001 // 外部API错误
	RateLimitExceeded ErrorCode = 995007 // 超出频率限制
)

// ErrorInfo 错误信息
type ErrorInfo struct {
	Code       ErrorCode   `json:"code"`
	Message    string      `json:"message"`
	HTTPStatus int         `json:"-"`
	Details    interface{} `json:"details,omitempty"`
}

// DefaultErrorMessages 默认错误消息
var DefaultErrorMessages = map[ErrorCode]string{
	InvalidParams:     "请求参数无效",
	Unauthorized:      "未授权访问",
	Forbidden:         "禁止访问",
	NotFound:          "资源不存在",
	AlreadyExists:     "资源已存在",
	Conflict:          "请求冲突",
	RateLimitExceeded: "请求过于频繁",

	InvalidCredentials: "用户名或密码错误",
	TokenExpired:       "登录已过期，请重新登录",
	TokenInvalid:       "无效的登录凭证",
	PasswordTooWeak:    "密码强度不足",
	AccountLocked:      "账户已被锁定",
	AccountDisabled:    "账户已被禁用",

	InsufficientBalance: "余额不足",
	InsufficientQuota:   "配额不足",
	WalletFrozen:        "钱包已被冻结",
	ContentNotPublished: "内容尚未发布",
	ChapterLocked:       "章节已锁定",

	ContentPendingReview: "内容待审核",
	ContentRejected:      "内容审核未通过",
	ContentViolation:     "内容违规",

	InternalError:    "服务器内部错误",
	DatabaseError:    "数据库错误",
	RedisError:       "缓存服务错误",
	ExternalAPIError: "外部服务错误",
}

// DefaultHTTPStatus 默认HTTP状态码
var DefaultHTTPStatus = map[ErrorCode]int{
	InvalidParams:     http.StatusBadRequest,
	Unauthorized:      http.StatusUnauthorized,
	Forbidden:         http.StatusForbidden,
	NotFound:          http.StatusNotFound,
	AlreadyExists:     http.StatusConflict,
	Conflict:          http.StatusConflict,
	RateLimitExceeded: http.StatusTooManyRequests,

	InvalidCredentials: http.StatusUnauthorized,
	TokenExpired:       http.StatusUnauthorized,
	TokenInvalid:       http.StatusUnauthorized,
	PasswordTooWeak:    http.StatusBadRequest,
	AccountLocked:      http.StatusForbidden,
	AccountDisabled:    http.StatusForbidden,

	InsufficientBalance: http.StatusBadRequest,
	InsufficientQuota:   http.StatusBadRequest,
	WalletFrozen:        http.StatusForbidden,
	ContentNotPublished: http.StatusForbidden,
	ChapterLocked:       http.StatusForbidden,

	ContentPendingReview: http.StatusAccepted,
	ContentRejected:      http.StatusForbidden,
	ContentViolation:     http.StatusForbidden,

	InternalError:    http.StatusInternalServerError,
	DatabaseError:    http.StatusInternalServerError,
	RedisError:       http.StatusInternalServerError,
	ExternalAPIError: http.StatusBadGateway,
}

// GetHTTPStatus 获取错误码对应的HTTP状态码
func GetHTTPStatus(code ErrorCode) int {
	if status, ok := DefaultHTTPStatus[code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// GetDefaultMessage 获取默认错误消息
func GetDefaultMessage(code ErrorCode) string {
	if msg, ok := DefaultErrorMessages[code]; ok {
		return msg
	}
	return "未知错误"
}
